import React, { useState } from 'react';
import { useRouter } from 'next/router';
import Link from 'next/link';
import { FaUser, FaLock, FaSignInAlt } from 'react-icons/fa';
import Layout from '@/components/layout/Layout';
import Card from '@/components/common/Card';
import InputField from '@/components/common/InputField';
import Button from '@/components/common/Button';
import ErrorMessage from '@/components/common/ErrorMessage';
import { login } from '@/services/auth';
import { ApiError } from '@/types';

const LoginPage: React.FC = () => {
  const router = useRouter();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);

    try {
      await login({ email, password });
      window.location.href = '/onsen';
    } catch (err) {
      const apiError = err as ApiError;
      setError(apiError.message || 'ログインに失敗しました。もう一度お試しください。');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Layout title="ログイン | 湯録 (Yuroku)">
      <div className="max-w-md mx-auto">
        <Card
          title="ログイン"
          subtitle="アカウント情報を入力してログインしてください"
          className="mt-8"
        >
          {error && <ErrorMessage message={error} className="mb-4" />}
          
          <form onSubmit={handleSubmit}>
            <div className="space-y-4">
              <InputField
                label="メールアドレス"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="example@example.com"
                required
                icon={<FaUser />}
              />
              
              <InputField
                label="パスワード"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="パスワードを入力"
                required
                icon={<FaLock />}
              />
              
              <Button
                type="submit"
                fullWidth
                isLoading={isLoading}
                loadingText="ログイン中..."
                icon={<FaSignInAlt />}
              >
                ログイン
              </Button>
            </div>
          </form>
        </Card>
        
        <div className="text-center mt-4">
          <p className="text-gray-600">
            アカウントをお持ちでない方は{' '}
            <Link href="/auth/register" className="text-primary-600 hover:text-primary-700">
              新規登録
            </Link>
          </p>
        </div>
      </div>
    </Layout>
  );
};

export default LoginPage; 