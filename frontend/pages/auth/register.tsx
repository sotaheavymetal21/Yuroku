import React, { useState } from 'react';
import { useRouter } from 'next/router';
import Link from 'next/link';
import { FaUser, FaLock, FaEnvelope, FaUserPlus } from 'react-icons/fa';
import Layout from '@/components/layout/Layout';
import Card from '@/components/common/Card';
import InputField from '@/components/common/InputField';
import Button from '@/components/common/Button';
import ErrorMessage from '@/components/common/ErrorMessage';
import { register } from '@/services/auth';
import { ApiError } from '@/types';

const RegisterPage: React.FC = () => {
  const router = useRouter();
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [validationErrors, setValidationErrors] = useState<{
    name?: string;
    email?: string;
    password?: string;
    confirmPassword?: string;
  }>({});

  const validateForm = () => {
    const errors: {
      name?: string;
      email?: string;
      password?: string;
      confirmPassword?: string;
    } = {};
    let isValid = true;

    if (!name.trim()) {
      errors.name = '名前を入力してください';
      isValid = false;
    }

    if (!email.trim()) {
      errors.email = 'メールアドレスを入力してください';
      isValid = false;
    } else if (!/\S+@\S+\.\S+/.test(email)) {
      errors.email = '有効なメールアドレスを入力してください';
      isValid = false;
    }

    if (!password) {
      errors.password = 'パスワードを入力してください';
      isValid = false;
    } else if (password.length < 8) {
      errors.password = 'パスワードは8文字以上である必要があります';
      isValid = false;
    }

    if (password !== confirmPassword) {
      errors.confirmPassword = 'パスワードが一致しません';
      isValid = false;
    }

    setValidationErrors(errors);
    return isValid;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setValidationErrors({});

    if (!validateForm()) {
      return;
    }

    setIsLoading(true);

    try {
      await register({ name, email, password });
      router.push('/onsen');
    } catch (err) {
      const apiError = err as ApiError;
      setError(apiError.message || '登録に失敗しました。もう一度お試しください。');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Layout title="新規登録 | 湯録 (Yuroku)">
      <div className="max-w-md mx-auto">
        <Card
          title="新規登録"
          subtitle="アカウントを作成して湯録を始めましょう"
          className="mt-8"
        >
          {error && <ErrorMessage message={error} className="mb-4" />}
          
          <form onSubmit={handleSubmit}>
            <div className="space-y-4">
              <InputField
                label="名前"
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="山田 太郎"
                required
                icon={<FaUser />}
                error={validationErrors.name}
              />
              
              <InputField
                label="メールアドレス"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="example@example.com"
                required
                icon={<FaEnvelope />}
                error={validationErrors.email}
              />
              
              <InputField
                label="パスワード"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="8文字以上のパスワード"
                required
                icon={<FaLock />}
                error={validationErrors.password}
              />
              
              <InputField
                label="パスワード（確認）"
                type="password"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                placeholder="パスワードを再入力"
                required
                icon={<FaLock />}
                error={validationErrors.confirmPassword}
              />
              
              <Button
                type="submit"
                fullWidth
                isLoading={isLoading}
                loadingText="登録中..."
                icon={<FaUserPlus />}
              >
                アカウント作成
              </Button>
            </div>
          </form>
        </Card>
        
        <div className="text-center mt-4">
          <p className="text-gray-600">
            すでにアカウントをお持ちの方は{' '}
            <Link href="/auth/login" className="text-primary-600 hover:text-primary-700">
              ログイン
            </Link>
          </p>
        </div>
      </div>
    </Layout>
  );
};

export default RegisterPage; 