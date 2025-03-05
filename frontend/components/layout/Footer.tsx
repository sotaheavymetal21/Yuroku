import React from 'react';
import Link from 'next/link';
import { FaHotTub, FaGithub } from 'react-icons/fa';

const Footer: React.FC = () => {
  const currentYear = new Date().getFullYear();

  return (
    <footer className="bg-onsen-dark text-white py-6 mt-auto">
      <div className="container mx-auto px-4">
        <div className="flex flex-col md:flex-row justify-between items-center">
          <div className="mb-4 md:mb-0">
            <Link href="/" className="flex items-center space-x-2 text-xl font-bold">
              <FaHotTub className="text-onsen" />
              <span>湯録 (Yuroku)</span>
            </Link>
            <p className="text-sm mt-2 text-gray-300">
              温泉体験を記録・共有するためのアプリケーション
            </p>
          </div>

          <div className="flex flex-col items-center md:items-end">
            <div className="flex space-x-4 mb-2">
              <a
                href="https://github.com/yourusername/yuroku"
                target="_blank"
                rel="noopener noreferrer"
                className="text-gray-300 hover:text-white transition-colors"
              >
                <FaGithub size={20} />
              </a>
            </div>
            <p className="text-sm text-gray-300">
              &copy; {currentYear} 湯録 (Yuroku). All rights reserved.
            </p>
          </div>
        </div>
      </div>
    </footer>
  );
};

export default Footer; 