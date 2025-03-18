import { ApiError, ServerError } from '@/types';

/**
 * APIエラーを処理し、ユーザーフレンドリーなメッセージを返す
 */
export const handleApiError = (error: unknown): string => {
  // APIエラーの場合
  if (isApiError(error)) {
    return error.error.message;
  } 
  
  // サーバーエラーの場合
  if (isServerError(error)) {
    return error.message;
  }
  
  // その他のエラー
  if (error instanceof Error) {
    return error.message;
  }
  
  // 型が不明なエラー
  return '予期しないエラーが発生しました。もう一度お試しください。';
};

/**
 * エラーがApiErrorかどうかを判定
 */
export const isApiError = (error: unknown): error is ApiError => {
  return (
    typeof error === 'object' && 
    error !== null && 
    'error' in error && 
    typeof (error as ApiError).error === 'object' &&
    'code' in (error as ApiError).error &&
    'message' in (error as ApiError).error
  );
};

/**
 * エラーがServerErrorかどうかを判定
 */
export const isServerError = (error: unknown): error is ServerError => {
  return (
    typeof error === 'object' && 
    error !== null && 
    'status' in error && 
    'message' in error &&
    'isServerError' in error
  );
};

/**
 * 特定のエラーコードかどうかを判定
 */
export const isErrorCode = (error: unknown, code: string): boolean => {
  return isApiError(error) && error.error.code === code;
};

/**
 * 認証エラーかどうかを判定
 */
export const isAuthError = (error: unknown): boolean => {
  return (
    isErrorCode(error, 'AUTHENTICATION_REQUIRED') ||
    isErrorCode(error, 'INVALID_TOKEN') ||
    isErrorCode(error, 'TOKEN_EXPIRED') ||
    isErrorCode(error, 'SESSION_EXPIRED')
  );
};

/**
 * 入力エラーかどうかを判定
 */
export const isInputError = (error: unknown): boolean => {
  return isErrorCode(error, 'INVALID_INPUT');
};

/**
 * エラーメッセージのユーザーフレンドリーな表示を取得
 */
export const getFriendlyErrorMessage = (error: unknown): string => {
  // APIエラーの場合
  if (isApiError(error)) {
    // 認証エラー
    if (isAuthError(error)) {
      return 'セッションの期限が切れました。再度ログインしてください。';
    }
    
    // 入力エラー
    if (isInputError(error)) {
      return error.error.message || '入力内容に誤りがあります。';
    }
    
    // その他のAPIエラー
    return error.error.message;
  }
  
  // サーバーエラー
  if (isServerError(error)) {
    if (error.status === 0) {
      return 'サーバーに接続できません。インターネット接続を確認してください。';
    }
    return error.message;
  }
  
  // 一般的なエラー
  if (error instanceof Error) {
    return error.message;
  }
  
  // 不明なエラー
  return '予期しないエラーが発生しました。後でもう一度お試しください。';
}; 