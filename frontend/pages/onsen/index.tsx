import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/router';
import { FaPlus, FaSearch, FaFilter, FaDownload, FaSortAmountDown, FaSortAmountUp } from 'react-icons/fa';
import Layout from '@/components/layout/Layout';
import Card from '@/components/common/Card';
import Button from '@/components/common/Button';
import InputField from '@/components/common/InputField';
import Loading from '@/components/common/Loading';
import ErrorMessage from '@/components/common/ErrorMessage';
import OnsenCard from '@/components/onsen/OnsenCard';
import { getOnsenLogs, exportOnsenLogsAsJson, exportOnsenLogsAsCsv } from '@/services/onsenLog';
import { OnsenLog, PaginationParams, OnsenLogFilter } from '@/types';
import { useAuth } from '@/contexts/AuthContext';

const OnsenLogsPage: React.FC = () => {
  const router = useRouter();
  const { isLoggedIn, loading: authLoading } = useAuth();
  const [onsenLogs, setOnsenLogs] = useState<OnsenLog[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [searchTerm, setSearchTerm] = useState('');
  const [showFilters, setShowFilters] = useState(false);
  const [sortField, setSortField] = useState<'visitDate' | 'rating'>('visitDate');
  const [sortDirection, setSortDirection] = useState<'asc' | 'desc'>('desc');
  
  // ページネーション
  const [pagination, setPagination] = useState<PaginationParams>({
    page: 1,
    limit: 10,
    total: 0,
  });
  
  // フィルター
  const [filters, setFilters] = useState<OnsenLogFilter>({
    name: '',
    location: '',
    minRating: 0,
    maxRating: 5,
    fromDate: '',
    toDate: '',
  });

  // 認証チェック
  useEffect(() => {
    if (!authLoading && !isLoggedIn) {
      router.push('/auth/login');
    }
  }, [isLoggedIn, authLoading, router]);

  // データ取得
  useEffect(() => {
    if (isLoggedIn) {
      fetchOnsenLogs();
    }
  }, [isLoggedIn, pagination.page, sortField, sortDirection]);

  const fetchOnsenLogs = async () => {
    setLoading(true);
    setError('');
    
    try {
      const response = await getOnsenLogs({
        page: pagination.page,
        limit: pagination.limit,
        sortBy: sortField,
        sortDirection,
      });
      
      setOnsenLogs(response.data);
      setPagination({
        ...pagination,
        total: response.total || 0,
      });
    } catch (err) {
      setError('温泉メモの取得に失敗しました。もう一度お試しください。');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    // 検索処理
    fetchOnsenLogs();
  };

  const handleFilterChange = (key: keyof OnsenLogFilter, value: string | number) => {
    setFilters({
      ...filters,
      [key]: value,
    });
  };

  const applyFilters = () => {
    setPagination({
      ...pagination,
      page: 1, // フィルター適用時は1ページ目に戻る
    });
    fetchOnsenLogs();
    setShowFilters(false);
  };

  const resetFilters = () => {
    setFilters({
      name: '',
      location: '',
      minRating: 0,
      maxRating: 5,
      fromDate: '',
      toDate: '',
    });
  };

  const toggleSortDirection = () => {
    setSortDirection(sortDirection === 'asc' ? 'desc' : 'asc');
  };

  const handleExportJson = async () => {
    try {
      await exportOnsenLogsAsJson();
    } catch (err) {
      setError('エクスポートに失敗しました。');
      console.error(err);
    }
  };

  const handleExportCsv = async () => {
    try {
      await exportOnsenLogsAsCsv();
    } catch (err) {
      setError('エクスポートに失敗しました。');
      console.error(err);
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
    <Layout title="温泉メモ一覧 | 湯録 (Yuroku)">
      <div className="mb-6 flex flex-col md:flex-row md:items-center md:justify-between">
        <h1 className="text-2xl font-bold mb-4 md:mb-0">温泉メモ一覧</h1>
        <div className="flex flex-wrap gap-2">
          <Button
            onClick={() => router.push('/onsen/new')}
            icon={<FaPlus />}
          >
            新規メモ作成
          </Button>
          <Button
            variant="outline"
            onClick={() => setShowFilters(!showFilters)}
            icon={<FaFilter />}
          >
            フィルター
          </Button>
          <div className="relative">
            <Button
              variant="outline"
              onClick={toggleSortDirection}
              icon={sortDirection === 'asc' ? <FaSortAmountUp /> : <FaSortAmountDown />}
            >
              {sortField === 'visitDate' ? '訪問日' : '評価'}
            </Button>
          </div>
          <Button
            variant="outline"
            onClick={handleExportJson}
            icon={<FaDownload />}
          >
            JSON
          </Button>
          <Button
            variant="outline"
            onClick={handleExportCsv}
            icon={<FaDownload />}
          >
            CSV
          </Button>
        </div>
      </div>

      {/* 検索バー */}
      <form onSubmit={handleSearch} className="mb-6">
        <InputField
          placeholder="温泉名や場所で検索..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          icon={<FaSearch />}
          className="max-w-xl"
        />
      </form>

      {/* フィルターパネル */}
      {showFilters && (
        <Card className="mb-6">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            <InputField
              label="温泉名"
              value={filters.name}
              onChange={(e) => handleFilterChange('name', e.target.value)}
              placeholder="温泉名で絞り込み"
            />
            <InputField
              label="場所"
              value={filters.location}
              onChange={(e) => handleFilterChange('location', e.target.value)}
              placeholder="場所で絞り込み"
            />
            <InputField
              label="最低評価"
              type="number"
              min={0}
              max={5}
              value={filters.minRating}
              onChange={(e) => handleFilterChange('minRating', parseInt(e.target.value))}
            />
            <InputField
              label="最高評価"
              type="number"
              min={0}
              max={5}
              value={filters.maxRating}
              onChange={(e) => handleFilterChange('maxRating', parseInt(e.target.value))}
            />
            <InputField
              label="開始日"
              type="date"
              value={filters.fromDate}
              onChange={(e) => handleFilterChange('fromDate', e.target.value)}
            />
            <InputField
              label="終了日"
              type="date"
              value={filters.toDate}
              onChange={(e) => handleFilterChange('toDate', e.target.value)}
            />
          </div>
          <div className="flex justify-end mt-4 space-x-2">
            <Button variant="outline" onClick={resetFilters}>
              リセット
            </Button>
            <Button onClick={applyFilters}>
              適用
            </Button>
          </div>
        </Card>
      )}

      {/* エラーメッセージ */}
      {error && <ErrorMessage message={error} className="mb-6" />}

      {/* 読み込み中 */}
      {loading ? (
        <div className="flex justify-center items-center h-64">
          <Loading size="large" text="読み込み中..." />
        </div>
      ) : (
        <>
          {/* 温泉メモ一覧 */}
          {onsenLogs.length === 0 ? (
            <Card className="p-8 text-center">
              <p className="text-gray-500 mb-4">温泉メモがまだありません。</p>
              <Button
                onClick={() => router.push('/onsen/new')}
                icon={<FaPlus />}
              >
                最初のメモを作成
              </Button>
            </Card>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
              {onsenLogs.map((onsenLog) => (
                <OnsenCard key={onsenLog.id} onsenLog={onsenLog} />
              ))}
            </div>
          )}

          {/* ページネーション */}
          {pagination.total > pagination.limit && (
            <div className="flex justify-center mt-8">
              <div className="flex space-x-2">
                <Button
                  variant="outline"
                  disabled={pagination.page === 1}
                  onClick={() => setPagination({ ...pagination, page: pagination.page - 1 })}
                >
                  前へ
                </Button>
                <span className="flex items-center px-4 py-2 bg-gray-100 rounded">
                  {pagination.page} / {Math.ceil(pagination.total / pagination.limit)}
                </span>
                <Button
                  variant="outline"
                  disabled={pagination.page >= Math.ceil(pagination.total / pagination.limit)}
                  onClick={() => setPagination({ ...pagination, page: pagination.page + 1 })}
                >
                  次へ
                </Button>
              </div>
            </div>
          )}
        </>
      )}
    </Layout>
  );
};

export default OnsenLogsPage; 