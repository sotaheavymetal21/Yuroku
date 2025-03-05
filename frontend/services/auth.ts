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
  if (response.token) {
    Cookies.set('token', response.token, { expires: TOKEN_EXPIRY_DAYS });
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
  if (response.token) {
    Cookies.set('token', response.token, { expires: TOKEN_EXPIRY_DAYS });
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
  return !!Cookies.get('token');
};

// トークンの取得
export const getToken = (): string | undefined => {
  return Cookies.get('token');
}; 