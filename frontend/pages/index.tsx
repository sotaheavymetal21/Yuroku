import React from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { FaHotTub, FaBook, FaChartLine, FaCloudUploadAlt } from 'react-icons/fa';
import Layout from '@/components/layout/Layout';
import Button from '@/components/common/Button';
import { useAuth } from '@/contexts/AuthContext';

const HomePage: React.FC = () => {
  const { isLoggedIn } = useAuth();

  return (
    <Layout
      title="湯録 (Yuroku) - 温泉体験記録アプリ"
      description="湯録（Yuroku）は、あなたの温泉体験を記録・管理するためのアプリケーションです。訪れた温泉の情報や感想を簡単に記録し、思い出を大切に保存しましょう。"
    >
      {/* ヒーローセクション */}
      <section className="py-12 md:py-20">
        <div className="container mx-auto px-4">
          <div className="flex flex-col md:flex-row items-center">
            <div className="md:w-1/2 md:pr-12">
              <h1 className="text-4xl md:text-5xl font-bold text-gray-900 mb-4">
                あなたの温泉体験を<br />記録しよう
              </h1>
              <p className="text-xl text-gray-600 mb-8">
                湯録（Yuroku）は、訪れた温泉の情報や感想を簡単に記録し、
                思い出を大切に保存するためのアプリケーションです。
              </p>
              <div className="flex flex-wrap gap-4">
                {isLoggedIn ? (
                  <Link href="/onsen">
                    <Button size="large" icon={<FaBook />}>
                      温泉メモを見る
                    </Button>
                  </Link>
                ) : (
                  <>
                    <Link href="/auth/register">
                      <Button size="large" icon={<FaHotTub />}>
                        無料で始める
                      </Button>
                    </Link>
                    <Link href="/auth/login">
                      <Button size="large" variant="outline">
                        ログイン
                      </Button>
                    </Link>
                  </>
                )}
              </div>
            </div>
            <div className="md:w-1/2 mt-12 md:mt-0">
              <div className="relative h-80 md:h-96 w-full rounded-lg overflow-hidden shadow-xl">
                <Image
                  src="/images/onsen-hero.jpg"
                  alt="温泉のイメージ"
                  fill
                  style={{ objectFit: 'cover' }}
                  priority
                />
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* 特徴セクション */}
      <section className="py-12 bg-gray-50">
        <div className="container mx-auto px-4">
          <h2 className="text-3xl font-bold text-center mb-12">湯録の特徴</h2>
          <div className="grid md:grid-cols-3 gap-8">
            <div className="bg-white p-6 rounded-lg shadow-md">
              <div className="text-onsen text-4xl mb-4">
                <FaBook />
              </div>
              <h3 className="text-xl font-semibold mb-2">簡単記録</h3>
              <p className="text-gray-600">
                訪れた温泉の基本情報や感想を簡単に記録できます。写真も添付可能で、思い出をより鮮明に残せます。
              </p>
            </div>
            <div className="bg-white p-6 rounded-lg shadow-md">
              <div className="text-onsen text-4xl mb-4">
                <FaChartLine />
              </div>
              <h3 className="text-xl font-semibold mb-2">統計・分析</h3>
              <p className="text-gray-600">
                訪問した温泉の傾向や評価を分析。自分の好みの温泉タイプや訪問頻度などが一目でわかります。
              </p>
            </div>
            <div className="bg-white p-6 rounded-lg shadow-md">
              <div className="text-onsen text-4xl mb-4">
                <FaCloudUploadAlt />
              </div>
              <h3 className="text-xl font-semibold mb-2">クラウド保存</h3>
              <p className="text-gray-600">
                記録したデータはクラウドに安全に保存。どのデバイスからでもアクセスでき、データ消失の心配がありません。
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* CTAセクション */}
      <section className="py-16 bg-onsen-dark text-white">
        <div className="container mx-auto px-4 text-center">
          <h2 className="text-3xl font-bold mb-6">今すぐ湯録を始めましょう</h2>
          <p className="text-xl mb-8 max-w-2xl mx-auto">
            あなたの温泉体験を記録し、大切な思い出として残しませんか？
            無料でアカウントを作成して、湯録を始めましょう。
          </p>
          {!isLoggedIn && (
            <Link href="/auth/register">
              <Button size="large" variant="primary" className="bg-white text-onsen-dark hover:bg-gray-100">
                無料アカウント作成
              </Button>
            </Link>
          )}
        </div>
      </section>
    </Layout>
  );
};

export default HomePage; 