import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/router';
import Image from 'next/image';
import { FaStar, FaMapMarkerAlt, FaCalendarAlt, FaWater, FaEdit, FaTrash, FaArrowLeft } from 'react-icons/fa';
import Layout from '@/components/layout/Layout';
import Button from '@/components/common/Button';
import Loading from '@/components/common/Loading';
import ErrorMessage from '@/components/common/ErrorMessage';
import { getOnsenLog, deleteOnsenLog } from '@/services/onsenLog';
import { OnsenLog } from '@/types';
import { useAuth } from '@/contexts/AuthContext';

const OnsenDetailPage: React.FC = () => {
  const router = useRouter();
  const { id } = router.query;
  const { isLoggedIn, loading: authLoading } = useAuth();
  const [onsenLog, setOnsenLog] = useState<OnsenLog | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [deleteConfirm, setDeleteConfirm] = useState(false);

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
      setOnsenLog(data);
    } catch (err) {
      setError('温泉メモの取得に失敗しました。もう一度お試しください。');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  // 日付をフォーマット
  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('ja-JP', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  };

  // 評価を星で表示
  const renderRating = (rating: number | undefined) => {
    if (rating === undefined) return null;
    
    return (
      <div className="flex items-center">
        {[...Array(5)].map((_, i) => (
          <FaStar
            key={i}
            className={`w-5 h-5 ${
              i < rating ? 'text-yellow-400' : 'text-gray-300'
            }`}
          />
        ))}
        <span className="ml-2 text-lg">{rating}</span>
      </div>
    );
  };

  const handleDelete = async () => {
    if (!deleteConfirm) {
      setDeleteConfirm(true);
      return;
    }
    
    try {
      await deleteOnsenLog(id as string);
      router.push('/onsen');
    } catch (err) {
      setError('削除に失敗しました。もう一度お試しください。');
      console.error(err);
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

  if (error) {
    return (
      <Layout title="エラー | 湯録 (Yuroku)">
        <div className="mb-4">
          <Button
            variant="outline"
            onClick={() => router.back()}
            icon={<FaArrowLeft />}
          >
            戻る
          </Button>
        </div>
        <ErrorMessage message={error} />
      </Layout>
    );
  }

  if (!onsenLog) {
    return (
      <Layout title="温泉メモが見つかりません | 湯録 (Yuroku)">
        <div className="mb-4">
          <Button
            variant="outline"
            onClick={() => router.back()}
            icon={<FaArrowLeft />}
          >
            戻る
          </Button>
        </div>
        <div className="bg-white rounded-lg shadow-md p-6 text-center">
          <p className="text-gray-500 mb-4">温泉メモが見つかりませんでした。</p>
          <Button
            onClick={() => router.push('/onsen')}
          >
            一覧に戻る
          </Button>
        </div>
      </Layout>
    );
  }

  return (
    <Layout title={`${onsenLog.name} | 湯録 (Yuroku)`}>
      <div className="mb-4 flex justify-between items-center">
        <Button
          variant="outline"
          onClick={() => router.back()}
          icon={<FaArrowLeft />}
        >
          戻る
        </Button>
        <div className="flex space-x-2">
          <Button
            variant="outline"
            onClick={() => router.push(`/onsen/edit/${onsenLog.id}`)}
            icon={<FaEdit />}
          >
            編集
          </Button>
          <Button
            variant="danger"
            onClick={handleDelete}
            icon={<FaTrash />}
          >
            {deleteConfirm ? '削除確認' : '削除'}
          </Button>
        </div>
      </div>

      <div className="bg-white rounded-lg shadow-md overflow-hidden">
        <div className="p-6">
          <h1 className="text-3xl font-bold text-gray-800 mb-4">{onsenLog.name}</h1>
          
          <div className="mb-6">
            {renderRating(onsenLog.rating)}
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
            <div>
              {onsenLog.location && (
                <div className="flex items-center text-gray-600 mb-4">
                  <FaMapMarkerAlt className="mr-2 text-onsen" />
                  <span className="text-lg">{onsenLog.location}</span>
                </div>
              )}
              
              <div className="flex items-center text-gray-600 mb-4">
                <FaCalendarAlt className="mr-2 text-onsen" />
                <span className="text-lg">{formatDate(onsenLog.visit_date)}</span>
              </div>
              
              {onsenLog.spring_type && (
                <div className="flex items-center text-gray-600 mb-4">
                  <FaWater className="mr-2 text-onsen" />
                  <span className="text-lg">{onsenLog.spring_type}</span>
                </div>
              )}
            </div>
            
            {/* 画像があれば表示 */}
            {/* 実装予定 */}
          </div>
          
          {onsenLog.features && onsenLog.features.length > 0 && (
            <div className="mb-6">
              <h2 className="text-xl font-semibold mb-2">特徴</h2>
              <div className="flex flex-wrap gap-2">
                {onsenLog.features.map((feature, index) => (
                  <span
                    key={index}
                    className="inline-block bg-gray-100 rounded-full px-4 py-2 text-sm font-semibold text-gray-700"
                  >
                    {feature}
                  </span>
                ))}
              </div>
            </div>
          )}
          
          {onsenLog.comment && (
            <div className="mb-6">
              <h2 className="text-xl font-semibold mb-2">コメント</h2>
              <div className="bg-gray-50 p-4 rounded-lg whitespace-pre-wrap">
                {onsenLog.comment}
              </div>
            </div>
          )}
          
          <div className="text-gray-500 text-sm mt-8">
            <p>作成日: {new Date(onsenLog.created_at).toLocaleString('ja-JP')}</p>
            <p>更新日: {new Date(onsenLog.updated_at).toLocaleString('ja-JP')}</p>
          </div>
        </div>
      </div>
    </Layout>
  );
};

export default OnsenDetailPage; 