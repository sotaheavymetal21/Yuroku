package repository

import (
	"context"

	"github.com/yourusername/yuroku/internal/domain/entity"
)

// UserRepository はユーザー情報の永続化を担当するインターフェースです
type UserRepository interface {
	// Create は新しいユーザーを作成します
	Create(ctx context.Context, user *entity.User) error

	// FindByID はIDでユーザーを検索します
	FindByID(ctx context.Context, id string) (*entity.User, error)

	// FindByEmail はメールアドレスでユーザーを検索します
	FindByEmail(ctx context.Context, email string) (*entity.User, error)

	// Update はユーザー情報を更新します
	Update(ctx context.Context, user *entity.User) error

	// Delete はユーザーを削除します
	Delete(ctx context.Context, id string) error
}
