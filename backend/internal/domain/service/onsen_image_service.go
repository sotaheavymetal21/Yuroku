package service

import (
	"context"
	"errors"
	"io"
	"path/filepath"
	"strings"

	"github.com/yourusername/yuroku/internal/domain/entity"
	"github.com/yourusername/yuroku/internal/domain/repository"
)

// OnsenImageService は温泉画像に関するドメインサービスです
type OnsenImageService struct {
	imageRepo    repository.OnsenImageRepository
	onsenLogRepo repository.OnsenLogRepository
	storageRepo  repository.StorageRepository
}

// NewOnsenImageService は新しい温泉画像サービスを作成します
func NewOnsenImageService(imageRepo repository.OnsenImageRepository, onsenLogRepo repository.OnsenLogRepository, storageRepo repository.StorageRepository) *OnsenImageService {
	return &OnsenImageService{
		imageRepo:    imageRepo,
		onsenLogRepo: onsenLogRepo,
		storageRepo:  storageRepo,
	}
}

// UploadImage は温泉画像をアップロードします
func (s *OnsenImageService) UploadImage(ctx context.Context, onsenID, userID string, file io.Reader, fileName, contentType, description string) (*entity.OnsenImage, error) {
	// 温泉メモを取得
	onsenLog, err := s.onsenLogRepo.FindByID(ctx, onsenID)
	if err != nil {
		return nil, err
	}

	// ユーザーIDの検証
	if onsenLog.UserID != userID {
		return nil, errors.New("この温泉メモに画像をアップロードする権限がありません")
	}

	// 既存の画像数をチェック
	existingImages, err := s.imageRepo.FindByOnsenID(ctx, onsenID)
	if err != nil {
		return nil, err
	}

	// 最大3枚までの制限
	if len(existingImages) >= 3 {
		return nil, errors.New("画像は最大3枚までアップロードできます")
	}

	// ファイル名を生成（UUID + 元の拡張子）
	ext := filepath.Ext(fileName)
	if ext == "" {
		// Content-Typeからファイル拡張子を推測
		ext = getExtensionFromContentType(contentType)
	}

	// ファイルをストレージにアップロード
	fileURL, err := s.storageRepo.Upload(ctx, file, onsenID+"-"+fileName, contentType)
	if err != nil {
		return nil, err
	}

	// 画像情報をデータベースに保存
	onsenImage := entity.NewOnsenImage(onsenID, userID, fileURL, description)
	if err := s.imageRepo.Create(ctx, onsenImage); err != nil {
		// エラーが発生した場合、アップロードしたファイルを削除
		_ = s.storageRepo.Delete(ctx, fileURL)
		return nil, err
	}

	return onsenImage, nil
}

// GetImagesByOnsenID は温泉IDに紐づく画像を取得します
func (s *OnsenImageService) GetImagesByOnsenID(ctx context.Context, onsenID, userID string) ([]*entity.OnsenImage, error) {
	// 温泉メモを取得
	onsenLog, err := s.onsenLogRepo.FindByID(ctx, onsenID)
	if err != nil {
		return nil, err
	}

	// ユーザーIDの検証
	if onsenLog.UserID != userID {
		return nil, errors.New("この温泉メモの画像を閲覧する権限がありません")
	}

	// 画像を取得
	return s.imageRepo.FindByOnsenID(ctx, onsenID)
}

// DeleteImage は温泉画像を削除します
func (s *OnsenImageService) DeleteImage(ctx context.Context, imageID, userID string) error {
	// 画像を取得
	image, err := s.imageRepo.FindByID(ctx, imageID)
	if err != nil {
		return err
	}

	// 温泉メモを取得
	onsenLog, err := s.onsenLogRepo.FindByID(ctx, image.OnsenID)
	if err != nil {
		return err
	}

	// ユーザーIDの検証
	if onsenLog.UserID != userID {
		return errors.New("この画像を削除する権限がありません")
	}

	// ストレージから画像を削除
	if err := s.storageRepo.Delete(ctx, image.ImageURL); err != nil {
		return err
	}

	// データベースから画像情報を削除
	return s.imageRepo.Delete(ctx, imageID)
}

// getExtensionFromContentType はContent-Typeからファイル拡張子を推測します
func getExtensionFromContentType(contentType string) string {
	switch {
	case strings.Contains(contentType, "image/jpeg"):
		return ".jpg"
	case strings.Contains(contentType, "image/png"):
		return ".png"
	case strings.Contains(contentType, "image/gif"):
		return ".gif"
	case strings.Contains(contentType, "image/webp"):
		return ".webp"
	default:
		return ".bin"
	}
}
