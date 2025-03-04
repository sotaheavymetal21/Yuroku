package controller

import (
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
	// コンテキストからユーザーIDを取得
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "認証されていません",
			},
		})
		return
	}

	var input struct {
		Name        string            `json:"name" binding:"required"`
		Location    string            `json:"location" binding:"required"`
		VisitDate   string            `json:"visit_date" binding:"required"`
		SpringType  entity.SpringType `json:"spring_type" binding:"required"`
		Temperature float64           `json:"temperature"`
		Rating      int               `json:"rating" binding:"required,min=1,max=5"`
		Notes       string            `json:"notes"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_INPUT",
				"message": "入力データが無効です: " + err.Error(),
			},
		})
		return
	}

	// 日付文字列をパース
	visitDate, err := time.Parse("2006-01-02", input.VisitDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_DATE",
				"message": "日付形式が無効です。YYYY-MM-DD形式で入力してください。",
			},
		})
		return
	}

	// ユースケースを呼び出し
	onsenLog, err := c.onsenLogUseCase.CreateOnsenLog(
		ctx.Request.Context(),
		userID.(string),
		input.Name,
		input.Location,
		visitDate,
		input.SpringType,
		input.Temperature,
		input.Rating,
		input.Notes,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "CREATE_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"id":          onsenLog.ID,
			"name":        onsenLog.Name,
			"location":    onsenLog.Location,
			"visit_date":  onsenLog.VisitDate.Format("2006-01-02"),
			"spring_type": onsenLog.SpringType,
			"temperature": onsenLog.Temperature,
			"rating":      onsenLog.Rating,
			"notes":       onsenLog.Notes,
			"created_at":  onsenLog.CreatedAt,
		},
		"message": "温泉メモの作成に成功しました",
	})
}

// GetOnsenLog は指定されたIDの温泉メモを取得します
func (c *OnsenLogController) GetOnsenLog(ctx *gin.Context) {
	// コンテキストからユーザーIDを取得
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "認証されていません",
			},
		})
		return
	}

	// パスパラメータからIDを取得
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "MISSING_ID",
				"message": "温泉メモIDが指定されていません",
			},
		})
		return
	}

	// ユースケースを呼び出し
	onsenLog, err := c.onsenLogUseCase.GetOnsenLog(ctx.Request.Context(), id, userID.(string))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": err.Error(),
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":          onsenLog.ID,
			"name":        onsenLog.Name,
			"location":    onsenLog.Location,
			"visit_date":  onsenLog.VisitDate.Format("2006-01-02"),
			"spring_type": onsenLog.SpringType,
			"temperature": onsenLog.Temperature,
			"rating":      onsenLog.Rating,
			"notes":       onsenLog.Notes,
			"created_at":  onsenLog.CreatedAt,
			"updated_at":  onsenLog.UpdatedAt,
		},
		"message": "温泉メモの取得に成功しました",
	})
}

// GetOnsenLogs はユーザーの全温泉メモをページネーションで取得します
func (c *OnsenLogController) GetOnsenLogs(ctx *gin.Context) {
	// コンテキストからユーザーIDを取得
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "認証されていません",
			},
		})
		return
	}

	// クエリパラメータからページとリミットを取得
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	// ユースケースを呼び出し
	onsenLogs, totalCount, err := c.onsenLogUseCase.GetOnsenLogs(ctx.Request.Context(), userID.(string), page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "FETCH_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	// レスポンスデータを作成
	var logsData []gin.H
	for _, log := range onsenLogs {
		logsData = append(logsData, gin.H{
			"id":          log.ID,
			"name":        log.Name,
			"location":    log.Location,
			"visit_date":  log.VisitDate.Format("2006-01-02"),
			"spring_type": log.SpringType,
			"temperature": log.Temperature,
			"rating":      log.Rating,
			"notes":       log.Notes,
			"created_at":  log.CreatedAt,
			"updated_at":  log.UpdatedAt,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"logs":        logsData,
			"total_count": totalCount,
			"page":        page,
			"limit":       limit,
			"total_pages": (totalCount + limit - 1) / limit,
		},
		"message": "温泉メモの取得に成功しました",
	})
}

// GetFilteredOnsenLogs はフィルター条件に基づいて温泉メモを取得します
func (c *OnsenLogController) GetFilteredOnsenLogs(ctx *gin.Context) {
	// コンテキストからユーザーIDを取得
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "認証されていません",
			},
		})
		return
	}

	// クエリパラメータからフィルター条件を取得
	springType := entity.SpringType(ctx.Query("spring_type"))
	location := ctx.Query("location")

	minRating, err := strconv.Atoi(ctx.DefaultQuery("min_rating", "0"))
	if err != nil || minRating < 0 || minRating > 5 {
		minRating = 0
	}

	var startDate, endDate *time.Time
	if startDateStr := ctx.Query("start_date"); startDateStr != "" {
		date, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			startDate = &date
		}
	}

	if endDateStr := ctx.Query("end_date"); endDateStr != "" {
		date, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			endDate = &date
		}
	}

	// ページネーションパラメータ
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	// ユースケースを呼び出し
	onsenLogs, totalCount, err := c.onsenLogUseCase.GetFilteredOnsenLogs(
		ctx.Request.Context(),
		userID.(string),
		springType,
		location,
		minRating,
		startDate,
		endDate,
		page,
		limit,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "FETCH_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	// レスポンスデータを作成
	var logsData []gin.H
	for _, log := range onsenLogs {
		logsData = append(logsData, gin.H{
			"id":          log.ID,
			"name":        log.Name,
			"location":    log.Location,
			"visit_date":  log.VisitDate.Format("2006-01-02"),
			"spring_type": log.SpringType,
			"temperature": log.Temperature,
			"rating":      log.Rating,
			"notes":       log.Notes,
			"created_at":  log.CreatedAt,
			"updated_at":  log.UpdatedAt,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"logs":        logsData,
			"total_count": totalCount,
			"page":        page,
			"limit":       limit,
			"total_pages": (totalCount + limit - 1) / limit,
		},
		"message": "温泉メモの取得に成功しました",
	})
}

// UpdateOnsenLog は温泉メモを更新します
func (c *OnsenLogController) UpdateOnsenLog(ctx *gin.Context) {
	// コンテキストからユーザーIDを取得
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "認証されていません",
			},
		})
		return
	}

	// パスパラメータからIDを取得
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "MISSING_ID",
				"message": "温泉メモIDが指定されていません",
			},
		})
		return
	}

	var input struct {
		Name        string            `json:"name" binding:"required"`
		Location    string            `json:"location" binding:"required"`
		VisitDate   string            `json:"visit_date" binding:"required"`
		SpringType  entity.SpringType `json:"spring_type" binding:"required"`
		Temperature float64           `json:"temperature"`
		Rating      int               `json:"rating" binding:"required,min=1,max=5"`
		Notes       string            `json:"notes"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_INPUT",
				"message": "入力データが無効です: " + err.Error(),
			},
		})
		return
	}

	// 日付文字列をパース
	visitDate, err := time.Parse("2006-01-02", input.VisitDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_DATE",
				"message": "日付形式が無効です。YYYY-MM-DD形式で入力してください。",
			},
		})
		return
	}

	// ユースケースを呼び出し
	onsenLog, err := c.onsenLogUseCase.UpdateOnsenLog(
		ctx.Request.Context(),
		id,
		userID.(string),
		input.Name,
		input.Location,
		visitDate,
		input.SpringType,
		input.Temperature,
		input.Rating,
		input.Notes,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "UPDATE_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":          onsenLog.ID,
			"name":        onsenLog.Name,
			"location":    onsenLog.Location,
			"visit_date":  onsenLog.VisitDate.Format("2006-01-02"),
			"spring_type": onsenLog.SpringType,
			"temperature": onsenLog.Temperature,
			"rating":      onsenLog.Rating,
			"notes":       onsenLog.Notes,
			"updated_at":  onsenLog.UpdatedAt,
		},
		"message": "温泉メモの更新に成功しました",
	})
}

// DeleteOnsenLog は温泉メモを削除します
func (c *OnsenLogController) DeleteOnsenLog(ctx *gin.Context) {
	// コンテキストからユーザーIDを取得
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "認証されていません",
			},
		})
		return
	}

	// パスパラメータからIDを取得
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "MISSING_ID",
				"message": "温泉メモIDが指定されていません",
			},
		})
		return
	}

	// ユースケースを呼び出し
	err := c.onsenLogUseCase.DeleteOnsenLog(ctx.Request.Context(), id, userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "DELETE_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "温泉メモの削除に成功しました",
	})
}

// ExportOnsenLogs は温泉メモをエクスポートします
func (c *OnsenLogController) ExportOnsenLogs(ctx *gin.Context) {
	// コンテキストからユーザーIDを取得
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "認証されていません",
			},
		})
		return
	}

	// クエリパラメータからフォーマットを取得
	format := ctx.DefaultQuery("format", "json")
	if format != "json" && format != "csv" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_FORMAT",
				"message": "サポートされていないフォーマットです。json または csv を指定してください。",
			},
		})
		return
	}

	// ユースケースを呼び出し
	data, contentType, filename, err := c.onsenLogUseCase.ExportOnsenLogs(ctx.Request.Context(), userID.(string), format)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "EXPORT_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	// ファイルをダウンロードとして送信
	ctx.Header("Content-Disposition", "attachment; filename="+filename)
	ctx.Data(http.StatusOK, contentType, data)
}
