import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { useRouter } from 'next/router';
import Cookies from 'js-cookie';
import { isAuthenticated, logout } from '@/services/auth';

interface AuthContextType {
  isLoggedIn: boolean;
  loading: boolean;
  logout: () => void;
  refreshAuthState: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [isLoggedIn, setIsLoggedIn] = useState<boolean>(false);
  const [loading, setLoading] = useState<boolean>(true);
  const router = useRouter();

  // 認証状態の確認
  const checkAuth = () => {
    const authenticated = isAuthenticated();
    setIsLoggedIn(authenticated);
    setLoading(false);
  };

  // 初期化時と認証状態変更時に認証状態を確認
  useEffect(() => {
    checkAuth();

    // Cookieの変更を監視
    const tokenCheckInterval = setInterval(() => {
      const token = Cookies.get('token');
      const isCurrentlyLoggedIn = !!token;
      
      // ログイン状態が変わった場合に更新
      if (isCurrentlyLoggedIn !== isLoggedIn) {
        setIsLoggedIn(isCurrentlyLoggedIn);
      }
    }, 5000); // 5秒ごとにチェック

    return () => {
      clearInterval(tokenCheckInterval);
    };
  }, [isLoggedIn]);

  // 認証が必要なページへのアクセス制御
  useEffect(() => {
    if (!loading) {
      // 認証が必要なページのパス
      const authRequiredPaths = ['/onsen', '/profile'];
      
      // 現在のパスが認証必須かどうかをチェック
      const requiresAuth = authRequiredPaths.some(path => 
        router.pathname.startsWith(path)
      );

      // 認証が必要なのにログインしていない場合はログインページへリダイレクト
      if (requiresAuth && !isLoggedIn) {
        router.push('/auth/login');
        return;
      }

      // ログイン済みなのにログインページにアクセスした場合はホームへリダイレクト
      if (isLoggedIn && (router.pathname === '/auth/login' || router.pathname === '/auth/register')) {
        router.push('/onsen');
        return;
      }
    }
  }, [isLoggedIn, loading, router, router.pathname]);

  // 手動で認証状態を更新する関数
  const refreshAuthState = () => {
    checkAuth();
  };

  // ログアウト処理
  const handleLogout = () => {
    logout();
    setIsLoggedIn(false);
    router.push('/auth/login');
  };

  return (
    <AuthContext.Provider
      value={{
        isLoggedIn,
        loading,
        logout: handleLogout,
        refreshAuthState,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}; 