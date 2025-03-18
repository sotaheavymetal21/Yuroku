package common

import (
	"regexp"
	"strings"
)

// ValidationError はバリデーションエラーを表します
type ValidationError struct {
	Field   string
	Message string
}

// Error はエラーメッセージを返します
func (e ValidationError) Error() string {
	return e.Message
}

// ValidateEmail はメールアドレスの形式を検証します
func ValidateEmail(email string) *ValidationError {
	if email == "" {
		return &ValidationError{
			Field:   "email",
			Message: "メールアドレスは必須です",
		}
	}

	// 簡易的なメールアドレス形式をチェック
	pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	if !matched {
		return &ValidationError{
			Field:   "email",
			Message: "メールアドレスの形式が無効です",
		}
	}

	return nil
}

// ValidatePassword はパスワードの強度を検証します
func ValidatePassword(password string) *ValidationError {
	if password == "" {
		return &ValidationError{
			Field:   "password",
			Message: "パスワードは必須です",
		}
	}

	if len(password) < 8 {
		return &ValidationError{
			Field:   "password",
			Message: "パスワードは8文字以上である必要があります",
		}
	}

	// 英字と数字が含まれているか
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)

	if !hasLetter || !hasDigit {
		return &ValidationError{
			Field:   "password",
			Message: "パスワードには英字と数字の両方を含める必要があります",
		}
	}

	return nil
}

// ValidateRequired は必須フィールドを検証します
func ValidateRequired(field, value, fieldName string) *ValidationError {
	if strings.TrimSpace(value) == "" {
		return &ValidationError{
			Field:   field,
			Message: fieldName + "は必須です",
		}
	}
	return nil
}

// ValidateRange は値が指定された範囲内かどうかを検証します
func ValidateRange(field string, value, min, max int, fieldName string) *ValidationError {
	if value < min || value > max {
		return &ValidationError{
			Field:   field,
			Message: fieldName + "は" + string(min) + "から" + string(max) + "の間でなければなりません",
		}
	}
	return nil
}
