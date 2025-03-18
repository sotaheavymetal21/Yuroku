package common

import (
	"errors"
	"fmt"
)

// AppError はアプリケーションのエラーを表します
type AppError struct {
	Code    string
	Message string
	Err     error
}

// AppErrorCode はアプリケーションエラーコードを定義します
const (
	ErrNotFound       = "NOT_FOUND"
	ErrInvalidInput   = "INVALID_INPUT"
	ErrUnauthorized   = "UNAUTHORIZED"
	ErrForbidden      = "FORBIDDEN"
	ErrInternal       = "INTERNAL_ERROR"
	ErrDuplicate      = "DUPLICATE_ENTITY"
	ErrAuthentication = "AUTHENTICATION_ERROR"
	ErrTokenExpired   = "TOKEN_EXPIRED"
	ErrDatabaseError  = "DATABASE_ERROR"
	ErrValidation     = "VALIDATION_ERROR"
)

// Error はエラーメッセージを返します
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap はラップされたエラーを返します
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError は新しいAppErrorを作成します
func NewAppError(code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// IsAppError はエラーがAppErrorかどうかを判定します
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// GetAppError はエラーからAppErrorを取得します
func GetAppError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return nil
}

// GetErrorCode はエラーからコードを取得します
func GetErrorCode(err error) string {
	if appErr := GetAppError(err); appErr != nil {
		return appErr.Code
	}
	return ErrInternal
}

// NewNotFoundError は「リソースが見つからない」エラーを作成します
func NewNotFoundError(message string, err error) *AppError {
	return NewAppError(ErrNotFound, message, err)
}

// NewInvalidInputError は「入力が無効」エラーを作成します
func NewInvalidInputError(message string, err error) *AppError {
	return NewAppError(ErrInvalidInput, message, err)
}

// NewUnauthorizedError は「認証されていない」エラーを作成します
func NewUnauthorizedError(message string, err error) *AppError {
	return NewAppError(ErrUnauthorized, message, err)
}

// NewForbiddenError は「権限がない」エラーを作成します
func NewForbiddenError(message string, err error) *AppError {
	return NewAppError(ErrForbidden, message, err)
}

// NewInternalError は「内部エラー」を作成します
func NewInternalError(message string, err error) *AppError {
	return NewAppError(ErrInternal, message, err)
}

// NewDuplicateError は「重複エラー」を作成します
func NewDuplicateError(message string, err error) *AppError {
	return NewAppError(ErrDuplicate, message, err)
}

// NewAuthenticationError は「認証エラー」を作成します
func NewAuthenticationError(message string, err error) *AppError {
	return NewAppError(ErrAuthentication, message, err)
}

// NewTokenExpiredError は「トークン期限切れエラー」を作成します
func NewTokenExpiredError(message string, err error) *AppError {
	return NewAppError(ErrTokenExpired, message, err)
}

// NewDatabaseError は「データベースエラー」を作成します
func NewDatabaseError(message string, err error) *AppError {
	return NewAppError(ErrDatabaseError, message, err)
}

// NewValidationError は「バリデーションエラー」を作成します
func NewValidationError(message string, err error) *AppError {
	return NewAppError(ErrValidation, message, err)
}
