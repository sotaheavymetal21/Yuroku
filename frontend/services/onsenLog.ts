import { apiRequest } from './api';
import {
  OnsenLog,
  OnsenLogCreateRequest,
  OnsenLogUpdateRequest,
  OnsenLogResponse,
  OnsenLogFilter,
  PaginationParams,
} from '@/types';

// 温泉メモの作成
export const createOnsenLog = async (
  data: OnsenLogCreateRequest
): Promise<OnsenLogResponse> => {
  return apiRequest<OnsenLogResponse>({
    method: 'POST',
    url: '/onsen_logs',
    data,
  });
};

// 温泉メモの一覧取得
export const getOnsenLogs = async (
  pagination?: PaginationParams
): Promise<OnsenLog[]> => {
  const params = pagination
    ? { page: pagination.page, limit: pagination.limit }
    : {};

  return apiRequest<OnsenLog[]>({
    method: 'GET',
    url: '/onsen_logs',
    params,
  });
};

// 温泉メモの詳細取得
export const getOnsenLog = async (id: string): Promise<OnsenLog> => {
  return apiRequest<OnsenLog>({
    method: 'GET',
    url: `/onsen_logs/${id}`,
  });
};

// 温泉メモの更新
export const updateOnsenLog = async (
  id: string,
  data: OnsenLogUpdateRequest
): Promise<OnsenLogResponse> => {
  return apiRequest<OnsenLogResponse>({
    method: 'PUT',
    url: `/onsen_logs/${id}`,
    data,
  });
};

// 温泉メモの削除
export const deleteOnsenLog = async (
  id: string
): Promise<{ message: string }> => {
  return apiRequest<{ message: string }>({
    method: 'DELETE',
    url: `/onsen_logs/${id}`,
  });
};

// 温泉メモのフィルタリング
export const filterOnsenLogs = async (
  filter: OnsenLogFilter,
  pagination?: PaginationParams
): Promise<OnsenLog[]> => {
  const params = {
    ...filter,
    ...(pagination ? { page: pagination.page, limit: pagination.limit } : {}),
  };

  return apiRequest<OnsenLog[]>({
    method: 'GET',
    url: '/onsen_logs',
    params,
  });
};

// 温泉メモのエクスポート（JSON形式）
export const exportOnsenLogsAsJson = async (): Promise<Blob> => {
  const response = await apiRequest<Blob>({
    method: 'GET',
    url: '/onsen_logs/export',
    params: { format: 'json' },
    responseType: 'blob',
  });

  return response;
};

// 温泉メモのエクスポート（CSV形式）
export const exportOnsenLogsAsCsv = async (): Promise<Blob> => {
  const response = await apiRequest<Blob>({
    method: 'GET',
    url: '/onsen_logs/export',
    params: { format: 'csv' },
    responseType: 'blob',
  });

  return response;
}; 