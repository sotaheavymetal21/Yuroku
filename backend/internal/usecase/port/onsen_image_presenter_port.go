package port

import (
	"github.com/yourusername/yuroku/internal/domain/entity"
)

// OnsenImagePresenterPort は温泉画像関連のレスポンスを整形するためのインターフェースです
type OnsenImagePresenterPort interface {
	// PresentOnsenImage は単一の温泉画像レスポンスを整形します
	PresentOnsenImage(onsenImage *entity.OnsenImage) map[string]interface{}

	// PresentOnsenImages は複数の温泉画像レスポンスを整形します
	PresentOnsenImages(onsenImages []*entity.OnsenImage) map[string]interface{}

	// PresentError はエラーレスポンスを整形します
	PresentError(err error) map[string]interface{}
}
