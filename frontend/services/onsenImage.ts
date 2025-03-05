import { apiRequest } from './api';
import { OnsenImage, OnsenImageUploadResponse } from '@/types';

// 温泉画像のアップロード
export const uploadOnsenImage = async (
  onsenId: string,
  file: File
): Promise<OnsenImageUploadResponse> => {
  const formData = new FormData();
  formData.append('onsen_id', onsenId);
  formData.append('file', file);

  return apiRequest<OnsenImageUploadResponse>({
    method: 'POST',
    url: '/onsen_images',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
};

// 温泉画像の取得
export const getOnsenImages = async (
  onsenId: string
): Promise<OnsenImage[]> => {
  return apiRequest<OnsenImage[]>({
    method: 'GET',
    url: `/onsen_images/${onsenId}`,
  });
};

// 温泉画像の削除
export const deleteOnsenImage = async (
  imageId: string
): Promise<{ message: string }> => {
  return apiRequest<{ message: string }>({
    method: 'DELETE',
    url: `/onsen_images/${imageId}`,
  });
}; 