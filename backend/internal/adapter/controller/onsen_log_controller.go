package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/yuroku/internal/domain/entity"
	"github.com/yourusername/yuroku/internal/usecase/port"
)

// OnsenLogController は温泉メモ関連のコントローラーです
type OnsenLogController struct {
	onsenLogUseCase port.OnsenLogInputPort
}

// NewOnsenLogController は新しい温泉メモコントローラーを作成します
func NewOnsenLogController(onsenLogUseCase port.OnsenLogInputPort) *OnsenLogController {
	return &OnsenLogController{
		onsenLogUseCase: onsenLogUseCase,
	}
}

// CreateOnsenLog は新しい温泉メモを作成します
func (c *OnsenLogController) CreateOnsenLog(ctx *gin.Context) {
	// ユーザーIDを取得
	userID, ok := GetUserID(ctx)
	if !ok {
		return
	}

	// リクエストボディをバインド
	var input struct {
		Name       string            `json:"name" binding:"required"`
		Location   string            `json:"location" binding:"required"`
		SpringType entity.SpringType `json:"spring_type" binding:"required"`
		Features   []entity.Feature  `json:"features"`
		VisitDate  string            `json:"visit_date" binding:"required"`
		Rating     int               `json:"rating" binding:"required,min=1,max=5"`
		Comment    string            `json:"comment"`
	}

	if !ValidateBindJSON(ctx, &input) {
		return
	}

	// 日付をパース
	visitDate, err := time.Parse("2006-01-02", input.VisitDate)
	if err != nil {
		RespondWithError(ctx, http.StatusBadRequest, "INVALID_DATE", "日付の形式が無効です（YYYY-MM-DD）")
		return
	}

	// 入力データを作成
	createInput := port.CreateOnsenLogInput{
		UserID:     userID,
		Name:       input.Name,
		Location:   input.Location,
		SpringType: input.SpringType,
		Features:   input.Features,
		VisitDate:  visitDate,
		Rating:     input.Rating,
		Comment:    input.Comment,
	}

	// ユースケースを呼び出し
	onsenLog, err := c.onsenLogUseCase.CreateOnsenLog(
		ctx.Request.Context(),
		createInput,
	)

	if err != nil {
		RespondWithAppError(ctx, err)
		return
	}

	// レスポンスを返す
	RespondWithSuccess(ctx, http.StatusCreated, gin.H{
		"id":          onsenLog.ID,
		"name":        onsenLog.Name,
		"location":    onsenLog.Location,
		"spring_type": onsenLog.SpringType,
		"features":    onsenLog.Features,
		"visit_date":  onsenLog.VisitDate.Format("2006-01-02"),
		"rating":      onsenLog.Rating,
		"comment":     onsenLog.Comment,
		"created_at":  onsenLog.CreatedAt,
		"updated_at":  onsenLog.UpdatedAt,
		"images":      onsenLog.Images,
	}, "温泉メモを作成しました")
}

// GetOnsenLog は特定の温泉メモを取得します
func (c *OnsenLogController) GetOnsenLog(ctx *gin.Context) {
	// ユーザーIDを取得
	userID, ok := GetUserID(ctx)
	if !ok {
		return
	}

	// パスパラメータからIDを取得
	id, ok := ValidatePathParam(ctx, "id", "温泉メモIDが必要です")
	if !ok {
		return
	}

	// ユースケースを呼び出し
	onsenLog, err := c.onsenLogUseCase.GetOnsenLog(
		ctx.Request.Context(),
		id,
		userID,
	)

	if err != nil {
		RespondWithAppError(ctx, err)
		return
	}

	// レスポンスを返す
	RespondWithSuccess(ctx, http.StatusOK, gin.H{
		"id":          onsenLog.ID,
		"name":        onsenLog.Name,
		"location":    onsenLog.Location,
		"spring_type": onsenLog.SpringType,
		"features":    onsenLog.Features,
		"visit_date":  onsenLog.VisitDate.Format("2006-01-02"),
		"rating":      onsenLog.Rating,
		"comment":     onsenLog.Comment,
		"created_at":  onsenLog.CreatedAt,
		"updated_at":  onsenLog.UpdatedAt,
		"images":      onsenLog.Images,
	}, "温泉メモを取得しました")
}

// GetOnsenLogs はユーザーの温泉メモリストを取得します
func (c *OnsenLogController) GetOnsenLogs(ctx *gin.Context) {
	// ユーザーIDを取得
	userID, ok := GetUserID(ctx)
	if !ok {
		return
	}

	// クエリパラメータを取得
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")

	// 数値に変換
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	// ユースケースを呼び出し
	result, err := c.onsenLogUseCase.GetOnsenLogs(
		ctx.Request.Context(),
		userID,
		page,
		limit,
	)

	if err != nil {
		RespondWithAppError(ctx, err)
		return
	}

	// レスポンスを返す
	RespondWithSuccess(ctx, http.StatusOK, gin.H{
		"onsen_logs":  result.OnsenLogs,
		"total_count": result.TotalCount,
		"page":        result.Page,
		"limit":       result.Limit,
	}, "温泉メモリストを取得しました")
}

// GetFilteredOnsenLogs はフィルタリングされた温泉メモリストを取得します
func (c *OnsenLogController) GetFilteredOnsenLogs(ctx *gin.Context) {
	// ユーザーIDを取得
	userID, ok := GetUserID(ctx)
	if !ok {
		return
	}

	// クエリパラメータを取得
	springTypeStr := ctx.Query("spring_type")
	location := ctx.Query("location")
	minRatingStr := ctx.DefaultQuery("min_rating", "0")
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")

	// 数値に変換
	minRating, err := strconv.Atoi(minRatingStr)
	if err != nil || minRating < 0 || minRating > 5 {
		minRating = 0
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	// 日付をパース
	var startDate, endDate *time.Time
	if startDateStr != "" {
		parsedStartDate, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			startDate = &parsedStartDate
		}
	}

	if endDateStr != "" {
		parsedEndDate, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			endDate = &parsedEndDate
		}
	}

	// SpringTypeを変換
	var springType entity.SpringType
	if springTypeStr != "" {
		springType = entity.SpringType(springTypeStr)
	}

	// 入力データを作成
	filterInput := port.FilterOnsenLogsInput{
		UserID:     userID,
		SpringType: springType,
		Location:   location,
		MinRating:  minRating,
		StartDate:  startDate,
		EndDate:    endDate,
		Page:       page,
		Limit:      limit,
	}

	// ユースケースを呼び出し
	result, err := c.onsenLogUseCase.GetFilteredOnsenLogs(
		ctx.Request.Context(),
		filterInput,
	)

	if err != nil {
		RespondWithAppError(ctx, err)
		return
	}

	// レスポンスを返す
	RespondWithSuccess(ctx, http.StatusOK, gin.H{
		"onsen_logs":  result.OnsenLogs,
		"total_count": result.TotalCount,
		"page":        result.Page,
		"limit":       result.Limit,
	}, "フィルタリングされた温泉メモリストを取得しました")
}

// UpdateOnsenLog は温泉メモを更新します
func (c *OnsenLogController) UpdateOnsenLog(ctx *gin.Context) {
	// ユーザーIDを取得
	userID, ok := GetUserID(ctx)
	if !ok {
		return
	}

	// パスパラメータからIDを取得
	id, ok := ValidatePathParam(ctx, "id", "温泉メモIDが必要です")
	if !ok {
		return
	}

	// リクエストボディをバインド
	var input struct {
		Name       string            `json:"name" binding:"required"`
		Location   string            `json:"location" binding:"required"`
		SpringType entity.SpringType `json:"spring_type" binding:"required"`
		Features   []entity.Feature  `json:"features"`
		VisitDate  string            `json:"visit_date" binding:"required"`
		Rating     int               `json:"rating" binding:"required,min=1,max=5"`
		Comment    string            `json:"comment"`
	}

	if !ValidateBindJSON(ctx, &input) {
		return
	}

	// 日付をパース
	visitDate, err := time.Parse("2006-01-02", input.VisitDate)
	if err != nil {
		RespondWithError(ctx, http.StatusBadRequest, "INVALID_DATE", "日付の形式が無効です（YYYY-MM-DD）")
		return
	}

	// 入力データを作成
	updateInput := port.UpdateOnsenLogInput{
		ID:         id,
		UserID:     userID,
		Name:       input.Name,
		Location:   input.Location,
		SpringType: input.SpringType,
		Features:   input.Features,
		VisitDate:  visitDate,
		Rating:     input.Rating,
		Comment:    input.Comment,
	}

	// ユースケースを呼び出し
	onsenLog, err := c.onsenLogUseCase.UpdateOnsenLog(
		ctx.Request.Context(),
		updateInput,
	)

	if err != nil {
		RespondWithAppError(ctx, err)
		return
	}

	// レスポンスを返す
	RespondWithSuccess(ctx, http.StatusOK, gin.H{
		"id":          onsenLog.ID,
		"name":        onsenLog.Name,
		"location":    onsenLog.Location,
		"spring_type": onsenLog.SpringType,
		"features":    onsenLog.Features,
		"visit_date":  onsenLog.VisitDate.Format("2006-01-02"),
		"rating":      onsenLog.Rating,
		"comment":     onsenLog.Comment,
		"created_at":  onsenLog.CreatedAt,
		"updated_at":  onsenLog.UpdatedAt,
		"images":      onsenLog.Images,
	}, "温泉メモを更新しました")
}

// DeleteOnsenLog は温泉メモを削除します
func (c *OnsenLogController) DeleteOnsenLog(ctx *gin.Context) {
	// コンテキストからユーザーIDを取得
	userID, ok := GetUserID(ctx)
	if !ok {
		return
	}

	// パスパラメータからIDを取得
	id, ok := ValidatePathParam(ctx, "id", "温泉メモIDが指定されていません")
	if !ok {
		return
	}

	// ユースケースを呼び出し
	err := c.onsenLogUseCase.DeleteOnsenLog(ctx.Request.Context(), id, userID)
	if err != nil {
		RespondWithAppError(ctx, err)
		return
	}

	RespondWithSuccess(ctx, http.StatusOK, nil, "温泉メモの削除に成功しました")
}

// ExportOnsenLogs は温泉メモをエクスポートします
func (c *OnsenLogController) ExportOnsenLogs(ctx *gin.Context) {
	// ユーザーIDを取得
	userID, ok := GetUserID(ctx)
	if !ok {
		return
	}

	// クエリパラメータからフォーマットを取得
	format := ctx.DefaultQuery("format", "json")
	if format != "json" && format != "csv" {
		RespondWithError(ctx, http.StatusBadRequest, "INVALID_FORMAT", "サポートされているフォーマットは json または csv です")
		return
	}

	// ユースケースを呼び出し
	data, err := c.onsenLogUseCase.ExportOnsenLogs(
		ctx.Request.Context(),
		userID,
		format,
	)

	if err != nil {
		RespondWithAppError(ctx, err)
		return
	}

	// Content-Typeを設定
	var contentType string
	var filename string
	switch format {
	case "json":
		contentType = "application/json"
		filename = "onsen_logs.json"
	case "csv":
		contentType = "text/csv"
		filename = "onsen_logs.csv"
	}

	// レスポンスヘッダーを設定
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Data(http.StatusOK, contentType, data)
}
