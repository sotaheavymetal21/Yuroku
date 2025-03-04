package interactor

import (
	"context"
	"errors"

	"github.com/yourusername/yuroku/internal/domain/service"
	"github.com/yourusername/yuroku/internal/usecase/port"
)

// OnsenImageInteractor は温泉画像ユースケースのインタラクターです
type OnsenImageInteractor struct {
	onsenImageService *service.OnsenImageService
	outputPort        port.OnsenImageOutputPort
}

// NewOnsenImageInteractor は新しい温泉画像インタラクターを作成します
func NewOnsenImageInteractor(
	onsenImageService *service.OnsenImageService,
	outputPort port.OnsenImageOutputPort,
) *OnsenImageInteractor {
	return &OnsenImageInteractor{
		onsenImageService: onsenImageService,
		outputPort:        outputPort,
	}
}

// UploadImage は温泉画像をアップロードします
func (i *OnsenImageInteractor) UploadImage(ctx context.Context, input port.UploadImageInput) (port.ImageOutputData, error) {
	// 入力値のバリデーション
	if input.OnsenID == "" {
		err := errors.New("温泉IDは必須です")
		_ = i.outputPort.PresentError(ctx, err)
		return port.ImageOutputData{}, err
	}

	if input.File == nil {
		err := errors.New("ファイルは必須です")
		_ = i.outputPort.PresentError(ctx, err)
		return port.ImageOutputData{}, err
	}

	// ドメインサービスを呼び出し
	image, err := i.onsenImageService.UploadImage(
		ctx,
		input.OnsenID,
		input.UserID,
		input.File,
		input.FileName,
		input.ContentType,
	)
	if err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return port.ImageOutputData{}, err
	}

	// 出力データを作成
	outputData := port.ImageOutputData{
		ID:        image.UUID,
		OnsenID:   image.OnsenID,
		ImageURL:  image.ImageURL,
		CreatedAt: image.CreatedAt,
	}

	// 出力ポートを呼び出し
	if err := i.outputPort.PresentImage(ctx, outputData); err != nil {
		return port.ImageOutputData{}, err
	}

	return outputData, nil
}

// GetImagesByOnsenID は温泉IDに紐づく画像を取得します
func (i *OnsenImageInteractor) GetImagesByOnsenID(ctx context.Context, onsenID, userID string) ([]port.ImageOutputData, error) {
	// 入力値のバリデーション
	if onsenID == "" {
		err := errors.New("温泉IDは必須です")
		_ = i.outputPort.PresentError(ctx, err)
		return nil, err
	}

	// ドメインサービスを呼び出し
	images, err := i.onsenImageService.GetImagesByOnsenID(ctx, onsenID, userID)
	if err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return nil, err
	}

	// 出力データを作成
	outputData := make([]port.ImageOutputData, len(images))
	for i, image := range images {
		outputData[i] = port.ImageOutputData{
			ID:        image.UUID,
			OnsenID:   image.OnsenID,
			ImageURL:  image.ImageURL,
			CreatedAt: image.CreatedAt,
		}
	}

	// 出力ポートを呼び出し
	if err := i.outputPort.PresentImages(ctx, outputData); err != nil {
		return nil, err
	}

	return outputData, nil
}

// DeleteImage は温泉画像を削除します
func (i *OnsenImageInteractor) DeleteImage(ctx context.Context, imageID, userID string) error {
	// 入力値のバリデーション
	if imageID == "" {
		err := errors.New("画像IDは必須です")
		_ = i.outputPort.PresentError(ctx, err)
		return err
	}

	// ドメインサービスを呼び出し
	if err := i.onsenImageService.DeleteImage(ctx, imageID, userID); err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return err
	}

	return nil
}
