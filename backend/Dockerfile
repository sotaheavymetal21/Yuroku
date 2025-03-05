FROM golang:1.21-alpine AS builder

WORKDIR /app

# 依存関係をダウンロードするためのファイルをコピー
COPY go.mod ./
# go.sumファイルが存在しない場合でもエラーにならないようにする
RUN touch go.sum
RUN go mod download && go mod tidy

# ソースコードをコピー
COPY . .

# アプリケーションをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -o yuroku-api ./cmd/api

# 実行ステージ
FROM alpine:latest

WORKDIR /app

# ビルドしたバイナリをコピー
COPY --from=builder /app/yuroku-api .

# 必要なディレクトリを作成
RUN mkdir -p /app/storage

# APIサーバーのポートを公開
EXPOSE 8080

# アプリケーションを実行
CMD ["./yuroku-api"] 