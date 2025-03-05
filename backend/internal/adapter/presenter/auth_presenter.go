package presenter

import (
	"github.com/yourusername/yuroku/internal/domain/entity"
	"github.com/yourusername/yuroku/internal/usecase/port"
)

// AuthPresenter は認証関連のレスポンスをフォーマットするプレゼンターです
type AuthPresenter struct{}

// NewAuthPresenter は新しいAuthPresenterインスタンスを作成します
func NewAuthPresenter() port.AuthPresenterPort {
	return &AuthPresenter{}
}

// PresentUser はユーザー情報のレスポンスをフォーマットします
func (p *AuthPresenter) PresentUser(user *entity.User) map[string]interface{} {
	return map[string]interface{}{
		"id":        user.ID.Hex(),
		"username":  user.Username,
		"email":     user.Email,
		"createdAt": user.CreatedAt,
		"updatedAt": user.UpdatedAt,
	}
}

// PresentToken は認証トークンのレスポンスをフォーマットします
func (p *AuthPresenter) PresentToken(token, refreshToken string) map[string]interface{} {
	return map[string]interface{}{
		"token":        token,
		"refreshToken": refreshToken,
	}
}

// PresentError はエラーレスポンスをフォーマットします
func (p *AuthPresenter) PresentError(err error) map[string]interface{} {
	return map[string]interface{}{
		"error": map[string]interface{}{
			"message": err.Error(),
		},
	}
}
