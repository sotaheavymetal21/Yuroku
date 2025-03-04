# Yuroku（湯録）バックエンド

Yuroku（湯録）アプリケーションのバックエンドコードです。Go言語とGinフレームワークを使用し、クリーンアーキテクチャに基づいて実装されています。

## アーキテクチャ

このバックエンドはクリーンアーキテクチャの原則に従って設計されています。詳細な設計ドキュメントは[こちら](../docs/backend-architecture.md)を参照してください。

## 主要なディレクトリ構造

```
backend/
├── cmd/
│   └── api/
│       └── main.go       # エントリーポイント
├── internal/
│   ├── domain/           # ドメイン層（エンティティとビジネスルール）
│   │   ├── entity/       # ドメインエンティティ
│   │   ├── repository/   # リポジトリインターフェース
│   │   └── service/      # ドメインサービス
│   ├── usecase/          # ユースケース層（アプリケーションロジック）
│   │   ├── interactor/   # ユースケース実装
│   │   └── port/         # 入出力ポート
│   ├── adapter/          # アダプター層（インターフェースアダプター）
│   │   ├── controller/   # コントローラー
│   │   ├── gateway/      # リポジトリ実装
│   │   └── presenter/    # プレゼンター
│   └── infrastructure/   # インフラストラクチャ層（フレームワークとドライバー）
│       ├── database/     # データベース接続
│       ├── middleware/   # ミドルウェア
│       ├── router/       # ルーティング
│       └── storage/      # ストレージ
└── pkg/                  # 外部パッケージ
```

## セットアップ

### 前提条件

- Go 1.21以上
- MongoDB
- Docker（オプション）

### 環境変数

`.env`ファイルを作成し、必要な環境変数を設定してください。`.env.example`を参考にしてください。

### 依存関係のインストール

```bash
go mod download
```

### 実行方法

```bash
# 直接実行
go run cmd/api/main.go

# ビルドして実行
go build -o yuroku-api cmd/api/main.go
./yuroku-api
```

### Dockerでの実行

```bash
docker build -t yuroku-api .
docker run -p 8080:8080 --env-file .env yuroku-api
```

## API仕様

APIの詳細な仕様については、Swaggerドキュメントを参照してください。サーバー起動後、`/swagger/index.html`にアクセスすることで確認できます。

## テスト

```bash
go test ./...
``` 