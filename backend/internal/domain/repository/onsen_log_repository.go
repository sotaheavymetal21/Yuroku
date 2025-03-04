package repository

import (
	"context"
	"time"

	"github.com/yourusername/yuroku/internal/domain/entity"
)

// OnsenLogRepository は温泉メモの永続化を担当するインターフェースです
type OnsenLogRepository interface {
	// Create は新しい温泉メモを作成します
	Create(ctx context.Context, onsenLog *entity.OnsenLog) error

	// FindByID はIDで温泉メモを検索します
	FindByID(ctx context.Context, id string) (*entity.OnsenLog, error)

	// FindByUserID はユーザーIDに紐づく温泉メモを検索します
	FindByUserID(ctx context.Context, userID string) ([]*entity.OnsenLog, error)

	// FindByUserIDWithPagination はユーザーIDに紐づく温泉メモをページネーションで検索します
	FindByUserIDWithPagination(ctx context.Context, userID string, page, limit int) ([]*entity.OnsenLog, int, error)

	// FindByUserIDAndFilter はユーザーIDと条件に紐づく温泉メモを検索します
	FindByUserIDAndFilter(ctx context.Context, userID string, springType entity.SpringType, location string, minRating int, startDate, endDate *time.Time, page, limit int) ([]*entity.OnsenLog, int, error)

	// Update は温泉メモを更新します
	Update(ctx context.Context, onsenLog *entity.OnsenLog) error

	// Delete は温泉メモを削除します
	Delete(ctx context.Context, id string) error

	// DeleteByUserID はユーザーIDに紐づく温泉メモをすべて削除します
	DeleteByUserID(ctx context.Context, userID string) error
}
