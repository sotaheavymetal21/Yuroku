package interactor

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/yourusername/yuroku/internal/domain/entity"
	"github.com/yourusername/yuroku/internal/domain/service"
	"github.com/yourusername/yuroku/internal/usecase/port"
)

// OnsenLogInteractor は温泉メモユースケースのインタラクターです
type OnsenLogInteractor struct {
	onsenLogService   *service.OnsenLogService
	onsenImageService *service.OnsenImageService
	outputPort        port.OnsenLogOutputPort
}

// NewOnsenLogInteractor は新しい温泉メモインタラクターを作成します
func NewOnsenLogInteractor(
	onsenLogService *service.OnsenLogService,
	onsenImageService *service.OnsenImageService,
	outputPort port.OnsenLogOutputPort,
) *OnsenLogInteractor {
	return &OnsenLogInteractor{
		onsenLogService:   onsenLogService,
		onsenImageService: onsenImageService,
		outputPort:        outputPort,
	}
}

// CreateOnsenLog は新しい温泉メモを作成します
func (i *OnsenLogInteractor) CreateOnsenLog(ctx context.Context, input port.CreateOnsenLogInput) (port.OnsenLogOutputData, error) {
	// 入力値のバリデーション
	if input.Name == "" {
		err := errors.New("温泉名は必須です")
		_ = i.outputPort.PresentError(ctx, err)
		return port.OnsenLogOutputData{}, err
	}

	// ドメインサービスを呼び出し
	onsenLog, err := i.onsenLogService.CreateOnsenLog(
		ctx,
		input.UserID,
		input.Name,
		input.Location,
		input.SpringType,
		input.Features,
		input.VisitDate,
		input.Rating,
		input.Comment,
	)
	if err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return port.OnsenLogOutputData{}, err
	}

	// 出力データを作成
	outputData := port.OnsenLogOutputData{
		ID:         onsenLog.UUID,
		UserID:     onsenLog.UserID,
		Name:       onsenLog.Name,
		Location:   onsenLog.Location,
		SpringType: onsenLog.SpringType,
		Features:   onsenLog.Features,
		VisitDate:  onsenLog.VisitDate,
		Rating:     onsenLog.Rating,
		Comment:    onsenLog.Comment,
		CreatedAt:  onsenLog.CreatedAt,
		UpdatedAt:  onsenLog.UpdatedAt,
	}

	// 出力ポートを呼び出し
	if err := i.outputPort.PresentOnsenLog(ctx, outputData); err != nil {
		return port.OnsenLogOutputData{}, err
	}

	return outputData, nil
}

// GetOnsenLog は温泉メモを取得します
func (i *OnsenLogInteractor) GetOnsenLog(ctx context.Context, id, userID string) (port.OnsenLogOutputData, error) {
	// ドメインサービスを呼び出し
	onsenLog, err := i.onsenLogService.GetOnsenLog(ctx, id)
	if err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return port.OnsenLogOutputData{}, err
	}

	// ユーザーIDの検証
	if onsenLog.UserID != userID {
		err := errors.New("この温泉メモを閲覧する権限がありません")
		_ = i.outputPort.PresentError(ctx, err)
		return port.OnsenLogOutputData{}, err
	}

	// 画像を取得
	images, err := i.onsenImageService.GetImagesByOnsenID(ctx, id, userID)
	if err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return port.OnsenLogOutputData{}, err
	}

	// 画像の出力データを作成
	imageOutputData := make([]port.ImageOutputData, len(images))
	for i, image := range images {
		imageOutputData[i] = port.ImageOutputData{
			ID:          image.UUID,
			OnsenID:     image.OnsenID,
			URL:         image.ImageURL,
			Description: image.Description,
			CreatedAt:   image.CreatedAt,
		}
	}

	// 出力データを作成
	outputData := port.OnsenLogOutputData{
		ID:         onsenLog.UUID,
		UserID:     onsenLog.UserID,
		Name:       onsenLog.Name,
		Location:   onsenLog.Location,
		SpringType: onsenLog.SpringType,
		Features:   onsenLog.Features,
		VisitDate:  onsenLog.VisitDate,
		Rating:     onsenLog.Rating,
		Comment:    onsenLog.Comment,
		CreatedAt:  onsenLog.CreatedAt,
		UpdatedAt:  onsenLog.UpdatedAt,
		Images:     imageOutputData,
	}

	// 出力ポートを呼び出し
	if err := i.outputPort.PresentOnsenLog(ctx, outputData); err != nil {
		return port.OnsenLogOutputData{}, err
	}

	return outputData, nil
}

// GetOnsenLogs はユーザーIDに紐づく温泉メモを取得します
func (i *OnsenLogInteractor) GetOnsenLogs(ctx context.Context, userID string, page, limit int) (port.OnsenLogsOutputData, error) {
	// ドメインサービスを呼び出し
	onsenLogs, totalCount, err := i.onsenLogService.GetOnsenLogsByUserIDWithPagination(ctx, userID, page, limit)
	if err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return port.OnsenLogsOutputData{}, err
	}

	// 出力データを作成
	onsenLogOutputData := make([]port.OnsenLogOutputData, len(onsenLogs))
	for i, onsenLog := range onsenLogs {
		onsenLogOutputData[i] = port.OnsenLogOutputData{
			ID:         onsenLog.UUID,
			UserID:     onsenLog.UserID,
			Name:       onsenLog.Name,
			Location:   onsenLog.Location,
			SpringType: onsenLog.SpringType,
			Features:   onsenLog.Features,
			VisitDate:  onsenLog.VisitDate,
			Rating:     onsenLog.Rating,
			Comment:    onsenLog.Comment,
			CreatedAt:  onsenLog.CreatedAt,
			UpdatedAt:  onsenLog.UpdatedAt,
		}
	}

	outputData := port.OnsenLogsOutputData{
		OnsenLogs:  onsenLogOutputData,
		TotalCount: totalCount,
		Page:       page,
		Limit:      limit,
	}

	// 出力ポートを呼び出し
	if err := i.outputPort.PresentOnsenLogs(ctx, outputData); err != nil {
		return port.OnsenLogsOutputData{}, err
	}

	return outputData, nil
}

// GetFilteredOnsenLogs はユーザーIDと条件に紐づく温泉メモを取得します
func (i *OnsenLogInteractor) GetFilteredOnsenLogs(ctx context.Context, input port.FilterOnsenLogsInput) (port.OnsenLogsOutputData, error) {
	// ドメインサービスを呼び出し
	onsenLogs, totalCount, err := i.onsenLogService.GetOnsenLogsByUserIDAndFilter(
		ctx,
		input.UserID,
		input.SpringType,
		input.Location,
		input.MinRating,
		input.StartDate,
		input.EndDate,
		input.Page,
		input.Limit,
	)
	if err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return port.OnsenLogsOutputData{}, err
	}

	// 出力データを作成
	onsenLogOutputData := make([]port.OnsenLogOutputData, len(onsenLogs))
	for i, onsenLog := range onsenLogs {
		onsenLogOutputData[i] = port.OnsenLogOutputData{
			ID:         onsenLog.UUID,
			UserID:     onsenLog.UserID,
			Name:       onsenLog.Name,
			Location:   onsenLog.Location,
			SpringType: onsenLog.SpringType,
			Features:   onsenLog.Features,
			VisitDate:  onsenLog.VisitDate,
			Rating:     onsenLog.Rating,
			Comment:    onsenLog.Comment,
			CreatedAt:  onsenLog.CreatedAt,
			UpdatedAt:  onsenLog.UpdatedAt,
		}
	}

	outputData := port.OnsenLogsOutputData{
		OnsenLogs:  onsenLogOutputData,
		TotalCount: totalCount,
		Page:       input.Page,
		Limit:      input.Limit,
	}

	// 出力ポートを呼び出し
	if err := i.outputPort.PresentOnsenLogs(ctx, outputData); err != nil {
		return port.OnsenLogsOutputData{}, err
	}

	return outputData, nil
}

// UpdateOnsenLog は温泉メモを更新します
func (i *OnsenLogInteractor) UpdateOnsenLog(ctx context.Context, input port.UpdateOnsenLogInput) (port.OnsenLogOutputData, error) {
	// 入力値のバリデーション
	if input.Name == "" {
		err := errors.New("温泉名は必須です")
		_ = i.outputPort.PresentError(ctx, err)
		return port.OnsenLogOutputData{}, err
	}

	// ドメインサービスを呼び出し
	onsenLog, err := i.onsenLogService.UpdateOnsenLog(
		ctx,
		input.ID,
		input.UserID,
		input.Name,
		input.Location,
		input.SpringType,
		input.Features,
		input.VisitDate,
		input.Rating,
		input.Comment,
	)
	if err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return port.OnsenLogOutputData{}, err
	}

	// 画像を取得
	images, err := i.onsenImageService.GetImagesByOnsenID(ctx, input.ID, input.UserID)
	if err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return port.OnsenLogOutputData{}, err
	}

	// 画像の出力データを作成
	imageOutputData := make([]port.ImageOutputData, len(images))
	for i, image := range images {
		imageOutputData[i] = port.ImageOutputData{
			ID:          image.UUID,
			OnsenID:     image.OnsenID,
			URL:         image.ImageURL,
			Description: image.Description,
			CreatedAt:   image.CreatedAt,
		}
	}

	// 出力データを作成
	outputData := port.OnsenLogOutputData{
		ID:         onsenLog.UUID,
		UserID:     onsenLog.UserID,
		Name:       onsenLog.Name,
		Location:   onsenLog.Location,
		SpringType: onsenLog.SpringType,
		Features:   onsenLog.Features,
		VisitDate:  onsenLog.VisitDate,
		Rating:     onsenLog.Rating,
		Comment:    onsenLog.Comment,
		CreatedAt:  onsenLog.CreatedAt,
		UpdatedAt:  onsenLog.UpdatedAt,
		Images:     imageOutputData,
	}

	// 出力ポートを呼び出し
	if err := i.outputPort.PresentOnsenLog(ctx, outputData); err != nil {
		return port.OnsenLogOutputData{}, err
	}

	return outputData, nil
}

// DeleteOnsenLog は温泉メモを削除します
func (i *OnsenLogInteractor) DeleteOnsenLog(ctx context.Context, id, userID string) error {
	// ドメインサービスを呼び出し
	if err := i.onsenLogService.DeleteOnsenLog(ctx, id, userID); err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return err
	}

	return nil
}

// ExportOnsenLogs はユーザーIDに紐づく温泉メモをエクスポートします
func (i *OnsenLogInteractor) ExportOnsenLogs(ctx context.Context, userID string, format string) ([]byte, error) {
	// ドメインサービスを呼び出し
	onsenLogs, err := i.onsenLogService.GetOnsenLogsByUserID(ctx, userID)
	if err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return nil, err
	}

	// フォーマットに応じてエクスポート
	var data []byte
	switch strings.ToLower(format) {
	case "json":
		data, err = i.exportAsJSON(onsenLogs)
	case "csv":
		data, err = i.exportAsCSV(onsenLogs)
	default:
		err = fmt.Errorf("unsupported format: %s", format)
	}

	if err != nil {
		_ = i.outputPort.PresentError(ctx, err)
		return nil, err
	}

	// 出力ポートを呼び出し
	if err := i.outputPort.PresentExportedData(ctx, data, format); err != nil {
		return nil, err
	}

	return data, nil
}

// exportAsJSON はJSONフォーマットでエクスポートします
func (i *OnsenLogInteractor) exportAsJSON(onsenLogs []*entity.OnsenLog) ([]byte, error) {
	// 出力データを作成
	onsenLogOutputData := make([]port.OnsenLogOutputData, len(onsenLogs))
	for i, onsenLog := range onsenLogs {
		onsenLogOutputData[i] = port.OnsenLogOutputData{
			ID:         onsenLog.UUID,
			UserID:     onsenLog.UserID,
			Name:       onsenLog.Name,
			Location:   onsenLog.Location,
			SpringType: onsenLog.SpringType,
			Features:   onsenLog.Features,
			VisitDate:  onsenLog.VisitDate,
			Rating:     onsenLog.Rating,
			Comment:    onsenLog.Comment,
			CreatedAt:  onsenLog.CreatedAt,
			UpdatedAt:  onsenLog.UpdatedAt,
		}
	}

	// JSONにエンコード
	return json.MarshalIndent(onsenLogOutputData, "", "  ")
}

// exportAsCSV はCSVフォーマットでエクスポートします
func (i *OnsenLogInteractor) exportAsCSV(onsenLogs []*entity.OnsenLog) ([]byte, error) {
	// CSVデータを作成
	var sb strings.Builder
	writer := csv.NewWriter(&sb)

	// ヘッダーを書き込み
	header := []string{"ID", "温泉名", "所在地", "泉質", "特徴", "訪問日", "評価", "コメント", "作成日", "更新日"}
	if err := writer.Write(header); err != nil {
		return nil, err
	}

	// データを書き込み
	for _, onsenLog := range onsenLogs {
		// 特徴を文字列に変換
		features := make([]string, len(onsenLog.Features))
		for i, feature := range onsenLog.Features {
			features[i] = string(feature)
		}
		featuresStr := strings.Join(features, ", ")

		// 行データを作成
		row := []string{
			onsenLog.UUID,
			onsenLog.Name,
			onsenLog.Location,
			string(onsenLog.SpringType),
			featuresStr,
			onsenLog.VisitDate.Format("2006-01-02"),
			fmt.Sprintf("%d", onsenLog.Rating),
			onsenLog.Comment,
			onsenLog.CreatedAt.Format("2006-01-02 15:04:05"),
			onsenLog.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return []byte(sb.String()), nil
}
