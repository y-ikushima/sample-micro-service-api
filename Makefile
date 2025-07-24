.PHONY: up down logs shell migrate-up migrate-down migrate-reset seed-db wire-gen

# Docker Compose コマンド
up:
	docker compose up --build

watch:
	docker compose watch

down:
	docker compose down

logs:
	docker compose logs -f

# Wire関連コマンド
wire-gen:
	docker compose exec app-service sh -c "cd /apps/backend/app-service && wire gen ./internal/wire"

# マイグレーション関連コマンド（コンテナ内で実行）
migrate-up:
	docker compose exec app-service sh -c "cd /package-go/database && go run cmd/main.go -migrate-up"

migrate-down:
	docker compose exec app-service sh -c "cd /package-go/database && go run cmd/main.go -migrate-down"

migrate-reset:
	docker compose exec app-service sh -c "cd /package-go/database && go run cmd/main.go -migrate-reset"

seed-db:
	docker compose exec app-service sh -c "cd /package-go/database && go run cmd/main.go -seed-db"

test-db:
	docker compose exec app-service sh -c "cd /package-go/database && go run cmd/main.go -test-db"

# コンテナ内にシェルでアクセス
shell:
	docker compose exec app-service sh

# PostgreSQLコンテナにアクセス
psql:
	docker compose exec postgres psql -U postgres -d sample_micro_service

# ヘルプ
help:
	@echo "利用可能なコマンド:"
	@echo "  make watch       - docker compose watchでサービス起動"
	@echo "  make up          - docker compose upでサービス起動"
	@echo "  make down        - サービス停止"
	@echo "  make logs        - ログ表示"
	@echo "  make wire-gen    - Wire依存性注入コード生成"
	@echo "  make migrate-up  - マイグレーション実行"
	@echo "  make migrate-down - マイグレーション巻き戻し"
	@echo "  make migrate-reset - マイグレーションリセット"
	@echo "  make seed-db     - テストデータ投入"
	@echo "  make test-db     - DB接続テスト"
	@echo "  make shell       - app-serviceコンテナ内シェル"
	@echo "  make psql        - PostgreSQLコンテナ接続" 