package presenter

import (
	"github.com/yourusername/yuroku/internal/domain/entity"
	"github.com/yourusername/yuroku/internal/usecase/port"
)

// 共通レスポンス構造体
type StandardResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   *ErrorData  `json:"error,omitempty"`
}

// エラーデータ構造体
type ErrorData struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// AuthPresenter は認証関連のレスポンスをフォーマットするプレゼンターです
type AuthPresenter struct{}

// NewAuthPresenter は新しいAuthPresenterインスタンスを作成します
func NewAuthPresenter() port.AuthPresenterPort {
	return &AuthPresenter{}
}

// PresentUser はユーザー情報のレスポンスをフォーマットします
func (p *AuthPresenter) PresentUser(user *entity.User) map[string]interface{} {
	userData := map[string]interface{}{
		"id":         user.ID.Hex(),
		"uuid":       user.UUID,
		"name":       user.Name,
		"email":      user.Email,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}

	return map[string]interface{}{
		"data":    userData,
		"message": "ユーザー情報を取得しました",
	}
}

// PresentToken は認証トークンのレスポンスをフォーマットします
func (p *AuthPresenter) PresentToken(token, refreshToken string) map[string]interface{} {
	tokenData := map[string]interface{}{
		"access_token":  token,
		"refresh_token": refreshToken,
	}

	return map[string]interface{}{
		"data":    tokenData,
		"message": "認証に成功しました",
	}
}

// PresentError はエラーレスポンスをフォーマットします
func (p *AuthPresenter) PresentError(err error) map[string]interface{} {
	// エラーコードを決定（デフォルトはGENERAL_ERROR）
	errorCode := "GENERAL_ERROR"

	// エラーメッセージによってコードを判断
	errMsg := err.Error()
	switch {
	case contains(errMsg, "パスワード"):
		errorCode = "INVALID_PASSWORD"
	case contains(errMsg, "メールアドレス"):
		errorCode = "INVALID_EMAIL"
	case contains(errMsg, "ユーザー"):
		errorCode = "USER_ERROR"
	case contains(errMsg, "トークン"):
		errorCode = "TOKEN_ERROR"
	case contains(errMsg, "認証"):
		errorCode = "AUTHENTICATION_ERROR"
	}

	return map[string]interface{}{
		"error": map[string]interface{}{
			"code":    errorCode,
			"message": errMsg,
		},
	}
}

// contains は文字列に特定の部分文字列が含まれるかをチェックします
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[0:len(substr)] == substr
}
