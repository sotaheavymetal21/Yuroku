/**
 * バリデーションユーティリティ関数
 */

export interface ValidationResult {
  isValid: boolean;
  error?: string;
}

/**
 * メールアドレスのバリデーション
 */
export const validateEmail = (email: string): ValidationResult => {
  if (!email) {
    return { isValid: false, error: 'メールアドレスを入力してください' };
  }

  const emailRegex = /\S+@\S+\.\S+/;
  if (!emailRegex.test(email)) {
    return { isValid: false, error: '有効なメールアドレスを入力してください' };
  }

  return { isValid: true };
};

/**
 * パスワードのバリデーション
 */
export const validatePassword = (password: string, minLength: number = 8): ValidationResult => {
  if (!password) {
    return { isValid: false, error: 'パスワードを入力してください' };
  }

  if (password.length < minLength) {
    return { isValid: false, error: `パスワードは${minLength}文字以上で入力してください` };
  }

  return { isValid: true };
};

/**
 * 名前のバリデーション
 */
export const validateName = (name: string): ValidationResult => {
  if (!name) {
    return { isValid: false, error: '名前を入力してください' };
  }

  if (name.length < 2) {
    return { isValid: false, error: '名前は2文字以上で入力してください' };
  }

  return { isValid: true };
};

/**
 * 必須フィールドのバリデーション
 */
export const validateRequired = (value: string, fieldName: string): ValidationResult => {
  if (!value || value.trim() === '') {
    return { isValid: false, error: `${fieldName}を入力してください` };
  }

  return { isValid: true };
};

/**
 * パスワード確認のバリデーション
 */
export const validatePasswordConfirmation = (
  password: string,
  confirmPassword: string
): ValidationResult => {
  if (!confirmPassword) {
    return { isValid: false, error: '確認用パスワードを入力してください' };
  }

  if (password !== confirmPassword) {
    return { isValid: false, error: 'パスワードが一致しません' };
  }

  return { isValid: true };
};

/**
 * 複数のバリデーション結果を統合
 */
export const combineValidationResults = (
  ...results: ValidationResult[]
): ValidationResult => {
  for (const result of results) {
    if (!result.isValid) {
      return result;
    }
  }
  return { isValid: true };
};
