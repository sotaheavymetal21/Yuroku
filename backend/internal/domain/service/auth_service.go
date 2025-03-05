package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yourusername/yuroku/internal/domain/entity"
	"github.com/yourusername/yuroku/internal/domain/repository"
)

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

// VerifyToken はJWTトークンを検証します
func (s *AuthService) VerifyToken(ctx context.Context, tokenString string) (string, error) {
	// トークンを解析
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 署名アルゴリズムを検証
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("無効な署名アルゴリズムです")
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return "", errors.New("トークンの検証に失敗しました")
	}

	// トークンの有効性を検証
	if !token.Valid {
		return "", errors.New("無効なトークンです")
	}

	// クレームを取得
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("無効なトークンクレームです")
	}

	// 有効期限を検証
	exp, ok := claims["exp"].(float64)
	if !ok {
		return "", errors.New("無効な有効期限です")
	}

	if time.Unix(int64(exp), 0).Before(time.Now()) {
		return "", errors.New("トークンの有効期限が切れています")
	}

	// ユーザーIDを取得
	userID, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("無効なユーザーIDです")
	}

	return userID, nil
}
