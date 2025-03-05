import React, { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/router';
import { FaHotTub, FaSignOutAlt, FaUser, FaBars, FaTimes } from 'react-icons/fa';
import { useAuth } from '@/contexts/AuthContext';

const Header: React.FC = () => {
  const { isLoggedIn, logout } = useAuth();
  const router = useRouter();
  const [isMenuOpen, setIsMenuOpen] = useState(false);

  const toggleMenu = () => {
    setIsMenuOpen(!isMenuOpen);
  };

  return (
    <header className="bg-onsen-dark text-white shadow-md">
      <div className="container mx-auto px-4 py-3">
        <div className="flex justify-between items-center">
          {/* ロゴ */}
          <Link href="/" className="flex items-center space-x-2 text-xl font-bold">
            <FaHotTub className="text-onsen" />
            <span>湯録 (Yuroku)</span>
          </Link>

          {/* モバイルメニューボタン */}
          <button
            className="md:hidden text-white focus:outline-none"
            onClick={toggleMenu}
          >
            {isMenuOpen ? <FaTimes size={24} /> : <FaBars size={24} />}
          </button>

          {/* デスクトップナビゲーション */}
          <nav className="hidden md:flex items-center space-x-6">
            {isLoggedIn ? (
              <>
                <Link
                  href="/onsen"
                  className={`hover:text-onsen transition-colors ${
                    router.pathname.startsWith('/onsen') ? 'text-onsen' : ''
                  }`}
                >
                  温泉メモ
                </Link>
                <button
                  onClick={logout}
                  className="flex items-center space-x-1 hover:text-onsen transition-colors"
                >
                  <FaSignOutAlt />
                  <span>ログアウト</span>
                </button>
              </>
            ) : (
              <>
                <Link
                  href="/auth/login"
                  className={`hover:text-onsen transition-colors ${
                    router.pathname === '/auth/login' ? 'text-onsen' : ''
                  }`}
                >
                  ログイン
                </Link>
                <Link
                  href="/auth/register"
                  className={`hover:text-onsen transition-colors ${
                    router.pathname === '/auth/register' ? 'text-onsen' : ''
                  }`}
                >
                  新規登録
                </Link>
              </>
            )}
          </nav>
        </div>

        {/* モバイルメニュー */}
        {isMenuOpen && (
          <nav className="md:hidden mt-4 pb-2">
            <ul className="space-y-3">
              {isLoggedIn ? (
                <>
                  <li>
                    <Link
                      href="/onsen"
                      className={`block hover:text-onsen transition-colors ${
                        router.pathname.startsWith('/onsen') ? 'text-onsen' : ''
                      }`}
                      onClick={() => setIsMenuOpen(false)}
                    >
                      温泉メモ
                    </Link>
                  </li>
                  <li>
                    <button
                      onClick={() => {
                        setIsMenuOpen(false);
                        logout();
                      }}
                      className="flex items-center space-x-1 hover:text-onsen transition-colors"
                    >
                      <FaSignOutAlt />
                      <span>ログアウト</span>
                    </button>
                  </li>
                </>
              ) : (
                <>
                  <li>
                    <Link
                      href="/auth/login"
                      className={`block hover:text-onsen transition-colors ${
                        router.pathname === '/auth/login' ? 'text-onsen' : ''
                      }`}
                      onClick={() => setIsMenuOpen(false)}
                    >
                      ログイン
                    </Link>
                  </li>
                  <li>
                    <Link
                      href="/auth/register"
                      className={`block hover:text-onsen transition-colors ${
                        router.pathname === '/auth/register' ? 'text-onsen' : ''
                      }`}
                      onClick={() => setIsMenuOpen(false)}
                    >
                      新規登録
                    </Link>
                  </li>
                </>
              )}
            </ul>
          </nav>
        )}
      </div>
    </header>
  );
};

export default Header; 