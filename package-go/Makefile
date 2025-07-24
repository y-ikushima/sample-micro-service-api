.PHONY: help install sqlc migrate-up migrate-down migrate-create build run test clean

# デフォルトのターゲット
help:
	@echo "利用可能なコマンド:"
	@echo "  install       - 依存関係をインストール"
	@echo "  sqlc          - SQLからGoコードを生成"
	@echo "  migrate-up    - データベースマイグレーションを実行"
	@echo "  migrate-down  - データベースマイグレーションをロールバック"
	@echo "  migrate-create - 新しいマイグレーションファイルを作成"
	@echo "  build         - アプリケーションをビルド"
	@echo "  run           - アプリケーションを実行"
	@echo "  test          - テストを実行"
	@echo "  clean         - ビルド成果物を削除"

# 依存関係のインストール
install:
	go mod tidy
	go mod download

# sqlcを使用してSQLからGoコード生成
sqlc:
	@which sqlc > /dev/null || (echo "sqlcがインストールされていません。go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest でインストールしてください" && exit 1)
	cd database && sqlc generate

# データベースマイグレーション実行
migrate-up:
	@if [ -z "$(DATABASE_URL)" ]; then echo "DATABASE_URL環境変数が設定されていません"; exit 1; fi
	cd database && migrate -path migrations -database "$(DATABASE_URL)" up

# データベースマイグレーションロールバック
migrate-down:
	@if [ -z "$(DATABASE_URL)" ]; then echo "DATABASE_URL環境変数が設定されていません"; exit 1; fi
	cd database && migrate -path migrations -database "$(DATABASE_URL)" down

# 新しいマイグレーションファイル作成
migrate-create:
	@if [ -z "$(name)" ]; then echo "使用法: make migrate-create name=migration_name"; exit 1; fi
	cd database && migrate create -ext sql -dir migrations -seq $(name)

# アプリケーションビルド
build:
	cd database && go build -o ../bin/db-tool cmd/main.go

# アプリケーション実行
run:
	cd database && go run cmd/main.go

# テスト実行
test:
	go test -v ./...

# ビルド成果物削除
clean:
	rm -rf bin/ 