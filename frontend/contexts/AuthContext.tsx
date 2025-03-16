import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { useRouter } from 'next/router';
import { isAuthenticated, logout } from '@/services/auth';

interface AuthContextType {
  isLoggedIn: boolean;
  loading: boolean;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [isLoggedIn, setIsLoggedIn] = useState<boolean>(false);
  const [loading, setLoading] = useState<boolean>(true);
  const router = useRouter();

  useEffect(() => {
    // 認証状態の確認
    const checkAuth = () => {
      const authenticated = isAuthenticated();
      setIsLoggedIn(authenticated);
      setLoading(false);
    };

    checkAuth();
  }, []);

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
        window.location.href = '/auth/login';
        return;
      }

      // ログイン済みなのにログインページにアクセスした場合はホームへリダイレクト
      if (isLoggedIn && (router.pathname === '/auth/login' || router.pathname === '/auth/register')) {
        window.location.href = '/onsen';
        return;
      }
    }
  }, [isLoggedIn, loading, router.pathname]);

  const handleLogout = () => {
    logout();
    setIsLoggedIn(false);
    window.location.href = '/auth/login';
  };

  return (
    <AuthContext.Provider
      value={{
        isLoggedIn,
        loading,
        logout: handleLogout,
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