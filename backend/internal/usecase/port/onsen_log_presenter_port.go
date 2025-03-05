package port

import (
	"github.com/yourusername/yuroku/internal/domain/entity"
)

// OnsenLogPresenterPort は温泉ログ関連のレスポンスを整形するためのインターフェースです
type OnsenLogPresenterPort interface {
	// PresentOnsenLog は単一の温泉ログレスポンスを整形します
	PresentOnsenLog(onsenLog *entity.OnsenLog) map[string]interface{}

	// PresentOnsenLogs は複数の温泉ログレスポンスを整形します
	PresentOnsenLogs(onsenLogs []*entity.OnsenLog, total int64, page int, limit int) map[string]interface{}

	// PresentError はエラーレスポンスを整形します
	PresentError(err error) map[string]interface{}
}
