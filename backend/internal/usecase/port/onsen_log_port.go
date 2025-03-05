package port

import (
	"context"
	"time"

	"github.com/yourusername/yuroku/internal/domain/entity"
)

// OnsenLogInputPort は温泉メモユースケースの入力ポートです
type OnsenLogInputPort interface {
	// CreateOnsenLog は新しい温泉メモを作成します
	CreateOnsenLog(ctx context.Context, input CreateOnsenLogInput) (OnsenLogOutputData, error)

	// GetOnsenLog は温泉メモを取得します
	GetOnsenLog(ctx context.Context, id, userID string) (OnsenLogOutputData, error)

	// GetOnsenLogs はユーザーIDに紐づく温泉メモを取得します
	GetOnsenLogs(ctx context.Context, userID string, page, limit int) (OnsenLogsOutputData, error)

	// GetFilteredOnsenLogs はユーザーIDと条件に紐づく温泉メモを取得します
	GetFilteredOnsenLogs(ctx context.Context, input FilterOnsenLogsInput) (OnsenLogsOutputData, error)

	// UpdateOnsenLog は温泉メモを更新します
	UpdateOnsenLog(ctx context.Context, input UpdateOnsenLogInput) (OnsenLogOutputData, error)

	// DeleteOnsenLog は温泉メモを削除します
	DeleteOnsenLog(ctx context.Context, id, userID string) error

	// ExportOnsenLogs はユーザーIDに紐づく温泉メモをエクスポートします
	ExportOnsenLogs(ctx context.Context, userID string, format string) ([]byte, error)
}

// OnsenLogOutputPort は温泉メモユースケースの出力ポートです
type OnsenLogOutputPort interface {
	// PresentOnsenLog は温泉メモを表示します
	PresentOnsenLog(ctx context.Context, data OnsenLogOutputData) error

	// PresentOnsenLogs は温泉メモのリストを表示します
	PresentOnsenLogs(ctx context.Context, data OnsenLogsOutputData) error

	// PresentExportedData はエクスポートされたデータを表示します
	PresentExportedData(ctx context.Context, data []byte, format string) error

	// PresentError はエラーを表示します
	PresentError(ctx context.Context, err error) error
}

// CreateOnsenLogInput は温泉メモ作成の入力データです
type CreateOnsenLogInput struct {
	UserID     string            `json:"user_id"`
	Name       string            `json:"name"`
	Location   string            `json:"location"`
	SpringType entity.SpringType `json:"spring_type"`
	Features   []entity.Feature  `json:"features"`
	VisitDate  time.Time         `json:"visit_date"`
	Rating     int               `json:"rating"`
	Comment    string            `json:"comment"`
}

// UpdateOnsenLogInput は温泉メモ更新の入力データです
type UpdateOnsenLogInput struct {
	ID         string            `json:"id"`
	UserID     string            `json:"user_id"`
	Name       string            `json:"name"`
	Location   string            `json:"location"`
	SpringType entity.SpringType `json:"spring_type"`
	Features   []entity.Feature  `json:"features"`
	VisitDate  time.Time         `json:"visit_date"`
	Rating     int               `json:"rating"`
	Comment    string            `json:"comment"`
}

// FilterOnsenLogsInput は温泉メモフィルタリングの入力データです
type FilterOnsenLogsInput struct {
	UserID     string            `json:"user_id"`
	SpringType entity.SpringType `json:"spring_type"`
	Location   string            `json:"location"`
	MinRating  int               `json:"min_rating"`
	StartDate  *time.Time        `json:"start_date"`
	EndDate    *time.Time        `json:"end_date"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
}

// OnsenLogOutputData は温泉メモの出力データです
type OnsenLogOutputData struct {
	ID         string            `json:"id"`
	UserID     string            `json:"user_id"`
	Name       string            `json:"name"`
	Location   string            `json:"location"`
	SpringType entity.SpringType `json:"spring_type"`
	Features   []entity.Feature  `json:"features"`
	VisitDate  time.Time         `json:"visit_date"`
	Rating     int               `json:"rating"`
	Comment    string            `json:"comment"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
	Images     []ImageOutputData `json:"images,omitempty"`
}

// OnsenLogsOutputData は温泉メモリストの出力データです
type OnsenLogsOutputData struct {
	OnsenLogs  []OnsenLogOutputData `json:"onsen_logs"`
	TotalCount int                  `json:"total_count"`
	Page       int                  `json:"page"`
	Limit      int                  `json:"limit"`
}
