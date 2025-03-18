import axios, { AxiosError, AxiosRequestConfig, AxiosResponse } from 'axios';
import Cookies from 'js-cookie';
import { ApiError } from '@/types';

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
    
    // 認証エラーの場合
    if (error.response?.status === 401 && originalRequest) {
      // リフレッシュトークンが存在する場合は、トークン更新を試みる
      const refreshToken = Cookies.get('refresh_token');
      
      if (refreshToken) {
        try {
          // リフレッシュトークンを使用して新しいアクセストークンを取得
          const refreshResponse = await axios.post(`${API_URL}/auth/refresh`, {
            refresh_token: refreshToken
          });
          
          // 新しいトークンを保存
          if (refreshResponse.data && refreshResponse.data.data && refreshResponse.data.data.access_token) {
            Cookies.set('token', refreshResponse.data.data.access_token, { 
              expires: 7,
              path: '/',
              sameSite: 'lax'  // Cookieの制限を緩和
            });
            
            // 新しいトークンで元のリクエストを再試行
            if (originalRequest.headers) {
              originalRequest.headers.Authorization = `Bearer ${refreshResponse.data.data.access_token}`;
            }
            return axios(originalRequest);
          }
        } catch (refreshError) {
          // リフレッシュトークンによる更新に失敗した場合はログアウト
          Cookies.remove('token');
          Cookies.remove('refresh_token');
          
          // 非同期でリダイレクト
          // ここではクライアントサイドのみの操作を確実にするため
          if (typeof window !== 'undefined') {
            setTimeout(() => {
              window.location.href = '/auth/login';
            }, 100);
          }
        }
      } else {
        // リフレッシュトークンがない場合はログアウト
        Cookies.remove('token');
        
        if (typeof window !== 'undefined') {
          setTimeout(() => {
            window.location.href = '/auth/login';
          }, 100);
        }
      }
    }
    
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
  } catch (error) {
    if (axios.isAxiosError(error) && error.response) {
      console.error('API Error:', error.response.data);
      throw error.response.data;
    }
    console.error('Unexpected error:', error);
    throw error;
  }
};

export default api; 