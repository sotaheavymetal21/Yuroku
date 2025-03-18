import Cookies from 'js-cookie';
import { apiRequest } from './api';
import { AuthResponse, LoginRequest, RegisterRequest } from '@/types';

// トークンの有効期限（7日）
const TOKEN_EXPIRY_DAYS = 7;

// Cookieの設定オプション
const cookieOptions = {
  expires: TOKEN_EXPIRY_DAYS,
  path: '/',
  sameSite: 'lax', // strictからlaxに変更して他ドメインとの連携を許可
  secure: process.env.NODE_ENV === 'production' // 開発環境ではHTTPでも使えるように
};

// トークンをCookieに保存する関数
const saveTokensToCookie = (accessToken: string, refreshToken?: string): void => {
  // JWTトークンをCookieに保存
  Cookies.set('token', accessToken, cookieOptions);
  
  // リフレッシュトークンも保存
  if (refreshToken) {
    Cookies.set('refresh_token', refreshToken, {
      ...cookieOptions,
      expires: 30 // リフレッシュトークンの有効期限を30日に設定
    });
  }

  // ローカルストレージにもバックアップ（開発環境用）
  if (typeof window !== 'undefined' && process.env.NODE_ENV !== 'production') {
    localStorage.setItem('debug_token', accessToken);
    if (refreshToken) {
      localStorage.setItem('debug_refresh_token', refreshToken);
    }
  }
};

// ユーザー登録
export const register = async (data: RegisterRequest): Promise<AuthResponse> => {
  try {
    const response = await apiRequest<AuthResponse>({
      method: 'POST',
      url: '/auth/register',
      data,
    });

    // トークンをCookieに保存
    if (response.data && response.data.access_token) {
      saveTokensToCookie(
        response.data.access_token, 
        response.data.refresh_token
      );
    }

    return response;
  } catch (error) {
    console.error('登録エラー:', error);
    throw error;
  }
};

// ログイン
export const login = async (data: LoginRequest): Promise<AuthResponse> => {
  try {
    const response = await apiRequest<AuthResponse>({
      method: 'POST',
      url: '/auth/login',
      data,
    });

    // トークンをCookieに保存
    if (response.data && response.data.access_token) {
      saveTokensToCookie(
        response.data.access_token, 
        response.data.refresh_token
      );
    }

    return response;
  } catch (error) {
    console.error('ログインエラー:', error);
    throw error;
  }
};

// ログアウト
export const logout = (): void => {
  // サーバーにログアウトリクエストを送信
  try {
    apiRequest({
      method: 'POST',
      url: '/auth/logout',
    }).catch(err => console.error('ログアウトAPI呼び出しエラー:', err));
  } catch (error) {
    console.error('ログアウトリクエストエラー:', error);
  }

  // Cookie削除
  Cookies.remove('token', { path: '/' });
  Cookies.remove('refresh_token', { path: '/' });

  // ローカルストレージのクリーンアップ（開発環境用）
  if (typeof window !== 'undefined') {
    localStorage.removeItem('debug_token');
    localStorage.removeItem('debug_refresh_token');
  }
};

// 認証状態の確認
export const isAuthenticated = (): boolean => {
  // Cookieとローカルストレージの両方をチェック
  const cookieToken = Cookies.get('token');
  
  if (cookieToken) {
    return true;
  }
  
  // バックアップチェック（開発環境用）
  if (typeof window !== 'undefined' && process.env.NODE_ENV !== 'production') {
    const localToken = localStorage.getItem('debug_token');
    if (localToken) {
      // ローカルストレージからトークンを発見した場合、Cookieに復元
      Cookies.set('token', localToken, cookieOptions);
      return true;
    }
  }
  
  return false;
};

// トークンの取得
export const getToken = (): string | undefined => {
  return Cookies.get('token');
};