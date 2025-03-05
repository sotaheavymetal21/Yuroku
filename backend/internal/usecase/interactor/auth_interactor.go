package interactor

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yourusername/yuroku/internal/domain/service"
	"github.com/yourusername/yuroku/internal/usecase/port"
)

// AuthInteractor は認証ユースケースのインタラクターです
type AuthInteractor struct {
	authService   *service.AuthService
	outputPort    port.AuthOutputPort
	jwtSecret     string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

// NewAuthInteractor は新しい認証インタラクターを作成します
func NewAuthInteractor(
	authService *service.AuthService,
	outputPort port.AuthOutputPort,
	jwtSecret string,
	accessExpiry time.Duration,
	refreshExpiry time.Duration,
) *AuthInteractor {
	return &AuthInteractor{
		authService:   authService,
		outputPort:    outputPort,
		jwtSecret:     jwtSecret,
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}
}

// Register はユーザー登録を行います
func (i *AuthInteractor) Register(ctx context.Context, name, email, password string) error {
	// 入力値のバリデーション
	if name == "" || email == "" || password == "" {
		err := errors.New("名前、メールアドレス、パスワードは必須です")
		_ = i.outputPort.PresentError(ctx, err)
		return err
	}

	// ドメインサービスを呼び出し
	user, err := i.authService.Register(ctx, name, email, password)
	if err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return err
	}

	// JWTトークンを生成
	accessToken, err := i.generateAccessToken(user.UUID)
	if err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return err
	}

	refreshToken, err := i.generateRefreshToken(user.UUID)
	if err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return err
	}

	// 出力データを作成
	outputData := port.AuthOutputData{
		UserID:       user.UUID,
		Email:        user.Email,
		Name:         user.Name,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	// 出力ポートを呼び出し
	if err := i.outputPort.PresentRegister(ctx, outputData); err != nil {
		return err
	}

	return nil
}

// Login はユーザーログインを行います
func (i *AuthInteractor) Login(ctx context.Context, email, password string) (string, string, error) {
	// 入力値のバリデーション
	if email == "" || password == "" {
		err := errors.New("メールアドレスとパスワードは必須です")
		_ = i.outputPort.PresentError(ctx, err)
		return "", "", err
	}

	// ドメインサービスを呼び出し
	user, err := i.authService.Login(ctx, email, password)
	if err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return "", "", err
	}

	// JWTトークンを生成
	accessToken, err := i.generateAccessToken(user.UUID)
	if err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return "", "", err
	}

	refreshToken, err := i.generateRefreshToken(user.UUID)
	if err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return "", "", err
	}

	// 出力データを作成
	outputData := port.AuthOutputData{
		UserID:       user.UUID,
		Email:        user.Email,
		Name:         user.Name,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	// 出力ポートを呼び出し
	if err := i.outputPort.PresentLogin(ctx, outputData); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// VerifyToken はトークンを検証します
func (i *AuthInteractor) VerifyToken(ctx context.Context, tokenString string) (string, error) {
	// トークンを解析
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 署名アルゴリズムを検証
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(i.jwtSecret), nil
	})

	if err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return "", err
	}

	// トークンの有効性を検証
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// ユーザーIDを取得
		userID, ok := claims["sub"].(string)
		if !ok {
			err := errors.New("invalid token claims")
			_ = i.outputPort.PresentError(ctx, err)
			return "", err
		}
		return userID, nil
	}

	err = errors.New("invalid token")
	_ = i.outputPort.PresentError(ctx, err)
	return "", err
}

// RefreshToken はトークンを更新します
func (i *AuthInteractor) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	// リフレッシュトークンを検証
	userID, err := i.VerifyToken(ctx, refreshToken)
	if err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return "", "", err
	}

	// ユーザー情報を取得
	user, err := i.authService.GetUserByID(ctx, userID)
	if err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return "", "", err
	}

	// 新しいトークンを生成
	accessToken, err := i.generateAccessToken(user.UUID)
	if err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return "", "", err
	}

	newRefreshToken, err := i.generateRefreshToken(user.UUID)
	if err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return "", "", err
	}

	// 出力データを作成
	outputData := port.AuthOutputData{
		UserID:       user.UUID,
		Email:        user.Email,
		Name:         user.Name,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}

	// 出力ポートを呼び出し
	if err := i.outputPort.PresentRefreshToken(ctx, outputData); err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

// generateAccessToken はアクセストークンを生成します
func (i *AuthInteractor) generateAccessToken(userID string) (string, error) {
	// トークンの有効期限を設定
	expirationTime := time.Now().Add(i.accessExpiry)

	// クレームを作成
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": expirationTime.Unix(),
		"iat": time.Now().Unix(),
		"typ": "access",
	}

	// トークンを生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// トークンに署名
	tokenString, err := token.SignedString([]byte(i.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// generateRefreshToken はリフレッシュトークンを生成します
func (i *AuthInteractor) generateRefreshToken(userID string) (string, error) {
	// トークンの有効期限を設定
	expirationTime := time.Now().Add(i.refreshExpiry)

	// クレームを作成
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": expirationTime.Unix(),
		"iat": time.Now().Unix(),
		"typ": "refresh",
	}

	// トークンを生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// トークンに署名
	tokenString, err := token.SignedString([]byte(i.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Logout はユーザーログアウトを行います
func (i *AuthInteractor) Logout(ctx context.Context, token string) error {
	// トークンを検証
	_, err := i.VerifyToken(ctx, token)
	if err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return err
	}

	// 実際のアプリケーションでは、トークンをブラックリストに追加するなどの処理を行う
	// ここではシンプルに成功を返す

	return nil
}

// GetUserByID はIDでユーザーを取得します
func (i *AuthInteractor) GetUserByID(ctx context.Context, id string) (port.UserOutputData, error) {
	// ドメインサービスを呼び出し
	user, err := i.authService.GetUserByID(ctx, id)
	if err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return port.UserOutputData{}, err
	}

	// 出力データを作成
	outputData := port.UserOutputData{
		ID:    user.UUID,
		Name:  user.Name,
		Email: user.Email,
	}

	return outputData, nil
}

// UpdateProfile はユーザープロフィールを更新します
func (i *AuthInteractor) UpdateProfile(ctx context.Context, userID, name, email string) error {
	// 入力値のバリデーション
	if userID == "" || name == "" || email == "" {
		err := errors.New("ユーザーID、名前、メールアドレスは必須です")
		_ = i.outputPort.PresentError(ctx, err)
		return err
	}

	// ユーザーを取得
	user, err := i.authService.GetUserByID(ctx, userID)
	if err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return err
	}

	// プロフィールを更新
	user.UpdateProfile(name, email)

	// ユーザーを保存
	if err := i.authService.UpdateUser(ctx, user); err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return err
	}

	return nil
}

// ChangePassword はパスワードを変更します
func (i *AuthInteractor) ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error {
	// 入力値のバリデーション
	if userID == "" || currentPassword == "" || newPassword == "" {
		err := errors.New("ユーザーID、現在のパスワード、新しいパスワードは必須です")
		_ = i.outputPort.PresentError(ctx, err)
		return err
	}

	// ユーザーを取得
	user, err := i.authService.GetUserByID(ctx, userID)
	if err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return err
	}

	// 現在のパスワードを検証
	if !user.ComparePassword(currentPassword) {
		err := errors.New("現在のパスワードが正しくありません")
		_ = i.outputPort.PresentError(ctx, err)
		return err
	}

	// パスワードを更新
	if err := user.UpdatePassword(newPassword); err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return err
	}

	// ユーザーを保存
	if err := i.authService.UpdateUser(ctx, user); err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return err
	}

	return nil
}

// DeleteAccount はアカウントを削除します
func (i *AuthInteractor) DeleteAccount(ctx context.Context, userID, password string) error {
	// 入力値のバリデーション
	if userID == "" || password == "" {
		err := errors.New("ユーザーIDとパスワードは必須です")
		_ = i.outputPort.PresentError(ctx, err)
		return err
	}

	// ユーザーを取得
	user, err := i.authService.GetUserByID(ctx, userID)
	if err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return err
	}

	// パスワードを検証
	if !user.ComparePassword(password) {
		err := errors.New("パスワードが正しくありません")
		_ = i.outputPort.PresentError(ctx, err)
		return err
	}

	// ユーザーを削除
	if err := i.authService.DeleteUser(ctx, userID); err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return err
	}

	return nil
}
