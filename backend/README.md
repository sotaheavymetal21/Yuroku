# 湯録 (Yuroku) バックエンド

湯録アプリケーションのバックエンドAPI実装です。

## 技術スタック

- [Go](https://golang.org/)
- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [MongoDB](https://www.mongodb.com/)
- クリーンアーキテクチャによる設計

## セットアップ方法

### 必要環境

- Go 1.16以上
- MongoDB
- Docker & Docker Compose (オプション)

### ローカル開発環境の構築

1. リポジトリをクローン
```
git clone https://github.com/yourusername/yuroku.git
cd yuroku/backend
```

2. 依存パッケージのインストール
```
go mod download
```

3. 環境変数の設定
```
cp .env.example .env
```
`.env`ファイルを編集して必要な環境変数を設定します。

4. アプリケーションの実行
```
go run cmd/api/main.go
```

### Dockerでの実行

```
docker-compose up -d
```

## APIドキュメント

湯録アプリケーションのRESTful API仕様書です。

### 認証API

#### ユーザー登録

新しいユーザーを登録します。

- **URL**: `/api/auth/register`
- **Method**: `POST`
- **認証**: 不要

**リクエスト**:
```json
{
  "name": "テスト太郎",
  "email": "test@example.com",
  "password": "password123"
}
```

**レスポンス (成功)**:
```json
{
  "message": "ユーザー登録が完了しました"
}
```

**レスポンス (エラー)**:
```json
{
  "error": {
    "code": "REGISTRATION_FAILED",
    "message": "このメールアドレスは既に登録されています"
  }
}
```

#### ログイン

ユーザーアカウントにログインし、JWTトークンを取得します。

- **URL**: `/api/auth/login`
- **Method**: `POST`
- **認証**: 不要

**リクエスト**:
```json
{
  "email": "test@example.com",
  "password": "password123"
}
```

**レスポンス (成功)**:
```json
{
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  },
  "message": "ログインに成功しました"
}
```

**レスポンス (エラー)**:
```json
{
  "error": {
    "code": "LOGIN_FAILED",
    "message": "メールアドレスまたはパスワードが正しくありません"
  }
}
```

#### トークン更新

リフレッシュトークンを使用して、新しいアクセストークンを取得します。

- **URL**: `/api/auth/refresh`
- **Method**: `POST`
- **認証**: 不要

**リクエスト**:
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**レスポンス (成功)**:
```json
{
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  },
  "message": "トークンの更新に成功しました"
}
```

**レスポンス (エラー)**:
```json
{
  "error": {
    "code": "TOKEN_REFRESH_FAILED",
    "message": "無効なリフレッシュトークンです"
  }
}
```

#### ログアウト

現在のセッションからログアウトします。

- **URL**: `/api/auth/logout`
- **Method**: `POST`
- **認証**: 必要

**リクエスト**:
```
Authorization: Bearer {access_token}
```

**レスポンス (成功)**:
```json
{
  "message": "ログアウトしました"
}
```

### 温泉メモAPI

#### 温泉メモ一覧の取得

ユーザーの温泉メモ一覧を取得します。

- **URL**: `/api/onsen-logs`
- **Method**: `GET`
- **認証**: 必要

**クエリパラメータ**:
- `page`: ページ番号（デフォルト: 1）
- `limit`: 1ページあたりの件数（デフォルト: 10）
- `spring_type`: 泉質でフィルタリング（オプション）
- `rating`: 最小評価でフィルタリング（オプション）
- `start_date`: 訪問日の開始日でフィルタリング（オプション、ISO8601形式）
- `end_date`: 訪問日の終了日でフィルタリング（オプション、ISO8601形式）
- `keyword`: 名前や所在地でフィルタリング（オプション）

**レスポンス (成功)**:
```json
{
  "data": {
    "onsen_logs": [
      {
        "id": "60a1b2c3d4e5f6a7b8c9d0e1",
        "user_id": "60a1b2c3d4e5f6a7b8c9d0e2",
        "name": "草津温泉",
        "location": "群馬県吾妻郡草津町",
        "spring_type": "酸性泉",
        "features": ["露天風呂あり", "景色が良い"],
        "visit_date": "2023-01-15T00:00:00Z",
        "rating": 5,
        "comment": "とても良い温泉でした。また行きたいです。",
        "created_at": "2023-01-16T15:30:45Z",
        "updated_at": "2023-01-16T15:30:45Z"
      }
    ],
    "total": 42,
    "page": 1,
    "limit": 10
  },
  "message": "温泉メモ一覧を取得しました"
}
```

#### 温泉メモの作成

新しい温泉メモを作成します。

- **URL**: `/api/onsen-logs`
- **Method**: `POST`
- **認証**: 必要

**リクエスト**:
```json
{
  "name": "草津温泉",
  "location": "群馬県吾妻郡草津町",
  "spring_type": "酸性泉",
  "features": ["露天風呂あり", "景色が良い"],
  "visit_date": "2023-01-15T00:00:00Z",
  "rating": 5,
  "comment": "とても良い温泉でした。また行きたいです。"
}
```

**レスポンス (成功)**:
```json
{
  "data": {
    "id": "60a1b2c3d4e5f6a7b8c9d0e1",
    "name": "草津温泉",
    "location": "群馬県吾妻郡草津町",
    "spring_type": "酸性泉",
    "features": ["露天風呂あり", "景色が良い"],
    "visit_date": "2023-01-15T00:00:00Z",
    "rating": 5,
    "comment": "とても良い温泉でした。また行きたいです。",
    "created_at": "2023-01-16T15:30:45Z",
    "updated_at": "2023-01-16T15:30:45Z"
  },
  "message": "温泉メモを作成しました"
}
```

#### 温泉メモの詳細取得

特定の温泉メモの詳細を取得します。

- **URL**: `/api/onsen-logs/{id}`
- **Method**: `GET`
- **認証**: 必要

**レスポンス (成功)**:
```json
{
  "data": {
    "id": "60a1b2c3d4e5f6a7b8c9d0e1",
    "user_id": "60a1b2c3d4e5f6a7b8c9d0e2",
    "name": "草津温泉",
    "location": "群馬県吾妻郡草津町",
    "spring_type": "酸性泉",
    "features": ["露天風呂あり", "景色が良い"],
    "visit_date": "2023-01-15T00:00:00Z",
    "rating": 5,
    "comment": "とても良い温泉でした。また行きたいです。",
    "created_at": "2023-01-16T15:30:45Z",
    "updated_at": "2023-01-16T15:30:45Z",
    "images": [
      {
        "image_id": "60a1b2c3d4e5f6a7b8c9d0e3",
        "image_url": "https://yuroku.example.com/uploads/images/60a1b2c3d4e5f6a7b8c9d0e3.jpg"
      }
    ]
  },
  "message": "温泉メモを取得しました"
}
```

#### 温泉メモの更新

温泉メモの情報を更新します。

- **URL**: `/api/onsen-logs/{id}`
- **Method**: `PUT`
- **認証**: 必要

**リクエスト**:
```json
{
  "name": "草津温泉（改訂版）",
  "rating": 4,
  "comment": "2回目の訪問でした。やはり良い温泉です。"
}
```

**レスポンス (成功)**:
```json
{
  "data": {
    "id": "60a1b2c3d4e5f6a7b8c9d0e1",
    "name": "草津温泉（改訂版）",
    "location": "群馬県吾妻郡草津町",
    "spring_type": "酸性泉",
    "features": ["露天風呂あり", "景色が良い"],
    "visit_date": "2023-01-15T00:00:00Z",
    "rating": 4,
    "comment": "2回目の訪問でした。やはり良い温泉です。",
    "created_at": "2023-01-16T15:30:45Z",
    "updated_at": "2023-01-17T10:22:33Z"
  },
  "message": "温泉メモを更新しました"
}
```

#### 温泉メモの削除

温泉メモを削除します。

- **URL**: `/api/onsen-logs/{id}`
- **Method**: `DELETE`
- **認証**: 必要

**レスポンス (成功)**:
```json
{
  "message": "温泉メモを削除しました"
}
```

### 温泉画像API

#### 画像のアップロード

温泉メモに画像をアップロードします。

- **URL**: `/api/onsen-logs/{id}/images`
- **Method**: `POST`
- **認証**: 必要
- **Content-Type**: `multipart/form-data`

**リクエスト**:
```
image: (ファイル)
```

**レスポンス (成功)**:
```json
{
  "data": {
    "image_id": "60a1b2c3d4e5f6a7b8c9d0e3",
    "image_url": "https://yuroku.example.com/uploads/images/60a1b2c3d4e5f6a7b8c9d0e3.jpg"
  },
  "message": "画像をアップロードしました"
}
```

#### 画像の削除

温泉メモから画像を削除します。

- **URL**: `/api/onsen-logs/{onsen_id}/images/{image_id}`
- **Method**: `DELETE`
- **認証**: 必要

**レスポンス (成功)**:
```json
{
  "message": "画像を削除しました"
}
```

## エラーコード一覧

APIレスポンスのエラーコードと意味の対応表です。

| エラーコード | 説明 |
|-------------|------|
| REGISTRATION_FAILED | ユーザー登録に失敗しました |
| LOGIN_FAILED | ログインに失敗しました |
| TOKEN_REFRESH_FAILED | トークン更新に失敗しました |
| MISSING_TOKEN | 認証トークンがありません |
| INVALID_TOKEN_FORMAT | トークン形式が無効です |
| INVALID_TOKEN | 無効なトークンです |
| TOKEN_EXPIRED | トークンの有効期限が切れています |
| AUTHENTICATION_REQUIRED | 認証が必要です |
| NOT_FOUND | リソースが見つかりません |
| INVALID_INPUT | 入力データが無効です |
| FORBIDDEN | このリソースにアクセスする権限がありません |
| SERVER_ERROR | サーバー内部エラーが発生しました |

## 開発者向け情報

### コードの構成

プロジェクトはクリーンアーキテクチャに基づいて設計されています：

- **エンティティ層** (`internal/domain/entity`): ビジネスオブジェクトと基本的なルール
- **ユースケース層** (`internal/usecase`): アプリケーションに特化したビジネスルール
- **インターフェースアダプター層** (`internal/adapter`): 外部インターフェースとの連携
- **フレームワーク・ドライバー層** (`internal/infrastructure`): 外部フレームワークとの連携

### テスト

テストの実行方法：

```
go test ./...
```

カバレッジレポートの作成：

```
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
``` 