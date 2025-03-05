# 湯録 (Yuroku) - 温泉体験記録アプリ

湯録（Yuroku）は、あなたの温泉体験を記録・管理するためのアプリケーションです。訪れた温泉の情報や感想を簡単に記録し、思い出を大切に保存しましょう。

## 機能

- 温泉訪問記録の作成・編集・削除
- 温泉の評価・コメント機能
- 温泉写真のアップロード
- 訪問記録の検索・フィルタリング
- データのエクスポート（JSON/CSV）

## 技術スタック

### バックエンド

- Go
- Gin Webフレームワーク
- MongoDB
- クリーンアーキテクチャ

### フロントエンド

- Next.js
- React
- TypeScript
- Tailwind CSS

## プロジェクト構成

```
yuroku/
├── backend/             # Goバックエンド
│   ├── cmd/             # エントリーポイント
│   ├── internal/        # 内部パッケージ
│   │   ├── domain/      # ドメイン層
│   │   ├── usecase/     # ユースケース層
│   │   ├── adapter/     # アダプター層
│   │   └── infrastructure/ # インフラ層
│   └── pkg/             # 外部パッケージ
└── frontend/            # Next.jsフロントエンド
    ├── components/      # Reactコンポーネント
    ├── contexts/        # Contextプロバイダー
    ├── hooks/           # カスタムフック
    ├── pages/           # ページコンポーネント
    ├── public/          # 静的ファイル
    ├── services/        # APIサービス
    ├── styles/          # スタイル
    ├── types/           # 型定義
    └── utils/           # ユーティリティ
```

## セットアップ手順

### 前提条件

- Go 1.21以上
- Node.js 18以上
- MongoDB 6.0以上
- Docker (オプション)

### バックエンドのセットアップ

1. リポジトリをクローン
   ```
   git clone https://github.com/yourusername/yuroku.git
   cd yuroku/backend
   ```

2. 依存関係のインストール
   ```
   go mod download
   ```

3. 環境変数の設定
   ```
   cp .env.example .env
   # .envファイルを編集して必要な設定を行う
   ```

4. サーバーの起動
   ```
   go run cmd/api/main.go
   ```

### フロントエンドのセットアップ

1. フロントエンドディレクトリに移動
   ```
   cd ../frontend
   ```

2. 依存関係のインストール
   ```
   npm install
   ```

3. 環境変数の設定
   ```
   cp .env.example .env.local
   # .env.localファイルを編集して必要な設定を行う
   ```

4. 開発サーバーの起動
   ```
   npm run dev
   ```

## 本番環境へのデプロイ

### バックエンド

1. バイナリのビルド
   ```
   go build -o yuroku-api cmd/api/main.go
   ```

2. サーバーへの転送とサービス設定
   ```
   # サーバー固有の手順に従ってください
   ```

### フロントエンド

1. 本番用ビルド
   ```
   npm run build
   ```

2. 静的ファイルのデプロイ
   ```
   npm run start
   # または、Vercel、Netlifyなどのサービスを利用
   ```

## ライセンス

MIT

## 貢献

プロジェクトへの貢献は大歓迎です。Issue報告や機能提案、プルリクエストなどお気軽にどうぞ。
