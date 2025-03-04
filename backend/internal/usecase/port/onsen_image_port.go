package port

import (
	"context"
	"io"
)

// OnsenImageInputPort は温泉画像ユースケースの入力ポートです
type OnsenImageInputPort interface {
	// UploadImage は温泉画像をアップロードします
	UploadImage(ctx context.Context, input UploadImageInput) (ImageOutputData, error)

	// GetImagesByOnsenID は温泉IDに紐づく画像を取得します
	GetImagesByOnsenID(ctx context.Context, onsenID, userID string) ([]ImageOutputData, error)

	// DeleteImage は温泉画像を削除します
	DeleteImage(ctx context.Context, imageID, userID string) error
}

// OnsenImageOutputPort は温泉画像ユースケースの出力ポートです
type OnsenImageOutputPort interface {
	// PresentImage は温泉画像を表示します
	PresentImage(ctx context.Context, data ImageOutputData) error

	// PresentImages は温泉画像のリストを表示します
	PresentImages(ctx context.Context, data []ImageOutputData) error

	// PresentError はエラーを表示します
	PresentError(ctx context.Context, err error) error
}

// UploadImageInput は画像アップロードの入力データです
type UploadImageInput struct {
	OnsenID     string    `json:"onsen_id"`
	UserID      string    `json:"user_id"`
	File        io.Reader `json:"-"`
	FileName    string    `json:"file_name"`
	ContentType string    `json:"content_type"`
}
