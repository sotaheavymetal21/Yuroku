package service

import (
	"context"
	"errors"
	"time"

	"github.com/yourusername/yuroku/internal/domain/entity"
	"github.com/yourusername/yuroku/internal/domain/repository"
)

// OnsenLogService は温泉メモに関するドメインサービスです
type OnsenLogService struct {
	onsenLogRepo repository.OnsenLogRepository
	imageRepo    repository.OnsenImageRepository
}

// NewOnsenLogService は新しい温泉メモサービスを作成します
func NewOnsenLogService(onsenLogRepo repository.OnsenLogRepository, imageRepo repository.OnsenImageRepository) *OnsenLogService {
	return &OnsenLogService{
		onsenLogRepo: onsenLogRepo,
		imageRepo:    imageRepo,
	}
}

// CreateOnsenLog は新しい温泉メモを作成します
func (s *OnsenLogService) CreateOnsenLog(ctx context.Context, userID, name, location string, springType entity.SpringType, features []entity.Feature, visitDate time.Time, rating int, comment string) (*entity.OnsenLog, error) {
	// 評価値のバリデーション
	if !entity.ValidateRating(rating) {
		return nil, errors.New("評価は0から5の間で指定してください")
	}

	// 温泉名のバリデーション
	if name == "" {
		return nil, errors.New("温泉名は必須です")
	}

	// 新しい温泉メモを作成
	onsenLog := entity.NewOnsenLog(userID, name, location, springType, features, visitDate, rating, comment)

	// 温泉メモを保存
	if err := s.onsenLogRepo.Create(ctx, onsenLog); err != nil {
		return nil, err
	}

	return onsenLog, nil
}

// GetOnsenLog は温泉メモを取得します
func (s *OnsenLogService) GetOnsenLog(ctx context.Context, id string) (*entity.OnsenLog, error) {
	return s.onsenLogRepo.FindByID(ctx, id)
}

// GetOnsenLogsByUserID はユーザーIDに紐づく温泉メモを取得します
func (s *OnsenLogService) GetOnsenLogsByUserID(ctx context.Context, userID string) ([]*entity.OnsenLog, error) {
	return s.onsenLogRepo.FindByUserID(ctx, userID)
}

// GetOnsenLogsByUserIDWithPagination はユーザーIDに紐づく温泉メモをページネーションで取得します
func (s *OnsenLogService) GetOnsenLogsByUserIDWithPagination(ctx context.Context, userID string, page, limit int) ([]*entity.OnsenLog, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return s.onsenLogRepo.FindByUserIDWithPagination(ctx, userID, page, limit)
}

// GetOnsenLogsByUserIDAndFilter はユーザーIDと条件に紐づく温泉メモを取得します
func (s *OnsenLogService) GetOnsenLogsByUserIDAndFilter(ctx context.Context, userID string, springType entity.SpringType, location string, minRating int, startDate, endDate *time.Time, page, limit int) ([]*entity.OnsenLog, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return s.onsenLogRepo.FindByUserIDAndFilter(ctx, userID, springType, location, minRating, startDate, endDate, page, limit)
}

// UpdateOnsenLog は温泉メモを更新します
func (s *OnsenLogService) UpdateOnsenLog(ctx context.Context, id, userID, name, location string, springType entity.SpringType, features []entity.Feature, visitDate time.Time, rating int, comment string) (*entity.OnsenLog, error) {
	// 評価値のバリデーション
	if !entity.ValidateRating(rating) {
		return nil, errors.New("評価は0から5の間で指定してください")
	}

	// 温泉名のバリデーション
	if name == "" {
		return nil, errors.New("温泉名は必須です")
	}

	// 温泉メモを取得
	onsenLog, err := s.onsenLogRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// ユーザーIDの検証
	if onsenLog.UserID != userID {
		return nil, errors.New("この温泉メモを編集する権限がありません")
	}

	// 温泉メモを更新
	onsenLog.Update(name, location, springType, features, visitDate, rating, comment)

	// 更新を保存
	if err := s.onsenLogRepo.Update(ctx, onsenLog); err != nil {
		return nil, err
	}

	return onsenLog, nil
}

// DeleteOnsenLog は温泉メモを削除します
func (s *OnsenLogService) DeleteOnsenLog(ctx context.Context, id, userID string) error {
	// 温泉メモを取得
	onsenLog, err := s.onsenLogRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// ユーザーIDの検証
	if onsenLog.UserID != userID {
		return errors.New("この温泉メモを削除する権限がありません")
	}

	// 関連する画像を削除
	if err := s.imageRepo.DeleteByOnsenID(ctx, id); err != nil {
		return err
	}

	// 温泉メモを削除
	return s.onsenLogRepo.Delete(ctx, id)
}
