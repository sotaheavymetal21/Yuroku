package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/yourusername/yuroku/internal/adapter/controller"
	"github.com/yourusername/yuroku/internal/infrastructure/middleware"
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
	// APIバージョンプレフィックス
	api := r.engine.Group("/api/v1")

	// 認証関連のルート
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

		// 温泉画像関連のルート
		images := onsenLogs.Group("/:onsen_id/images")
		{
			images.POST("", r.onsenImageController.UploadImage)
			images.GET("", r.onsenImageController.GetImagesByOnsenID)
			images.DELETE("/:image_id", r.onsenImageController.DeleteImage)
		}
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
