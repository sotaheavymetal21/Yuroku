package service

import (
	"context"
	"errors"

	"github.com/yourusername/yuroku/internal/domain/entity"
	"github.com/yourusername/yuroku/internal/domain/repository"
)

// AuthService は認証に関するドメインサービスです
type AuthService struct {
	userRepo repository.UserRepository
}

// NewAuthService は新しい認証サービスを作成します
func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// Register はユーザー登録を行います
func (s *AuthService) Register(ctx context.Context, email, password string) (*entity.User, error) {
	// メールアドレスの重複チェック
	existingUser, err := s.userRepo.FindByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, errors.New("このメールアドレスは既に登録されています")
	}

	// 新しいユーザーを作成
	user, err := entity.NewUser(email, password)
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
