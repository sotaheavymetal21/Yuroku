// ユーザー関連の型定義
export interface User {
  id: string;
  email: string;
  createdAt: string;
  updatedAt: string;
}

export interface AuthResponse {
  user_id: string;
  token: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
}

// 温泉メモ関連の型定義
export interface OnsenLog {
  id: string;
  user_id: string;
  name: string;
  location?: string;
  spring_type?: string;
  features?: string[];
  visit_date: string;
  rating?: number;
  comment?: string;
  created_at: string;
  updated_at: string;
  // フロントエンド表示用
  _id?: string;
  waterType?: string;
}

export interface OnsenLogCreateRequest {
  name: string;
  location?: string;
  spring_type?: string;
  features?: string[];
  visit_date: string;
  rating?: number;
  comment?: string;
}

export interface OnsenLogUpdateRequest {
  name?: string;
  location?: string;
  spring_type?: string;
  features?: string[];
  visit_date?: string;
  rating?: number;
  comment?: string;
}

export interface OnsenLogResponse {
  id: string;
  message: string;
}

// 温泉画像関連の型定義
export interface OnsenImage {
  image_id: string;
  image_url: string;
}

export interface OnsenImageUploadResponse {
  image_id: string;
  image_url: string;
}

// フィルタリング関連の型定義
export interface OnsenLogFilter {
  keyword?: string;
  spring_type?: string;
  rating?: number;
  start_date?: string;
  end_date?: string;
  // フロントエンド用の追加フィールド
  name?: string;
  location?: string;
  minRating?: number;
  maxRating?: number;
  fromDate?: string;
  toDate?: string;
}

// ページネーション関連の型定義
export interface PaginationParams {
  page: number;
  limit: number;
  // フロントエンド用の追加フィールド
  total?: number;
  sortBy?: string;
  sortDirection?: 'asc' | 'desc';
}

// APIレスポンスの共通型
export interface ApiResponse<T> {
  data: T;
  message?: string;
  total?: number;
}

export interface ApiError {
  error: {
    code: string;
    message: string;
  };
} 