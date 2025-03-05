# Yuroku（湯録）フロントエンド

Yuroku（湯録）アプリケーションのフロントエンドコードです。React と Next.js を使用して実装されています。

## 技術スタック

- React 18
- Next.js 14
- TypeScript
- Tailwind CSS
- React Hook Form
- Axios
- React Query
- React Icons
- React Datepicker
- React Toastify

## 機能

- ユーザー認証（登録・ログイン・ログアウト）
- 温泉メモの作成・閲覧・編集・削除
- 温泉画像のアップロード・表示・削除
- 温泉メモの検索・フィルタリング
- 温泉メモのエクスポート（JSON/CSV）
- レスポンシブデザイン

## セットアップ

### 前提条件

- Node.js 18以上
- npm または yarn

### 依存関係のインストール

```bash
# npmの場合
npm install

# yarnの場合
yarn install
```

### 開発サーバーの起動

```bash
# npmの場合
npm run dev

# yarnの場合
yarn dev
```

開発サーバーは http://localhost:3000 で起動します。

### ビルド

```bash
# npmの場合
npm run build

# yarnの場合
yarn build
```

### 本番モードでの起動

```bash
# npmの場合
npm run start

# yarnの場合
yarn start
```

## 環境変数

`.env.local`ファイルを作成し、以下の環境変数を設定してください：

```
NEXT_PUBLIC_API_URL=http://localhost:8080/v1
```

## ディレクトリ構造

```
frontend/
├── components/           # 再利用可能なコンポーネント
│   ├── common/           # 共通コンポーネント
│   ├── layout/           # レイアウト関連
│   └── onsen/            # 温泉関連コンポーネント
├── hooks/                # カスタムフック
├── pages/                # ページコンポーネント
│   ├── api/              # APIルート
│   ├── auth/             # 認証関連ページ
│   └── onsen/            # 温泉関連ページ
├── public/               # 静的ファイル
├── styles/               # グローバルスタイル
├── utils/                # ユーティリティ関数
├── contexts/             # Contextプロバイダー
├── services/             # APIサービス
└── types/                # 型定義
``` 