package repository

import (
	"context"

	"github.com/yourusername/yuroku/internal/domain/entity"
)

// OnsenImageRepository は温泉画像の永続化を担当するインターフェースです
type OnsenImageRepository interface {
	// Create は新しい温泉画像を作成します
	Create(ctx context.Context, onsenImage *entity.OnsenImage) error

	// FindByID はIDで温泉画像を検索します
	FindByID(ctx context.Context, id string) (*entity.OnsenImage, error)

	// FindByOnsenID は温泉IDに紐づく画像を検索します
	FindByOnsenID(ctx context.Context, onsenID string) ([]*entity.OnsenImage, error)

	// Delete は温泉画像を削除します
	Delete(ctx context.Context, id string) error

	// DeleteByOnsenID は温泉IDに紐づく画像をすべて削除します
	DeleteByOnsenID(ctx context.Context, onsenID string) error
}
