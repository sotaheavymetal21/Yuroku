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
func (i *AuthInteractor) Register(ctx context.Context, email, password string) (port.AuthOutputData, error) {
	// 入力値のバリデーション
	if email == "" || password == "" {
		err := errors.New("メールアドレスとパスワードは必須です")
		_ = i.outputPort.PresentError(ctx, err)
		return port.AuthOutputData{}, err
	}

	// ドメインサービスを呼び出し
	user, err := i.authService.Register(ctx, email, password)
	if err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return port.AuthOutputData{}, err
	}

	// JWTトークンを生成
	accessToken, err := i.generateAccessToken(user.UUID)
	if err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return port.AuthOutputData{}, err
	}

	refreshToken, err := i.generateRefreshToken(user.UUID)
	if err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return port.AuthOutputData{}, err
	}

	// 出力データを作成
	outputData := port.AuthOutputData{
		UserID:       user.UUID,
		Email:        user.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	// 出力ポートを呼び出し
	if err := i.outputPort.PresentRegister(ctx, outputData); err != nil {
		return port.AuthOutputData{}, err
	}

	return outputData, nil
}

// Login はユーザーログインを行います
func (i *AuthInteractor) Login(ctx context.Context, email, password string) (port.AuthOutputData, error) {
	// 入力値のバリデーション
	if email == "" || password == "" {
		err := errors.New("メールアドレスとパスワードは必須です")
		_ = i.outputPort.PresentError(ctx, err)
		return port.AuthOutputData{}, err
	}

	// ドメインサービスを呼び出し
	user, err := i.authService.Login(ctx, email, password)
	if err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return port.AuthOutputData{}, err
	}

	// JWTトークンを生成
	accessToken, err := i.generateAccessToken(user.UUID)
	if err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return port.AuthOutputData{}, err
	}

	refreshToken, err := i.generateRefreshToken(user.UUID)
	if err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return port.AuthOutputData{}, err
	}

	// 出力データを作成
	outputData := port.AuthOutputData{
		UserID:       user.UUID,
		Email:        user.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	// 出力ポートを呼び出し
	if err := i.outputPort.PresentLogin(ctx, outputData); err != nil {
		return port.AuthOutputData{}, err
	}

	return outputData, nil
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
func (i *AuthInteractor) RefreshToken(ctx context.Context, refreshToken string) (port.AuthOutputData, error) {
	// リフレッシュトークンを検証
	userID, err := i.VerifyToken(ctx, refreshToken)
	if err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return port.AuthOutputData{}, err
	}

	// ユーザー情報を取得
	user, err := i.authService.GetUserByID(ctx, userID)
	if err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return port.AuthOutputData{}, err
	}

	// 新しいトークンを生成
	accessToken, err := i.generateAccessToken(user.UUID)
	if err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return port.AuthOutputData{}, err
	}

	newRefreshToken, err := i.generateRefreshToken(user.UUID)
	if err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return port.AuthOutputData{}, err
	}

	// 出力データを作成
	outputData := port.AuthOutputData{
		UserID:       user.UUID,
		Email:        user.Email,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}

	// 出力ポートを呼び出し
	if err := i.outputPort.PresentRefreshToken(ctx, outputData); err != nil {
		return port.AuthOutputData{}, err
	}

	return outputData, nil
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
