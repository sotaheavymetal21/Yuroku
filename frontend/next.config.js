/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  images: {
    domains: ['localhost', 'storage.googleapis.com'],
  },
  env: {
    API_URL: process.env.API_URL || 'http://localhost:8080/api',
  },
  // Docker環境用の設定
  output: 'standalone',
};

module.exports = nextConfig;