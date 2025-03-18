package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yourusername/yuroku/internal/domain/entity"
	"github.com/yourusername/yuroku/internal/domain/repository"
)

// TokenType はトークンの種類を表します
type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

// TokenClaims はJWTトークンのクレームを表します
type TokenClaims struct {
	UserID string    `json:"sub"`
	Type   TokenType `json:"type"`
	jwt.RegisteredClaims
}

// AuthService は認証に関するドメインサービスです
type AuthService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

// NewAuthService は新しい認証サービスを作成します
func NewAuthService(userRepo repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// Register はユーザー登録を行います
func (s *AuthService) Register(ctx context.Context, name, email, password string) (*entity.User, error) {
	// メールアドレスの重複チェック
	existingUser, err := s.userRepo.FindByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, errors.New("このメールアドレスは既に登録されています")
	}

	// 新しいユーザーを作成
	user, err := entity.NewUser(name, email, password)
	if err != nil {
		return nil, err
	}

	// ユーザーを保存
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login はユーザーログインを行います
func (s *AuthService) Login(ctx context.Context, email, password string) (*entity.User, error) {
	// メールアドレスでユーザーを検索
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("メールアドレスまたはパスワードが正しくありません")
	}

	// パスワードを検証
	if !user.ComparePassword(password) {
		return nil, errors.New("メールアドレスまたはパスワードが正しくありません")
	}

	return user, nil
}

// GetUserByID はIDでユーザーを取得します
func (s *AuthService) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

// UpdateUser はユーザー情報を更新します
func (s *AuthService) UpdateUser(ctx context.Context, user *entity.User) error {
	return s.userRepo.Update(ctx, user)
}

// DeleteUser はユーザーを削除します
func (s *AuthService) DeleteUser(ctx context.Context, id string) error {
	return s.userRepo.Delete(ctx, id)
}

// GenerateTokens はアクセストークンとリフレッシュトークンを生成します
func (s *AuthService) GenerateTokens(userID string, accessExp, refreshExp time.Duration) (string, string, error) {
	// アクセストークンを生成
	accessToken, err := s.generateToken(userID, AccessToken, accessExp)
	if err != nil {
		return "", "", fmt.Errorf("アクセストークンの生成に失敗しました: %w", err)
	}

	// リフレッシュトークンを生成
	refreshToken, err := s.generateToken(userID, RefreshToken, refreshExp)
	if err != nil {
		return "", "", fmt.Errorf("リフレッシュトークンの生成に失敗しました: %w", err)
	}

	return accessToken, refreshToken, nil
}

// generateToken は指定した種類のトークンを生成します
func (s *AuthService) generateToken(userID string, tokenType TokenType, expiry time.Duration) (string, error) {
	now := time.Now()
	claims := TokenClaims{
		UserID: userID,
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "yuroku-api",
			Subject:   userID,
		},
	}

	// トークンを生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// トークンに署名
	signedToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// VerifyToken はJWTトークンを検証してユーザーIDを返します
func (s *AuthService) VerifyToken(ctx context.Context, tokenString string) (string, error) {
	// トークンを解析
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 署名アルゴリズムを検証
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("無効な署名アルゴリズムです")
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return "", fmt.Errorf("トークンの検証に失敗しました: %w", err)
	}

	// トークンの有効性を検証
	if !token.Valid {
		return "", errors.New("無効なトークンです")
	}

	// クレームを取得
	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return "", errors.New("無効なトークンクレームです")
	}

	// アクセストークンとリフレッシュトークンの両方を許可する
	// トークン種別のチェックを削除

	// ユーザーIDを返す
	return claims.UserID, nil
}

// VerifyRefreshToken はリフレッシュトークンを検証してユーザーIDを返します
func (s *AuthService) VerifyRefreshToken(ctx context.Context, tokenString string) (string, error) {
	// トークンを解析
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 署名アルゴリズムを検証
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("無効な署名アルゴリズムです")
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return "", fmt.Errorf("リフレッシュトークンの検証に失敗しました: %w", err)
	}

	// トークンの有効性を検証
	if !token.Valid {
		return "", errors.New("無効なリフレッシュトークンです")
	}

	// クレームを取得
	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return "", errors.New("無効なトークンクレームです")
	}

	// トークン種別を確認
	if claims.Type != RefreshToken {
		return "", errors.New("無効なリフレッシュトークンです")
	}

	// ユーザーの存在確認
	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil || user == nil {
		return "", errors.New("ユーザーが見つかりません")
	}

	// ユーザーIDを返す
	return claims.UserID, nil
}
