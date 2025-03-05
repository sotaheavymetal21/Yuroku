# Yuroku（湯録）アプリケーション用Makefile

# デフォルト環境変数
ENV ?= dev

# ヘルプコマンド
.PHONY: help
help:
	@echo "使用可能なコマンド:"
	@echo "  make up              - 開発環境のコンテナを起動"
	@echo "  make down            - 開発環境のコンテナを停止"
	@echo "  make restart         - 開発環境のコンテナを再起動"
	@echo "  make logs            - コンテナのログを表示"
	@echo "  make ps              - 実行中のコンテナを表示"
	@echo "  make build           - コンテナをビルド"
	@echo "  make prod-up         - 本番環境のコンテナを起動"
	@echo "  make prod-down       - 本番環境のコンテナを停止"
	@echo "  make prod-restart    - 本番環境のコンテナを再起動"
	@echo "  make clean           - 未使用のDockerリソースを削除"
	@echo "  make frontend-shell  - フロントエンドコンテナのシェルに接続"
	@echo "  make backend-shell   - バックエンドコンテナのシェルに接続"
	@echo "  make mongo-shell     - MongoDBコンテナのシェルに接続"

# 開発環境コマンド
.PHONY: up
up:
	@echo "開発環境のコンテナを起動しています..."
	docker compose up -d

.PHONY: down
down:
	@echo "開発環境のコンテナを停止しています..."
	docker compose down

.PHONY: restart
restart: down up

.PHONY: logs
logs:
	docker compose logs -f

.PHONY: ps
ps:
	docker compose ps

.PHONY: build
build:
	@echo "コンテナをビルドしています..."
	docker compose build

# 本番環境コマンド
.PHONY: prod-up
prod-up:
	@echo "本番環境のコンテナを起動しています..."
	docker compose -f docker-compose.prod.yml up -d

.PHONY: prod-down
prod-down:
	@echo "本番環境のコンテナを停止しています..."
	docker compose -f docker-compose.prod.yml down

.PHONY: prod-restart
prod-restart: prod-down prod-up

# クリーンアップコマンド
.PHONY: clean
clean:
	@echo "未使用のDockerリソースを削除しています..."
	docker system prune -f

# シェルアクセスコマンド
.PHONY: frontend-shell
frontend-shell:
	docker compose exec frontend sh

.PHONY: backend-shell
backend-shell:
	docker compose exec backend sh

.PHONY: mongo-shell
mongo-shell:
	docker compose exec mongo mongo