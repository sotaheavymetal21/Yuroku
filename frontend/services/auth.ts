import Cookies from 'js-cookie';
import { apiRequest } from './api';
import { AuthResponse, LoginRequest, RegisterRequest } from '@/types';

// トークンの有効期限（7日）
const TOKEN_EXPIRY_DAYS = 7;

// ユーザー登録
export const register = async (data: RegisterRequest): Promise<AuthResponse> => {
  const response = await apiRequest<AuthResponse>({
    method: 'POST',
    url: '/auth/register',
    data,
  });

  // トークンをCookieに保存
  if (response.data && response.data.access_token) {
    // JWTトークンをCookieに保存
    Cookies.set('token', response.data.access_token, { 
      expires: TOKEN_EXPIRY_DAYS,
      path: '/',
      secure: process.env.NODE_ENV === 'production',
      sameSite: 'strict'
    });
    // リフレッシュトークンも保存
    if (response.data.refresh_token) {
      Cookies.set('refresh_token', response.data.refresh_token, {
        expires: 30, // リフレッシュトークンの有効期限を30日に設定
        path: '/',
        secure: process.env.NODE_ENV === 'production',
        sameSite: 'strict'
      });
    }
  }

  return response;
};

// ログイン
export const login = async (data: LoginRequest): Promise<AuthResponse> => {
  const response = await apiRequest<AuthResponse>({
    method: 'POST',
    url: '/auth/login',
    data,
  });

  // トークンをCookieに保存
  if (response.data && response.data.access_token) {
    // JWTトークンをCookieに保存
    Cookies.set('token', response.data.access_token, { 
      expires: TOKEN_EXPIRY_DAYS,
      path: '/',
      secure: process.env.NODE_ENV === 'production',
      sameSite: 'strict'
    });
    // リフレッシュトークンも保存
    if (response.data.refresh_token) {
      Cookies.set('refresh_token', response.data.refresh_token, {
        expires: 30, // リフレッシュトークンの有効期限を30日に設定
        path: '/',
        secure: process.env.NODE_ENV === 'production',
        sameSite: 'strict'
      });
    }
  }

  return response;
};

// ログアウト
export const logout = (): void => {
  Cookies.remove('token');
  if (typeof window !== 'undefined') {
    window.location.href = '/auth/login';
  }
};

// 認証状態の確認
export const isAuthenticated = (): boolean => {
  const token = Cookies.get('token');
  return !!token;
};

// トークンの取得
export const getToken = (): string | undefined => {
  return Cookies.get('token');
}; 