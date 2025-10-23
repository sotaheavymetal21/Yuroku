package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/yuroku/internal/common"
)

// ErrorResponse はエラーレスポンスの構造を表します
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail はエラーの詳細を表します
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// SuccessResponse は成功レスポンスの構造を表します
type SuccessResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
}

// RespondWithError は統一されたエラーレスポンスを返します
func RespondWithError(ctx *gin.Context, statusCode int, code, message string) {
	ctx.JSON(statusCode, ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
}

// RespondWithSuccess は統一された成功レスポンスを返します
func RespondWithSuccess(ctx *gin.Context, statusCode int, data interface{}, message string) {
	response := SuccessResponse{
		Message: message,
	}
	if data != nil {
		response.Data = data
	}
	ctx.JSON(statusCode, response)
}

// RespondWithAppError はAppErrorから適切なHTTPステータスコードとレスポンスを返します
func RespondWithAppError(ctx *gin.Context, err error) {
	if appErr := common.GetAppError(err); appErr != nil {
		statusCode := getHTTPStatusFromAppError(appErr)
		RespondWithError(ctx, statusCode, appErr.Code, appErr.Message)
		return
	}

	// AppErrorでない場合は内部エラーとして扱う
	RespondWithError(ctx, http.StatusInternalServerError, common.ErrInternal, err.Error())
}

// getHTTPStatusFromAppError はAppErrorのコードから適切なHTTPステータスコードを返します
func getHTTPStatusFromAppError(appErr *common.AppError) int {
	switch appErr.Code {
	case common.ErrNotFound:
		return http.StatusNotFound
	case common.ErrInvalidInput, common.ErrValidation:
		return http.StatusBadRequest
	case common.ErrUnauthorized, common.ErrAuthentication, common.ErrTokenExpired:
		return http.StatusUnauthorized
	case common.ErrForbidden:
		return http.StatusForbidden
	case common.ErrDuplicate:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

// GetUserID はコンテキストからユーザーIDを取得します
// ユーザーIDが存在しない場合はエラーレスポンスを返してfalseを返します
func GetUserID(ctx *gin.Context) (string, bool) {
	userID, exists := ctx.Get("userID")
	if !exists {
		RespondWithError(ctx, http.StatusUnauthorized, common.ErrUnauthorized, "認証が必要です")
		return "", false
	}
	return userID.(string), true
}

// ValidateBindJSON はJSONのバインドを行い、エラーがあればエラーレスポンスを返します
func ValidateBindJSON(ctx *gin.Context, input interface{}) bool {
	if err := ctx.ShouldBindJSON(input); err != nil {
		RespondWithError(ctx, http.StatusBadRequest, common.ErrInvalidInput, "入力データが無効です: "+err.Error())
		return false
	}
	return true
}

// ValidatePathParam はパスパラメータの存在を検証します
func ValidatePathParam(ctx *gin.Context, paramName, errorMessage string) (string, bool) {
	value := ctx.Param(paramName)
	if value == "" {
		RespondWithError(ctx, http.StatusBadRequest, common.ErrInvalidInput, errorMessage)
		return "", false
	}
	return value, true
}
