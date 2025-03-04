package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/yuroku/internal/usecase/port"
)

// AuthController は認証関連のコントローラーです
type AuthController struct {
	authUseCase port.AuthInputPort
}

// NewAuthController は新しい認証コントローラーを作成します
func NewAuthController(authUseCase port.AuthInputPort) *AuthController {
	return &AuthController{
		authUseCase: authUseCase,
	}
}

// Register はユーザー登録を処理します
func (c *AuthController) Register(ctx *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_INPUT",
				"message": "入力データが無効です: " + err.Error(),
			},
		})
		return
	}

	// ユースケースを呼び出し
	err := c.authUseCase.Register(ctx.Request.Context(), input.Name, input.Email, input.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "REGISTRATION_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "ユーザー登録が完了しました",
	})
}

// Login はユーザーログインを処理します
func (c *AuthController) Login(ctx *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_INPUT",
				"message": "入力データが無効です: " + err.Error(),
			},
		})
		return
	}

	// ユースケースを呼び出し
	token, refreshToken, err := c.authUseCase.Login(ctx.Request.Context(), input.Email, input.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "LOGIN_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"access_token":  token,
			"refresh_token": refreshToken,
		},
		"message": "ログインに成功しました",
	})
}

// RefreshToken はトークンを更新します
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_INPUT",
				"message": "入力データが無効です: " + err.Error(),
			},
		})
		return
	}

	// ユースケースを呼び出し
	newToken, newRefreshToken, err := c.authUseCase.RefreshToken(ctx.Request.Context(), input.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "TOKEN_REFRESH_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"access_token":  newToken,
			"refresh_token": newRefreshToken,
		},
		"message": "トークンの更新に成功しました",
	})
}

// Logout はユーザーログアウトを処理します
func (c *AuthController) Logout(ctx *gin.Context) {
	// Authorizationヘッダーからトークンを取得
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "MISSING_TOKEN",
				"message": "認証トークンがありません",
			},
		})
		return
	}

	// "Bearer "プレフィックスを削除
	token := authHeader[7:]

	// ユースケースを呼び出し
	err := c.authUseCase.Logout(ctx.Request.Context(), token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "LOGOUT_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "ログアウトに成功しました",
	})
}

// GetCurrentUser は現在のユーザー情報を取得します
func (c *AuthController) GetCurrentUser(ctx *gin.Context) {
	// コンテキストからユーザーIDを取得
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "認証されていません",
			},
		})
		return
	}

	// ユースケースを呼び出し
	user, err := c.authUseCase.GetUserByID(ctx.Request.Context(), userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "USER_FETCH_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
		"message": "ユーザー情報の取得に成功しました",
	})
}

// UpdateProfile はユーザープロフィールを更新します
func (c *AuthController) UpdateProfile(ctx *gin.Context) {
	// コンテキストからユーザーIDを取得
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "認証されていません",
			},
		})
		return
	}

	var input struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required,email"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_INPUT",
				"message": "入力データが無効です: " + err.Error(),
			},
		})
		return
	}

	// ユースケースを呼び出し
	err := c.authUseCase.UpdateProfile(ctx.Request.Context(), userID.(string), input.Name, input.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "PROFILE_UPDATE_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "プロフィールの更新に成功しました",
	})
}

// ChangePassword はパスワードを変更します
func (c *AuthController) ChangePassword(ctx *gin.Context) {
	// コンテキストからユーザーIDを取得
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "認証されていません",
			},
		})
		return
	}

	var input struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=8"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_INPUT",
				"message": "入力データが無効です: " + err.Error(),
			},
		})
		return
	}

	// ユースケースを呼び出し
	err := c.authUseCase.ChangePassword(ctx.Request.Context(), userID.(string), input.CurrentPassword, input.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "PASSWORD_CHANGE_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "パスワードの変更に成功しました",
	})
}

// DeleteAccount はアカウントを削除します
func (c *AuthController) DeleteAccount(ctx *gin.Context) {
	// コンテキストからユーザーIDを取得
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "認証されていません",
			},
		})
		return
	}

	var input struct {
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_INPUT",
				"message": "入力データが無効です: " + err.Error(),
			},
		})
		return
	}

	// ユースケースを呼び出し
	err := c.authUseCase.DeleteAccount(ctx.Request.Context(), userID.(string), input.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "ACCOUNT_DELETION_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "アカウントの削除に成功しました",
	})
}
