package presenter

import (
	"github.com/yourusername/yuroku/internal/domain/entity"
	"github.com/yourusername/yuroku/internal/usecase/port"
)

// OnsenImagePresenter は温泉画像関連のレスポンスを整形するプレゼンターです
type OnsenImagePresenter struct{}

// NewOnsenImagePresenter は新しいOnsenImagePresenterインスタンスを作成します
func NewOnsenImagePresenter() port.OnsenImagePresenterPort {
	return &OnsenImagePresenter{}
}

// PresentOnsenImage は単一の温泉画像レスポンスを整形します
func (p *OnsenImagePresenter) PresentOnsenImage(onsenImage *entity.OnsenImage) map[string]interface{} {
	return map[string]interface{}{
		"onsen_image": onsenImage,
	}
}

// PresentOnsenImages は複数の温泉画像レスポンスを整形します
func (p *OnsenImagePresenter) PresentOnsenImages(onsenImages []*entity.OnsenImage) map[string]interface{} {
	return map[string]interface{}{
		"onsen_images": onsenImages,
	}
}

// PresentError はエラーレスポンスを整形します
func (p *OnsenImagePresenter) PresentError(err error) map[string]interface{} {
	return map[string]interface{}{
		"error": err.Error(),
	}
}
