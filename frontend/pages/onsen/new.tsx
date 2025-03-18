import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/router';
import { FaSave, FaTimes, FaUpload } from 'react-icons/fa';
import Layout from '@/components/layout/Layout';
import Card from '@/components/common/Card';
import Button from '@/components/common/Button';
import InputField from '@/components/common/InputField';
import TextareaField from '@/components/common/TextareaField';
import ErrorMessage from '@/components/common/ErrorMessage';
import Loading from '@/components/common/Loading';
import { createOnsenLog } from '@/services/onsenLog';
import { uploadOnsenImage } from '@/services/onsenImage';
import { useAuth } from '@/contexts/AuthContext';

const NewOnsenLogPage: React.FC = () => {
  const router = useRouter();
  const { isLoggedIn, loading: authLoading } = useAuth();
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState('');
  const [images, setImages] = useState<File[]>([]);
  const [previewUrls, setPreviewUrls] = useState<string[]>([]);
  
  // フォームの状態
  const [formData, setFormData] = useState({
    name: '',
    location: '',
    visitDate: new Date().toISOString().split('T')[0],
    rating: 3,
    waterType: '',
    price: 0,
    facilities: '',
    comment: '',
  });

  // 認証チェック
  useEffect(() => {
    if (!authLoading && !isLoggedIn) {
      router.push('/auth/login');
    }
  }, [isLoggedIn, authLoading, router]);

  // 入力変更ハンドラー
  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setFormData({
      ...formData,
      [name]: name === 'rating' || name === 'price' ? Number(value) : value,
    });
  };

  // 画像アップロードハンドラー
  const handleImageUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      const newFiles = Array.from(e.target.files);
      setImages([...images, ...newFiles]);
      
      // プレビューURLの生成
      const newPreviewUrls = newFiles.map(file => URL.createObjectURL(file));
      setPreviewUrls([...previewUrls, ...newPreviewUrls]);
    }
  };

  // 画像削除ハンドラー
  const handleRemoveImage = (index: number) => {
    const newImages = [...images];
    newImages.splice(index, 1);
    setImages(newImages);
    
    const newPreviewUrls = [...previewUrls];
    URL.revokeObjectURL(newPreviewUrls[index]);
    newPreviewUrls.splice(index, 1);
    setPreviewUrls(newPreviewUrls);
  };

  // フォーム送信ハンドラー
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setIsSubmitting(true);

    try {
      // バリデーション
      if (!formData.name || !formData.location || !formData.waterType || !formData.visitDate) {
        setError('必須項目を入力してください');
        setIsSubmitting(false);
        return;
      }

      // バックエンドのAPI要件に合わせてデータを変換
      const requestData = {
        name: formData.name,
        location: formData.location,
        spring_type: formData.waterType, // waterTypeをspring_typeとして送信
        visit_date: formData.visitDate, // visitDateをvisit_dateとして送信
        rating: formData.rating,
        comment: formData.comment,
      };

      console.log('送信データ:', requestData);

      // 温泉メモの作成
      const response = await createOnsenLog(requestData);
      console.log('レスポンス:', response);

      if (!response || !response.data || !response.data.id) {
        throw new Error('無効なレスポンス: IDが取得できませんでした');
      }

      const onsenId = response.data.id;
      console.log('作成された温泉ID:', onsenId);

      // 画像のアップロード（複数ある場合は順次処理）
      if (images.length > 0) {
        try {
          for (const image of images) {
            console.log(`画像アップロード: ${image.name} (${onsenId})`);
            await uploadOnsenImage(onsenId, image);
          }
        } catch (uploadErr) {
          console.error('画像アップロードエラー:', uploadErr);
          // 画像アップロードが失敗しても温泉メモは作成されているので、エラーメッセージを表示するだけ
          setError('温泉メモは作成されましたが、画像のアップロードに失敗しました。');
          // 詳細ページに移動
          setTimeout(() => {
            router.push(`/onsen/${onsenId}`);
          }, 1000);
          return;
        }
      }

      // 認証状態を更新（必要に応じて）
      // useAuth().refreshAuthState();

      // 成功したら詳細ページへリダイレクト（遅延させてステート更新を確実に）
      setTimeout(() => {
        router.push(`/onsen/${onsenId}`);
      }, 500);
    } catch (err) {
      console.error('温泉メモ作成エラー:', err);
      
      // エラーメッセージを表示
      if (err && typeof err === 'object' && 'error' in err) {
        const apiError = err as { error: { message: string } };
        setError(apiError.error.message || '温泉メモの作成に失敗しました。もう一度お試しください。');
      } else {
        setError('温泉メモの作成に失敗しました。もう一度お試しください。');
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  if (authLoading) {
    return (
      <Layout title="読み込み中... | 湯録 (Yuroku)">
        <div className="flex justify-center items-center h-64">
          <Loading size="large" text="読み込み中..." />
        </div>
      </Layout>
    );
  }

  return (
    <Layout title="新規温泉メモ作成 | 湯録 (Yuroku)">
      <div className="max-w-3xl mx-auto">
        <h1 className="text-2xl font-bold mb-6">新規温泉メモ作成</h1>
        
        {error && <ErrorMessage message={error} className="mb-4" />}
        
        <form onSubmit={handleSubmit}>
          <Card className="mb-6">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <InputField
                label="温泉名 *"
                name="name"
                value={formData.name}
                onChange={handleChange}
                placeholder="例: ○○温泉"
                required
              />
              
              <InputField
                label="場所 *"
                name="location"
                value={formData.location}
                onChange={handleChange}
                placeholder="例: 北海道札幌市"
                required
              />
              
              <InputField
                label="訪問日 *"
                name="visitDate"
                type="date"
                value={formData.visitDate}
                onChange={handleChange}
                required
              />
              
              <InputField
                label="評価 (1-5) *"
                name="rating"
                type="number"
                min={1}
                max={5}
                value={formData.rating}
                onChange={handleChange}
                required
              />
              
              <InputField
                label="泉質 *"
                name="waterType"
                value={formData.waterType}
                onChange={handleChange}
                placeholder="例: 硫黄泉"
                required
              />
              
              <InputField
                label="料金 (円)"
                name="price"
                type="number"
                min={0}
                value={formData.price}
                onChange={handleChange}
                placeholder="例: 800"
              />
            </div>
            
            <div className="mt-4">
              <InputField
                label="設備・サービス"
                name="facilities"
                value={formData.facilities}
                onChange={handleChange}
                placeholder="例: サウナ、露天風呂、休憩室"
              />
            </div>
            
            <div className="mt-4">
              <TextareaField
                label="感想・メモ"
                name="comment"
                value={formData.comment}
                onChange={handleChange}
                placeholder="温泉の感想や思い出を自由に記録しましょう"
                rows={5}
              />
            </div>
          </Card>
          
          <Card title="写真" className="mb-6">
            <div className="mb-4">
              <label className="block text-sm font-medium text-gray-700 mb-1">
                温泉の写真
              </label>
              <div className="flex items-center">
                <label className="cursor-pointer bg-white border border-gray-300 rounded-md px-3 py-2 hover:bg-gray-50 transition-colors">
                  <span className="flex items-center text-gray-600">
                    <FaUpload className="mr-2" /> 写真を選択
                  </span>
                  <input
                    type="file"
                    accept="image/*"
                    multiple
                    onChange={handleImageUpload}
                    className="hidden"
                  />
                </label>
                <span className="ml-3 text-sm text-gray-500">
                  JPG, PNG, GIF形式（最大5MB）
                </span>
              </div>
            </div>
            
            {previewUrls.length > 0 && (
              <div className="grid grid-cols-2 md:grid-cols-3 gap-4 mt-4">
                {previewUrls.map((url, index) => (
                  <div key={index} className="relative">
                    <img
                      src={url}
                      alt={`プレビュー ${index + 1}`}
                      className="w-full h-32 object-cover rounded-md"
                    />
                    <button
                      type="button"
                      onClick={() => handleRemoveImage(index)}
                      className="absolute top-1 right-1 bg-red-500 text-white rounded-full p-1 hover:bg-red-600 transition-colors"
                    >
                      <FaTimes size={14} />
                    </button>
                  </div>
                ))}
              </div>
            )}
          </Card>
          
          <div className="flex justify-end space-x-4">
            <Button
              type="button"
              variant="outline"
              onClick={() => router.back()}
            >
              キャンセル
            </Button>
            <Button
              type="submit"
              isLoading={isSubmitting}
              loadingText="保存中..."
              icon={<FaSave />}
            >
              保存する
            </Button>
          </div>
        </form>
      </div>
    </Layout>
  );
};

export default NewOnsenLogPage; 