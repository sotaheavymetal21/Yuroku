package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/yourusername/yuroku/internal/adapter/controller"
	"github.com/yourusername/yuroku/internal/adapter/gateway"
	"github.com/yourusername/yuroku/internal/adapter/presenter"
	"github.com/yourusername/yuroku/internal/domain/service"
	"github.com/yourusername/yuroku/internal/infrastructure/database"
	"github.com/yourusername/yuroku/internal/infrastructure/middleware"
	"github.com/yourusername/yuroku/internal/infrastructure/router"
	"github.com/yourusername/yuroku/internal/infrastructure/storage"
	"github.com/yourusername/yuroku/internal/usecase/interactor"
)

func main() {
	// 環境変数を読み込み
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// MongoDBに接続
	mongoClient, err := database.NewMongoClient(os.Getenv("MONGODB_URI"))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())

	// データベースを取得
	db := mongoClient.Database(os.Getenv("MONGODB_DATABASE"))

	// ストレージを初期化
	fileStorage, err := storage.NewLocalFileStorage(os.Getenv("UPLOAD_DIR"))
	if err != nil {
		log.Fatalf("Failed to initialize file storage: %v", err)
	}

	// リポジトリを初期化
	userRepo := gateway.NewMongoUserRepository(db)
	onsenLogRepo := gateway.NewMongoOnsenLogRepository(db)
	onsenImageRepo := gateway.NewMongoOnsenImageRepository(db)

	// ドメインサービスを初期化
	jwtSecret := os.Getenv("JWT_SECRET")
	authService := service.NewAuthService(userRepo, jwtSecret)
	onsenLogService := service.NewOnsenLogService(onsenLogRepo, onsenImageRepo)
	onsenImageService := service.NewOnsenImageService(onsenImageRepo, onsenLogRepo, fileStorage)

	// プレゼンターを初期化
	authPresenter := presenter.NewAuthPresenter()
	onsenLogPresenter := presenter.NewOnsenLogPresenter()
	onsenImagePresenter := presenter.NewOnsenImagePresenter()

	// プレゼンターをOutputPortにアダプト
	authOutputPort := presenter.NewAuthOutputAdapter(authPresenter)
	onsenLogOutputPort := presenter.NewOnsenLogOutputAdapter(onsenLogPresenter)
	onsenImageOutputPort := presenter.NewOnsenImageOutputAdapter(onsenImagePresenter)

	// JWTの設定
	accessTokenDuration := 15 * time.Minute    // アクセストークンの有効期限
	refreshTokenDuration := 7 * 24 * time.Hour // リフレッシュトークンの有効期限

	// ユースケースを初期化
	authInteractor := interactor.NewAuthInteractor(authService, authOutputPort, jwtSecret, accessTokenDuration, refreshTokenDuration)
	onsenLogInteractor := interactor.NewOnsenLogInteractor(onsenLogService, onsenImageService, onsenLogOutputPort)
	onsenImageInteractor := interactor.NewOnsenImageInteractor(onsenImageService, onsenImageOutputPort)

	// コントローラーを初期化
	authController := controller.NewAuthController(authInteractor)
	onsenLogController := controller.NewOnsenLogController(onsenLogInteractor)
	onsenImageController := controller.NewOnsenImageController(onsenImageInteractor)

	// ミドルウェアを初期化
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// ルーターを初期化
	r := router.NewRouter(
		authMiddleware,
		authController,
		onsenLogController,
		onsenImageController,
	)

	// ルートを設定
	r.SetupRoutes()

	// サーバーを起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// シグナル処理を設定
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 別のゴルーチンでサーバーを起動
	go func() {
		log.Printf("Server is running on port %s", port)
		if err := r.Run(":" + port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// シグナルを待機
	<-quit
	log.Println("Shutting down server...")

	// グレースフルシャットダウンのためのコンテキスト
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// MongoDBの接続を閉じる
	if err := mongoClient.Disconnect(ctx); err != nil {
		log.Fatalf("Failed to disconnect from MongoDB: %v", err)
	}

	log.Println("Server exited properly")
}
