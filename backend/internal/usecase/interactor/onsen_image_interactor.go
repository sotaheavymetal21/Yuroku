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
	if input.OnsenID == "" || input.UserID == "" || input.File == nil {
		err := errors.New("温泉ID、ユーザーID、ファイルは必須です")
		_ = i.outputPort.PresentError(ctx, err)
		return port.ImageOutputData{}, err
	}

	// ドメインサービスを呼び出し
	onsenImage, err := i.onsenImageService.UploadImage(
		ctx,
		input.OnsenID,
		input.UserID,
		input.File,
		input.Filename,
		input.ContentType,
		input.Description,
	)
	if err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return port.ImageOutputData{}, err
	}

	// 出力データを作成
	outputData := port.ImageOutputData{
		ID:          onsenImage.UUID,
		OnsenID:     onsenImage.OnsenID,
		URL:         onsenImage.ImageURL,
		Description: onsenImage.Description,
		CreatedAt:   onsenImage.CreatedAt,
	}

	// 出力ポートを呼び出し
	if err := i.outputPort.PresentImage(ctx, outputData); err != nil {
		return port.ImageOutputData{}, err
	}

	return outputData, nil
}

// GetImagesByOnsenID は温泉IDに紐づく画像を取得します
func (i *OnsenImageInteractor) GetImagesByOnsenID(ctx context.Context, input port.GetImagesByOnsenIDInput) ([]port.ImageOutputData, error) {
	// 入力値のバリデーション
	if input.OnsenID == "" || input.UserID == "" {
		err := errors.New("温泉IDとユーザーIDは必須です")
		_ = i.outputPort.PresentError(ctx, err)
		return nil, err
	}

	// ドメインサービスを呼び出し
	images, err := i.onsenImageService.GetImagesByOnsenID(ctx, input.OnsenID, input.UserID)
	if err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return nil, err
	}

	// 出力データを作成
	outputData := make([]port.ImageOutputData, len(images))
	for i, image := range images {
		outputData[i] = port.ImageOutputData{
			ID:          image.UUID,
			OnsenID:     image.OnsenID,
			URL:         image.ImageURL,
			Description: image.Description,
			CreatedAt:   image.CreatedAt,
		}
	}

	// 出力ポートを呼び出し
	if err := i.outputPort.PresentImages(ctx, outputData); err != nil {
		return nil, err
	}

	return outputData, nil
}

// DeleteImage は温泉画像を削除します
func (i *OnsenImageInteractor) DeleteImage(ctx context.Context, input port.DeleteImageInput) error {
	// 入力値のバリデーション
	if input.ImageID == "" || input.UserID == "" {
		err := errors.New("画像IDとユーザーIDは必須です")
		_ = i.outputPort.PresentError(ctx, err)
		return err
	}

	// ドメインサービスを呼び出し
	err := i.onsenImageService.DeleteImage(ctx, input.ImageID, input.UserID)
	if err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return err
	}

	return nil
}
