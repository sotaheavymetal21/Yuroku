package presenter

import (
	"github.com/yourusername/yuroku/internal/domain/entity"
	"github.com/yourusername/yuroku/internal/usecase/port"
)

// OnsenLogPresenter は温泉ログ関連のレスポンスを整形するプレゼンターです
type OnsenLogPresenter struct{}

// NewOnsenLogPresenter は新しいOnsenLogPresenterインスタンスを作成します
func NewOnsenLogPresenter() port.OnsenLogPresenterPort {
	return &OnsenLogPresenter{}
}

// PresentOnsenLog は単一の温泉ログレスポンスを整形します
func (p *OnsenLogPresenter) PresentOnsenLog(onsenLog *entity.OnsenLog) map[string]interface{} {
	return map[string]interface{}{
		"onsen_log": onsenLog,
	}
}

// PresentOnsenLogs は複数の温泉ログレスポンスを整形します
func (p *OnsenLogPresenter) PresentOnsenLogs(onsenLogs []*entity.OnsenLog, total int64, page int, limit int) map[string]interface{} {
	return map[string]interface{}{
		"onsen_logs": onsenLogs,
		"pagination": map[string]interface{}{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	}
}

// PresentError はエラーレスポンスを整形します
func (p *OnsenLogPresenter) PresentError(err error) map[string]interface{} {
	return map[string]interface{}{
		"error": err.Error(),
	}
}
