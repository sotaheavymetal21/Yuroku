package port

import "context"

// AuthInputPort は認証ユースケースの入力ポートです
type AuthInputPort interface {
	// Register はユーザー登録を行います
	Register(ctx context.Context, name, email, password string) error

	// Login はユーザーログインを行います
	Login(ctx context.Context, email, password string) (string, string, error)

	// VerifyToken はトークンを検証します
	VerifyToken(ctx context.Context, token string) (string, error)

	// RefreshToken はトークンを更新します
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)

	// Logout はユーザーログアウトを行います
	Logout(ctx context.Context, token string) error

	// GetUserByID はIDでユーザーを取得します
	GetUserByID(ctx context.Context, id string) (UserOutputData, error)

	// UpdateProfile はユーザープロフィールを更新します
	UpdateProfile(ctx context.Context, userID, name, email string) error

	// ChangePassword はパスワードを変更します
	ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error

	// DeleteAccount はアカウントを削除します
	DeleteAccount(ctx context.Context, userID, password string) error
}

// AuthOutputPort は認証ユースケースの出力ポートです
type AuthOutputPort interface {
	// PresentRegister はユーザー登録結果を表示します
	PresentRegister(ctx context.Context, data AuthOutputData) error

	// PresentLogin はユーザーログイン結果を表示します
	PresentLogin(ctx context.Context, data AuthOutputData) error

	// PresentRefreshToken はトークン更新結果を表示します
	PresentRefreshToken(ctx context.Context, data AuthOutputData) error

	// PresentError はエラーを表示します
	PresentError(ctx context.Context, err error) error
}

// AuthOutputData は認証ユースケースの出力データです
type AuthOutputData struct {
	UserID       string `json:"user_id"`
	Email        string `json:"email"`
	Name         string `json:"name,omitempty"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// UserOutputData はユーザー情報の出力データです
type UserOutputData struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
