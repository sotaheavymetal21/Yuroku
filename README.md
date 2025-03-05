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

## Docker環境でのセットアップ（推奨）

### 前提条件

- Docker
- Docker Compose

### 手順

1. リポジトリをクローン
   ```
   git clone https://github.com/yourusername/yuroku.git
   cd yuroku
   ```

2. 環境変数ファイルの準備
   ```
   # バックエンド
   cp backend/.env.example backend/.env
   
   # フロントエンド
   cp frontend/.env.example frontend/.env.local
   ```

3. Dockerコンテナの起動
   ```
   docker-compose up
   ```

4. アプリケーションへのアクセス
   - フロントエンド: http://localhost:3000
   - バックエンドAPI: http://localhost:8080
   - MongoDB管理画面: http://localhost:8081

### 開発時の便利なコマンド

- コンテナをバックグラウンドで起動
  ```
  docker-compose up -d
  ```

- コンテナのログを確認
  ```
  docker-compose logs -f
  ```

- 特定のサービスのログを確認
  ```
  docker-compose logs -f frontend
  docker-compose logs -f backend
  ```

- コンテナの停止
  ```
  docker-compose down
  ```

- コンテナとボリュームの削除（データベースも削除）
  ```
  docker-compose down -v
  ```

## 従来の方法でのセットアップ

### 前提条件

- Go 1.21以上
- Node.js 18以上
- MongoDB 6.0以上

### バックエンドのセットアップ

1. バックエンドディレクトリに移動
   ```
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

### Dockerを使用したデプロイ

1. 本番用の環境変数を設定
   ```
   # バックエンド
   cp backend/.env.example backend/.env.prod
   # 本番用の設定を編集
   
   # フロントエンド
   cp frontend/.env.example frontend/.env.prod
   # 本番用の設定を編集
   ```

2. 本番用のDockerイメージをビルド
   ```
   docker-compose -f docker-compose.prod.yml build
   ```

3. 本番環境でコンテナを起動
   ```
   docker-compose -f docker-compose.prod.yml up -d
   ```

### 従来の方法でのデプロイ

#### バックエンド

1. バイナリのビルド
   ```
   cd backend
   go build -o yuroku-api cmd/api/main.go
   ```

2. サーバーへの転送とサービス設定
   ```
   # サーバー固有の手順に従ってください
   ```

#### フロントエンド

1. 本番用ビルド
   ```
   cd frontend
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

## Makefileを使った簡単な起動方法

プロジェクトのルートディレクトリに`Makefile`を用意しています。以下のコマンドで簡単にアプリケーションを起動・管理できます。

### 基本的なコマンド

```bash
# 開発環境のコンテナを起動
make up

# 開発環境のコンテナを停止
make down

# 開発環境のコンテナを再起動
make restart

# コンテナのログを表示
make logs

# 実行中のコンテナを表示
make ps

# コンテナをビルド
make build
```

### 本番環境用コマンド

```bash
# 本番環境のコンテナを起動
make prod-up

# 本番環境のコンテナを停止
make prod-down

# 本番環境のコンテナを再起動
make prod-restart
```

### その他のユーティリティコマンド

```bash
# 未使用のDockerリソースを削除
make clean

# フロントエンドコンテナのシェルに接続
make frontend-shell

# バックエンドコンテナのシェルに接続
make backend-shell

# MongoDBコンテナのシェルに接続
make mongo-shell

# 使用可能なコマンド一覧を表示
make help
```
