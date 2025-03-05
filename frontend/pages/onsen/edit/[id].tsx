import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/router';
import { FaSave, FaArrowLeft } from 'react-icons/fa';
import Layout from '@/components/layout/Layout';
import Button from '@/components/common/Button';
import InputField from '@/components/common/InputField';
import TextareaField from '@/components/common/TextareaField';
import SelectField from '@/components/common/SelectField';
import Loading from '@/components/common/Loading';
import ErrorMessage from '@/components/common/ErrorMessage';
import { getOnsenLog, updateOnsenLog } from '@/services/onsenLog';
import { OnsenLog, OnsenLogUpdateRequest } from '@/types';
import { useAuth } from '@/contexts/AuthContext';

// 温泉の種類のオプション
const springTypeOptions = [
  { value: '', label: '選択してください' },
  { value: '単純温泉', label: '単純温泉' },
  { value: '塩化物泉', label: '塩化物泉' },
  { value: '炭酸水素塩泉', label: '炭酸水素塩泉' },
  { value: '硫酸塩泉', label: '硫酸塩泉' },
  { value: '二酸化炭素泉', label: '二酸化炭素泉' },
  { value: '含鉄泉', label: '含鉄泉' },
  { value: '酸性泉', label: '酸性泉' },
  { value: '含よう素泉', label: '含よう素泉' },
  { value: '硫黄泉', label: '硫黄泉' },
  { value: '放射能泉', label: '放射能泉' },
  { value: 'その他', label: 'その他' },
];

// 温泉の特徴のオプション
const featureOptions = [
  { value: '露天風呂', label: '露天風呂' },
  { value: '貸切風呂', label: '貸切風呂' },
  { value: '岩風呂', label: '岩風呂' },
  { value: '檜風呂', label: '檜風呂' },
  { value: '混浴', label: '混浴' },
  { value: '日帰り入浴', label: '日帰り入浴' },
  { value: '温泉街', label: '温泉街' },
  { value: '秘湯', label: '秘湯' },
  { value: '景色が良い', label: '景色が良い' },
  { value: '食事が美味しい', label: '食事が美味しい' },
];

const EditOnsenLogPage: React.FC = () => {
  const router = useRouter();
  const { id } = router.query;
  const { isLoggedIn, loading: authLoading } = useAuth();
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');
  const [formData, setFormData] = useState<OnsenLogUpdateRequest>({
    name: '',
    location: '',
    spring_type: '',
    features: [],
    visit_date: '',
    rating: 0,
    comment: '',
  });

  // 認証チェック
  useEffect(() => {
    if (!authLoading && !isLoggedIn) {
      router.push('/auth/login');
    }
  }, [isLoggedIn, authLoading, router]);

  // データ取得
  useEffect(() => {
    if (isLoggedIn && id) {
      fetchOnsenLog();
    }
  }, [isLoggedIn, id]);

  const fetchOnsenLog = async () => {
    setLoading(true);
    setError('');
    
    try {
      const data = await getOnsenLog(id as string);
      setFormData({
        name: data.name,
        location: data.location || '',
        spring_type: data.spring_type || '',
        features: data.features || [],
        visit_date: data.visit_date,
        rating: data.rating || 0,
        comment: data.comment || '',
      });
    } catch (err) {
      setError('温泉メモの取得に失敗しました。もう一度お試しください。');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setFormData({
      ...formData,
      [name]: value,
    });
  };

  const handleRatingChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = parseInt(e.target.value);
    setFormData({
      ...formData,
      rating: isNaN(value) ? 0 : Math.max(0, Math.min(5, value)),
    });
  };

  const handleFeatureChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { value, checked } = e.target;
    
    if (checked) {
      setFormData({
        ...formData,
        features: [...(formData.features || []), value],
      });
    } else {
      setFormData({
        ...formData,
        features: (formData.features || []).filter(feature => feature !== value),
      });
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!formData.name) {
      setError('温泉名は必須です。');
      return;
    }
    
    if (!formData.visit_date) {
      setError('訪問日は必須です。');
      return;
    }
    
    setSubmitting(true);
    setError('');
    
    try {
      await updateOnsenLog(id as string, formData);
      router.push(`/onsen/${id}`);
    } catch (err) {
      setError('更新に失敗しました。もう一度お試しください。');
      console.error(err);
      setSubmitting(false);
    }
  };

  if (authLoading || loading) {
    return (
      <Layout title="読み込み中... | 湯録 (Yuroku)">
        <div className="flex justify-center items-center h-64">
          <Loading size="large" text="読み込み中..." />
        </div>
      </Layout>
    );
  }

  return (
    <Layout title="温泉メモの編集 | 湯録 (Yuroku)">
      <div className="mb-4">
        <Button
          variant="outline"
          onClick={() => router.back()}
          icon={<FaArrowLeft />}
        >
          戻る
        </Button>
      </div>

      <div className="bg-white rounded-lg shadow-md overflow-hidden">
        <div className="p-6">
          <h1 className="text-2xl font-bold text-gray-800 mb-6">温泉メモの編集</h1>
          
          {error && <ErrorMessage message={error} className="mb-6" />}
          
          <form onSubmit={handleSubmit}>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
              <InputField
                label="温泉名 *"
                name="name"
                value={formData.name}
                onChange={handleChange}
                required
              />
              
              <InputField
                label="場所"
                name="location"
                value={formData.location}
                onChange={handleChange}
                placeholder="例: 北海道札幌市"
              />
              
              <InputField
                label="訪問日 *"
                name="visit_date"
                type="date"
                value={formData.visit_date}
                onChange={handleChange}
                required
              />
              
              <SelectField
                label="温泉の種類"
                name="spring_type"
                value={formData.spring_type}
                onChange={handleChange}
                options={springTypeOptions}
              />
              
              <div className="md:col-span-2">
                <label className="block text-gray-700 font-medium mb-2">
                  評価
                </label>
                <div className="flex items-center">
                  <input
                    type="range"
                    name="rating"
                    min="0"
                    max="5"
                    step="0.5"
                    value={formData.rating}
                    onChange={handleRatingChange}
                    className="w-full max-w-xs mr-4"
                  />
                  <span className="text-lg font-medium">{formData.rating}</span>
                </div>
              </div>
            </div>
            
            <div className="mb-6">
              <label className="block text-gray-700 font-medium mb-2">
                特徴
              </label>
              <div className="grid grid-cols-2 md:grid-cols-4 gap-2">
                {featureOptions.map((option) => (
                  <div key={option.value} className="flex items-center">
                    <input
                      type="checkbox"
                      id={`feature-${option.value}`}
                      name="features"
                      value={option.value}
                      checked={(formData.features || []).includes(option.value)}
                      onChange={handleFeatureChange}
                      className="mr-2"
                    />
                    <label htmlFor={`feature-${option.value}`}>
                      {option.label}
                    </label>
                  </div>
                ))}
              </div>
            </div>
            
            <div className="mb-6">
              <TextareaField
                label="コメント"
                name="comment"
                value={formData.comment || ''}
                onChange={handleChange}
                rows={5}
                placeholder="温泉の感想や思い出を記録しましょう..."
              />
            </div>
            
            <div className="flex justify-end">
              <Button
                type="submit"
                disabled={submitting}
                icon={<FaSave />}
              >
                {submitting ? '保存中...' : '保存する'}
              </Button>
            </div>
          </form>
        </div>
      </div>
    </Layout>
  );
};

export default EditOnsenLogPage; 