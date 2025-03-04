package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/yuroku/internal/domain/service"
)

// AuthMiddleware は認証ミドルウェアを提供します
type AuthMiddleware struct {
	authService *service.AuthService
}

// NewAuthMiddleware は新しい認証ミドルウェアを作成します
func NewAuthMiddleware(authService *service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// RequireAuth はユーザー認証を要求するミドルウェアです
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Authorizationヘッダーを取得
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "MISSING_TOKEN",
					"message": "認証トークンがありません",
				},
			})
			return
		}

		// Bearer トークンを取得
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "INVALID_TOKEN_FORMAT",
					"message": "トークン形式が無効です",
				},
			})
			return
		}

		token := parts[1]

		// トークンを検証
		userID, err := m.authService.VerifyToken(ctx.Request.Context(), token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "INVALID_TOKEN",
					"message": err.Error(),
				},
			})
			return
		}

		// ユーザーIDをコンテキストに設定
		ctx.Set("userID", userID)
		ctx.Next()
	}
}

// OptionalAuth はユーザー認証を任意とするミドルウェアです
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Authorizationヘッダーを取得
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			// トークンがなければ続行
			ctx.Next()
			return
		}

		// Bearer トークンを取得
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			// 無効なトークン形式なら続行
			ctx.Next()
			return
		}

		token := parts[1]

		// トークンを検証
		userID, err := m.authService.VerifyToken(ctx.Request.Context(), token)
		if err != nil {
			// 無効なトークンなら続行
			ctx.Next()
			return
		}

		// ユーザーIDをコンテキストに設定
		ctx.Set("userID", userID)
		ctx.Next()
	}
}
