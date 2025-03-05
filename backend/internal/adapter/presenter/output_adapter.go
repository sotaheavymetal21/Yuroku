package presenter

import (
	"context"

	"github.com/yourusername/yuroku/internal/usecase/port"
)

// AuthOutputAdapter はAuthPresenterをAuthOutputPortに適応させるアダプターです
type AuthOutputAdapter struct {
	Presenter port.AuthPresenterPort
}

// NewAuthOutputAdapter は新しいAuthOutputAdapterインスタンスを作成します
func NewAuthOutputAdapter(presenter port.AuthPresenterPort) port.AuthOutputPort {
	return &AuthOutputAdapter{
		Presenter: presenter,
	}
}

// PresentRegister はユーザー登録結果を表示します
func (a *AuthOutputAdapter) PresentRegister(ctx context.Context, data port.AuthOutputData) error {
	a.Presenter.PresentToken(data.AccessToken, data.RefreshToken)
	return nil
}

// PresentLogin はログイン結果を表示します
func (a *AuthOutputAdapter) PresentLogin(ctx context.Context, data port.AuthOutputData) error {
	a.Presenter.PresentToken(data.AccessToken, data.RefreshToken)
	return nil
}

// PresentRefreshToken はトークン更新結果を表示します
func (a *AuthOutputAdapter) PresentRefreshToken(ctx context.Context, data port.AuthOutputData) error {
	a.Presenter.PresentToken(data.AccessToken, data.RefreshToken)
	return nil
}

// PresentError はエラーを表示します
func (a *AuthOutputAdapter) PresentError(ctx context.Context, err error) error {
	a.Presenter.PresentError(err)
	return nil
}

// OnsenLogOutputAdapter はOnsenLogPresenterをOnsenLogOutputPortに適応させるアダプターです
type OnsenLogOutputAdapter struct {
	Presenter port.OnsenLogPresenterPort
}

// NewOnsenLogOutputAdapter は新しいOnsenLogOutputAdapterインスタンスを作成します
func NewOnsenLogOutputAdapter(presenter port.OnsenLogPresenterPort) port.OnsenLogOutputPort {
	return &OnsenLogOutputAdapter{
		Presenter: presenter,
	}
}

// PresentOnsenLog は温泉メモを表示します
func (a *OnsenLogOutputAdapter) PresentOnsenLog(ctx context.Context, data port.OnsenLogOutputData) error {
	return nil
}

// PresentOnsenLogs は温泉メモのリストを表示します
func (a *OnsenLogOutputAdapter) PresentOnsenLogs(ctx context.Context, data port.OnsenLogsOutputData) error {
	return nil
}

// PresentExportedData はエクスポートされたデータを表示します
func (a *OnsenLogOutputAdapter) PresentExportedData(ctx context.Context, data []byte, format string) error {
	return nil
}

// PresentError はエラーを表示します
func (a *OnsenLogOutputAdapter) PresentError(ctx context.Context, err error) error {
	return nil
}

// OnsenImageOutputAdapter はOnsenImagePresenterをOnsenImageOutputPortに適応させるアダプターです
type OnsenImageOutputAdapter struct {
	Presenter port.OnsenImagePresenterPort
}

// NewOnsenImageOutputAdapter は新しいOnsenImageOutputAdapterインスタンスを作成します
func NewOnsenImageOutputAdapter(presenter port.OnsenImagePresenterPort) port.OnsenImageOutputPort {
	return &OnsenImageOutputAdapter{
		Presenter: presenter,
	}
}

// PresentImage は温泉画像を表示します
func (a *OnsenImageOutputAdapter) PresentImage(ctx context.Context, data port.ImageOutputData) error {
	return nil
}

// PresentImages は温泉画像のリストを表示します
func (a *OnsenImageOutputAdapter) PresentImages(ctx context.Context, data []port.ImageOutputData) error {
	return nil
}

// PresentError はエラーを表示します
func (a *OnsenImageOutputAdapter) PresentError(ctx context.Context, err error) error {
	return nil
}
