import axios, { AxiosError, AxiosRequestConfig, AxiosResponse } from 'axios';
import Cookies from 'js-cookie';
import { ApiError, ServerError } from '@/types';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

// Axiosインスタンスの作成
const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true,  // CORS対応のためにクッキーを送信
});

// リクエストインターセプター
api.interceptors.request.use(
  (config) => {
    const token = Cookies.get('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    
    // Content-Typeがmultipart/form-dataの場合はヘッダーを設定しない
    if (config.data instanceof FormData) {
      delete config.headers['Content-Type'];
    }
    
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// レスポンスインターセプター
api.interceptors.response.use(
  (response: AxiosResponse) => {
    return response;
  },
  async (error: AxiosError<ApiError>) => {
    const originalRequest = error.config;
    
    // 認証エラーの場合（tokenの期限切れ）
    if (error.response?.status === 401 && originalRequest && !originalRequest.headers?._retry) {
      // リトライフラグを設定して無限ループを防止
      originalRequest.headers = originalRequest.headers || {};
      originalRequest.headers._retry = true;
      
      // リフレッシュトークンが存在する場合は、トークン更新を試みる
      const refreshToken = Cookies.get('refresh_token');
      
      if (refreshToken) {
        try {
          // リフレッシュトークンを使用して新しいアクセストークンを取得
          const refreshResponse = await axios.post(`${API_URL}/auth/refresh`, {
            refresh_token: refreshToken
          }, {
            // リフレッシュトークンのリクエストには、認証ヘッダーをつけない
            headers: {
              'Content-Type': 'application/json'
            }
          });
          
          // 新しいトークンを保存
          if (refreshResponse.data?.data?.access_token) {
            const newAccessToken = refreshResponse.data.data.access_token;
            const newRefreshToken = refreshResponse.data.data.refresh_token;
            
            // HTTPOnlyクッキーを使用できない場合の次善策
            Cookies.set('token', newAccessToken, { 
              expires: 1/96, // 15分の有効期限
              path: '/',
              sameSite: 'strict',  // CSRFを防止するために厳格に設定
              secure: typeof window !== 'undefined' && window.location.protocol === 'https:'  // HTTPS接続時のみ
            });
            
            if (newRefreshToken) {
              Cookies.set('refresh_token', newRefreshToken, { 
                expires: 7,  // 7日間
                path: '/',
                sameSite: 'strict',
                secure: typeof window !== 'undefined' && window.location.protocol === 'https:'
              });
            }
            
            // 新しいトークンで元のリクエストを再試行
            if (originalRequest.headers) {
              originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;
            }
            return axios(originalRequest);
          }
        } catch (refreshError) {
          // リフレッシュトークンによる更新に失敗した場合は、すべてのトークンを削除
          Cookies.remove('token');
          Cookies.remove('refresh_token');
          
          // ログアウト状態に遷移
          if (typeof window !== 'undefined') {
            // 再ログインページへリダイレクト前にカスタムイベント発行
            window.dispatchEvent(new CustomEvent('auth:sessionExpired'));
            
            // 直接リダイレクトではなく、ユーザーに通知を表示
            setTimeout(() => {
              window.location.href = '/auth/login?expired=true';
            }, 100);
          }
          
          return Promise.reject({
            error: {
              code: 'SESSION_EXPIRED',
              message: 'セッションの有効期限が切れました。再度ログインしてください。'
            }
          });
        }
      }
      
      // リフレッシュトークンがない場合
      Cookies.remove('token');
      
      if (typeof window !== 'undefined') {
        window.dispatchEvent(new CustomEvent('auth:tokenMissing'));
        setTimeout(() => {
          window.location.href = '/auth/login';
        }, 100);
      }
      
      return Promise.reject({
        error: {
          code: 'AUTHENTICATION_REQUIRED',
          message: '認証が必要です。ログインしてください。'
        }
      });
    }
    
    // その他のエラー
    return Promise.reject(error);
  }
);

// APIリクエスト関数
export const apiRequest = async <T>(
  config: AxiosRequestConfig
): Promise<T> => {
  try {
    const response = await api(config);
    return response.data;
  } catch (error: unknown) {
    if (axios.isAxiosError(error)) {
      // レスポンスがある場合（サーバーからのエラー）
      if (error.response) {
        console.error('API Error:', error.response.data);
        throw error.response.data;
      } 
      // レスポンスがない場合（ネットワークエラーなど）
      else if (error.request) {
        console.error('Network Error:', error.message);
        const serverError: ServerError = {
          status: 0,
          message: 'サーバーに接続できません。ネットワーク接続を確認してください。',
          isServerError: true
        };
        throw serverError;
      }
    }
    // その他の予期しないエラー
    console.error('Unexpected error:', error);
    const unexpectedError: ServerError = {
      status: 500,
      message: '予期しないエラーが発生しました。しばらく経ってからもう一度お試しください。',
      isServerError: true
    };
    throw unexpectedError;
  }
};

export default api; 