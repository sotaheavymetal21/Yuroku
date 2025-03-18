package presenter

import (
	"github.com/yourusername/yuroku/internal/domain/entity"
	"github.com/yourusername/yuroku/internal/usecase/port"
)

// OnsenLogPresenter は温泉メモ関連のレスポンスをフォーマットするプレゼンターです
type OnsenLogPresenter struct{}

// NewOnsenLogPresenter は新しいOnsenLogPresenterインスタンスを作成します
func NewOnsenLogPresenter() port.OnsenLogPresenterPort {
	return &OnsenLogPresenter{}
}

// PresentOnsenLog は温泉メモのレスポンスをフォーマットします
func (p *OnsenLogPresenter) PresentOnsenLog(onsenLog *entity.OnsenLog) map[string]interface{} {
	logData := formatOnsenLog(onsenLog)

	return map[string]interface{}{
		"data":    logData,
		"message": "温泉メモを取得しました",
	}
}

// PresentOnsenLogs は温泉メモリストのレスポンスをフォーマットします
func (p *OnsenLogPresenter) PresentOnsenLogs(onsenLogs []*entity.OnsenLog, total int64, page int, limit int) map[string]interface{} {
	logsData := make([]map[string]interface{}, len(onsenLogs))

	for i, log := range onsenLogs {
		logsData[i] = formatOnsenLog(log)
	}

	return map[string]interface{}{
		"data": map[string]interface{}{
			"onsen_logs": logsData,
			"total":      total,
			"page":       page,
			"limit":      limit,
		},
		"message": "温泉メモ一覧を取得しました",
	}
}

// PresentError はエラーレスポンスをフォーマットします
func (p *OnsenLogPresenter) PresentError(err error) map[string]interface{} {
	// エラーコードを決定（デフォルトはGENERAL_ERROR）
	errorCode := "GENERAL_ERROR"

	// エラーメッセージによってコードを判断
	errMsg := err.Error()
	switch {
	case contains(errMsg, "見つかりません"):
		errorCode = "NOT_FOUND"
	case contains(errMsg, "権限"):
		errorCode = "PERMISSION_DENIED"
	case contains(errMsg, "無効"):
		errorCode = "INVALID_INPUT"
	}

	return map[string]interface{}{
		"error": map[string]interface{}{
			"code":    errorCode,
			"message": errMsg,
		},
	}
}

// formatOnsenLog は温泉メモエンティティをレスポンス用のマップに変換します
func formatOnsenLog(onsenLog *entity.OnsenLog) map[string]interface{} {
	return map[string]interface{}{
		"id":          onsenLog.ID.Hex(),
		"uuid":        onsenLog.UUID,
		"user_id":     onsenLog.UserID,
		"name":        onsenLog.Name,
		"location":    onsenLog.Location,
		"spring_type": onsenLog.SpringType,
		"features":    onsenLog.Features,
		"visit_date":  onsenLog.VisitDate,
		"rating":      onsenLog.Rating,
		"comment":     onsenLog.Comment,
		"created_at":  onsenLog.CreatedAt,
		"updated_at":  onsenLog.UpdatedAt,
	}
}
