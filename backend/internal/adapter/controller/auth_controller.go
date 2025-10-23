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

	if !ValidateBindJSON(ctx, &input) {
		return
	}

	// ユースケースを呼び出し
	err := c.authUseCase.Register(ctx.Request.Context(), input.Name, input.Email, input.Password)
	if err != nil {
		RespondWithAppError(ctx, err)
		return
	}

	RespondWithSuccess(ctx, http.StatusCreated, nil, "ユーザー登録が完了しました")
}

// Login はユーザーログインを処理します
func (c *AuthController) Login(ctx *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if !ValidateBindJSON(ctx, &input) {
		return
	}

	// ユースケースを呼び出し
	token, refreshToken, err := c.authUseCase.Login(ctx.Request.Context(), input.Email, input.Password)
	if err != nil {
		RespondWithAppError(ctx, err)
		return
	}

	RespondWithSuccess(ctx, http.StatusOK, gin.H{
		"access_token":  token,
		"refresh_token": refreshToken,
	}, "ログインに成功しました")
}

// RefreshToken はトークンを更新します
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if !ValidateBindJSON(ctx, &input) {
		return
	}

	// ユースケースを呼び出し
	newToken, newRefreshToken, err := c.authUseCase.RefreshToken(ctx.Request.Context(), input.RefreshToken)
	if err != nil {
		RespondWithAppError(ctx, err)
		return
	}

	RespondWithSuccess(ctx, http.StatusOK, gin.H{
		"access_token":  newToken,
		"refresh_token": newRefreshToken,
	}, "トークンの更新に成功しました")
}

// Logout はユーザーログアウトを処理します
func (c *AuthController) Logout(ctx *gin.Context) {
	// Authorizationヘッダーからトークンを取得
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		RespondWithError(ctx, http.StatusBadRequest, "MISSING_TOKEN", "認証トークンがありません")
		return
	}

	// "Bearer "プレフィックスを削除
	token := authHeader[7:]

	// ユースケースを呼び出し
	err := c.authUseCase.Logout(ctx.Request.Context(), token)
	if err != nil {
		RespondWithAppError(ctx, err)
		return
	}

	RespondWithSuccess(ctx, http.StatusOK, nil, "ログアウトに成功しました")
}

// GetCurrentUser は現在のユーザー情報を取得します
func (c *AuthController) GetCurrentUser(ctx *gin.Context) {
	// コンテキストからユーザーIDを取得
	userID, ok := GetUserID(ctx)
	if !ok {
		return
	}

	// ユースケースを呼び出し
	user, err := c.authUseCase.GetUserByID(ctx.Request.Context(), userID)
	if err != nil {
		RespondWithAppError(ctx, err)
		return
	}

	RespondWithSuccess(ctx, http.StatusOK, gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	}, "ユーザー情報の取得に成功しました")
}

// UpdateProfile はユーザープロフィールを更新します
func (c *AuthController) UpdateProfile(ctx *gin.Context) {
	// コンテキストからユーザーIDを取得
	userID, ok := GetUserID(ctx)
	if !ok {
		return
	}

	var input struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required,email"`
	}

	if !ValidateBindJSON(ctx, &input) {
		return
	}

	// ユースケースを呼び出し
	err := c.authUseCase.UpdateProfile(ctx.Request.Context(), userID, input.Name, input.Email)
	if err != nil {
		RespondWithAppError(ctx, err)
		return
	}

	RespondWithSuccess(ctx, http.StatusOK, nil, "プロフィールの更新に成功しました")
}

// ChangePassword はパスワードを変更します
func (c *AuthController) ChangePassword(ctx *gin.Context) {
	// コンテキストからユーザーIDを取得
	userID, ok := GetUserID(ctx)
	if !ok {
		return
	}

	var input struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=8"`
	}

	if !ValidateBindJSON(ctx, &input) {
		return
	}

	// ユースケースを呼び出し
	err := c.authUseCase.ChangePassword(ctx.Request.Context(), userID, input.CurrentPassword, input.NewPassword)
	if err != nil {
		RespondWithAppError(ctx, err)
		return
	}

	RespondWithSuccess(ctx, http.StatusOK, nil, "パスワードの変更に成功しました")
}

// DeleteAccount はアカウントを削除します
func (c *AuthController) DeleteAccount(ctx *gin.Context) {
	// コンテキストからユーザーIDを取得
	userID, ok := GetUserID(ctx)
	if !ok {
		return
	}

	var input struct {
		Password string `json:"password" binding:"required"`
	}

	if !ValidateBindJSON(ctx, &input) {
		return
	}

	// ユースケースを呼び出し
	err := c.authUseCase.DeleteAccount(ctx.Request.Context(), userID, input.Password)
	if err != nil {
		RespondWithAppError(ctx, err)
		return
	}

	RespondWithSuccess(ctx, http.StatusOK, nil, "アカウントの削除に成功しました")
}
