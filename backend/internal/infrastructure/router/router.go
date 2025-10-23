package router

import (
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/yourusername/yuroku/internal/adapter/controller"
	"github.com/yourusername/yuroku/internal/adapter/gateway"
	"github.com/yourusername/yuroku/internal/adapter/presenter"
	"github.com/yourusername/yuroku/internal/domain/service"
	"github.com/yourusername/yuroku/internal/infrastructure/database"
	"github.com/yourusername/yuroku/internal/infrastructure/middleware"
	"github.com/yourusername/yuroku/internal/infrastructure/storage"
	"github.com/yourusername/yuroku/internal/usecase/interactor"
)

// Router はAPIルーターを提供します
type Router struct {
	engine               *gin.Engine
	authMiddleware       *middleware.AuthMiddleware
	authController       *controller.AuthController
	onsenLogController   *controller.OnsenLogController
	onsenImageController *controller.OnsenImageController
}

// NewRouter は新しいAPIルーターを作成します
func NewRouter(
	authMiddleware *middleware.AuthMiddleware,
	authController *controller.AuthController,
	onsenLogController *controller.OnsenLogController,
	onsenImageController *controller.OnsenImageController,
) *Router {
	engine := gin.Default()

	// CORSミドルウェアを設定
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"} // フロントエンドのオリジン
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	config.AllowCredentials = true
	engine.Use(cors.New(config))

	return &Router{
		engine:               engine,
		authMiddleware:       authMiddleware,
		authController:       authController,
		onsenLogController:   onsenLogController,
		onsenImageController: onsenImageController,
	}
}

// SetupRoutes はAPIルートを設定します
func (r *Router) SetupRoutes() {
	// APIプレフィックス
	api := r.engine.Group("/api")

	// 認証関連のルート (/api/auth/...)
	auth := api.Group("/auth")
	{
		auth.POST("/register", r.authController.Register)
		auth.POST("/login", r.authController.Login)
		auth.POST("/refresh", r.authController.RefreshToken)
		auth.POST("/logout", r.authMiddleware.RequireAuth(), r.authController.Logout)

		// ユーザープロフィール関連
		profile := auth.Group("/profile", r.authMiddleware.RequireAuth())
		{
			profile.GET("", r.authController.GetCurrentUser)
			profile.PUT("", r.authController.UpdateProfile)
			profile.PUT("/password", r.authController.ChangePassword)
			profile.DELETE("", r.authController.DeleteAccount)
		}
	}

	// 温泉メモ関連のルート
	onsenLogs := api.Group("/onsen_logs", r.authMiddleware.RequireAuth())
	{
		onsenLogs.POST("", r.onsenLogController.CreateOnsenLog)
		onsenLogs.GET("", r.onsenLogController.GetOnsenLogs)
		onsenLogs.GET("/filter", r.onsenLogController.GetFilteredOnsenLogs)
		onsenLogs.GET("/export", r.onsenLogController.ExportOnsenLogs)
		onsenLogs.GET("/:id", r.onsenLogController.GetOnsenLog)
		onsenLogs.PUT("/:id", r.onsenLogController.UpdateOnsenLog)
		onsenLogs.DELETE("/:id", r.onsenLogController.DeleteOnsenLog)
	}

	// 温泉画像関連のルート
	onsenImages := api.Group("/onsen_images", r.authMiddleware.RequireAuth())
	{
		onsenImages.POST("/:onsen_id", r.onsenImageController.UploadImage)
		onsenImages.GET("/:onsen_id", r.onsenImageController.GetImagesByOnsenID)
		onsenImages.DELETE("/:image_id", r.onsenImageController.DeleteImage)
	}

	// ヘルスチェック
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
}

// Run はサーバーを起動します
func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}

// SetupDependencies はアプリケーションの依存関係を設定します
func SetupDependencies() (*Router, error) {
	// データベース接続を取得
	db, err := database.GetMongoDB()
	if err != nil {
		return nil, err
	}

	// ストレージを設定
	localStorage, err := storage.NewLocalFileStorage("./uploads")
	if err != nil {
		return nil, err
	}
	storageRepo := gateway.NewLocalStorageRepository(localStorage)

	// リポジトリを作成
	userRepo := gateway.NewMongoUserRepository(db)
	onsenLogRepo := gateway.NewMongoOnsenLogRepository(db)
	onsenImageRepo := gateway.NewMongoOnsenImageRepository(db)

	// JWT設定
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key" // 開発用のデフォルト値
	}
	accessTokenDuration := 15 * time.Minute    // アクセストークンの有効期限
	refreshTokenDuration := 7 * 24 * time.Hour // リフレッシュトークンの有効期限

	// ドメインサービスを作成
	authService := service.NewAuthService(userRepo, jwtSecret)
	onsenLogService := service.NewOnsenLogService(onsenLogRepo, onsenImageRepo)
	onsenImageService := service.NewOnsenImageService(onsenImageRepo, onsenLogRepo, storageRepo)

	// プレゼンターを作成
	authPresenter := presenter.NewAuthPresenter()
	onsenLogPresenter := presenter.NewOnsenLogPresenter()
	onsenImagePresenter := presenter.NewOnsenImagePresenter()

	// プレゼンターをOutputPortにアダプト
	authOutputPort := presenter.NewAuthOutputAdapter(authPresenter)
	onsenLogOutputPort := presenter.NewOnsenLogOutputAdapter(onsenLogPresenter)
	onsenImageOutputPort := presenter.NewOnsenImageOutputAdapter(onsenImagePresenter)

	// ユースケースを作成
	authInteractor := interactor.NewAuthInteractor(
		authService,
		authOutputPort,
		jwtSecret,
		accessTokenDuration,
		refreshTokenDuration,
	)
	onsenLogInteractor := interactor.NewOnsenLogInteractor(onsenLogService, onsenImageService, onsenLogOutputPort)
	onsenImageInteractor := interactor.NewOnsenImageInteractor(onsenImageService, onsenImageOutputPort)

	// コントローラーを作成
	authController := controller.NewAuthController(authInteractor)
	onsenLogController := controller.NewOnsenLogController(onsenLogInteractor)
	onsenImageController := controller.NewOnsenImageController(onsenImageInteractor)

	// 認証ミドルウェアを作成
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// ルーターを作成
	router := NewRouter(
		authMiddleware,
		authController,
		onsenLogController,
		onsenImageController,
	)

	return router, nil
}
