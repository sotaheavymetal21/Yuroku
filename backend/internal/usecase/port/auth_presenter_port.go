package port

import (
	"github.com/yourusername/yuroku/internal/domain/entity"
)

// AuthPresenterPort は認証関連のレスポンスをフォーマットするインターフェースです
type AuthPresenterPort interface {
	// PresentUser はユーザー情報のレスポンスをフォーマットします
	PresentUser(user *entity.User) map[string]interface{}

	// PresentToken は認証トークンのレスポンスをフォーマットします
	PresentToken(token, refreshToken string) map[string]interface{}

	// PresentError はエラーレスポンスをフォーマットします
	PresentError(err error) map[string]interface{}
}
