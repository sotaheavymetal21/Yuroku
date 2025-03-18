import { handleApiError, isApiError, isServerError, isErrorCode, isAuthError, isInputError, getFriendlyErrorMessage } from './error-handler';
import { ApiError, ServerError } from '@/types';

describe('Error Handler Utilities', () => {
  describe('isApiError', () => {
    it('should return true for valid ApiError', () => {
      const error: ApiError = {
        error: {
          code: 'TEST_ERROR',
          message: 'Test error message'
        }
      };
      expect(isApiError(error)).toBe(true);
    });

    it('should return false for non-ApiError', () => {
      expect(isApiError(null)).toBe(false);
      expect(isApiError(undefined)).toBe(false);
      expect(isApiError('string')).toBe(false);
      expect(isApiError(123)).toBe(false);
      expect(isApiError({})).toBe(false);
      expect(isApiError({ error: 'not an object' })).toBe(false);
      expect(isApiError({ error: {} })).toBe(false);
    });
  });

  describe('isServerError', () => {
    it('should return true for valid ServerError', () => {
      const error: ServerError = {
        status: 500,
        message: 'Internal Server Error',
        isServerError: true
      };
      expect(isServerError(error)).toBe(true);
    });

    it('should return false for non-ServerError', () => {
      expect(isServerError(null)).toBe(false);
      expect(isServerError(undefined)).toBe(false);
      expect(isServerError('string')).toBe(false);
      expect(isServerError(123)).toBe(false);
      expect(isServerError({})).toBe(false);
      expect(isServerError({ status: 500 })).toBe(false);
      expect(isServerError({ message: 'Error' })).toBe(false);
      expect(isServerError({ isServerError: true })).toBe(false);
    });
  });

  describe('isErrorCode', () => {
    it('should return true when error code matches', () => {
      const error: ApiError = {
        error: {
          code: 'MATCH_CODE',
          message: 'Error message'
        }
      };
      expect(isErrorCode(error, 'MATCH_CODE')).toBe(true);
    });

    it('should return false when error code does not match', () => {
      const error: ApiError = {
        error: {
          code: 'WRONG_CODE',
          message: 'Error message'
        }
      };
      expect(isErrorCode(error, 'MATCH_CODE')).toBe(false);
    });

    it('should return false for non-ApiError', () => {
      expect(isErrorCode(null, 'CODE')).toBe(false);
      expect(isErrorCode(undefined, 'CODE')).toBe(false);
      expect(isErrorCode('string', 'CODE')).toBe(false);
      expect(isErrorCode({}, 'CODE')).toBe(false);
    });
  });

  describe('isAuthError', () => {
    it('should return true for authentication errors', () => {
      const authErrors = [
        'AUTHENTICATION_REQUIRED',
        'INVALID_TOKEN',
        'TOKEN_EXPIRED',
        'SESSION_EXPIRED'
      ];

      authErrors.forEach(code => {
        const error: ApiError = {
          error: {
            code,
            message: 'Authentication error'
          }
        };
        expect(isAuthError(error)).toBe(true);
      });
    });

    it('should return false for non-authentication errors', () => {
      const error: ApiError = {
        error: {
          code: 'OTHER_ERROR',
          message: 'Some other error'
        }
      };
      expect(isAuthError(error)).toBe(false);
    });
  });

  describe('handleApiError', () => {
    it('should return message from ApiError', () => {
      const error: ApiError = {
        error: {
          code: 'TEST_ERROR',
          message: 'API error message'
        }
      };
      expect(handleApiError(error)).toBe('API error message');
    });

    it('should return message from ServerError', () => {
      const error: ServerError = {
        status: 500,
        message: 'Server error message',
        isServerError: true
      };
      expect(handleApiError(error)).toBe('Server error message');
    });

    it('should return message from Error object', () => {
      const error = new Error('Standard error message');
      expect(handleApiError(error)).toBe('Standard error message');
    });

    it('should return generic message for unknown error types', () => {
      expect(handleApiError(null)).toBe('予期しないエラーが発生しました。もう一度お試しください。');
      expect(handleApiError(undefined)).toBe('予期しないエラーが発生しました。もう一度お試しください。');
      expect(handleApiError('string error')).toBe('予期しないエラーが発生しました。もう一度お試しください。');
      expect(handleApiError(123)).toBe('予期しないエラーが発生しました。もう一度お試しください。');
    });
  });

  describe('getFriendlyErrorMessage', () => {
    it('should return friendly message for auth errors', () => {
      const error: ApiError = {
        error: {
          code: 'TOKEN_EXPIRED',
          message: 'Original message'
        }
      };
      expect(getFriendlyErrorMessage(error)).toBe('セッションの期限が切れました。再度ログインしてください。');
    });

    it('should return original message for input errors', () => {
      const error: ApiError = {
        error: {
          code: 'INVALID_INPUT',
          message: 'メールアドレスの形式が正しくありません'
        }
      };
      expect(getFriendlyErrorMessage(error)).toBe('メールアドレスの形式が正しくありません');
    });

    it('should handle server connection errors', () => {
      const error: ServerError = {
        status: 0,
        message: 'サーバーに接続できません',
        isServerError: true
      };
      expect(getFriendlyErrorMessage(error)).toBe('サーバーに接続できません。インターネット接続を確認してください。');
    });
  });
}); 