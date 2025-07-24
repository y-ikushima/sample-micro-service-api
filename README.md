# Sample Micro Service API

Go 言語で構築されたマイクロサービスアーキテクチャのサンプル実装です。

## 構成

```
sample-micro-service-api/
├── package-go/              # 共通ライブラリ
│   ├── database/            # データベース操作ライブラリ
│   │   ├── cmd/            # CLI ツール
│   │   ├── internal/       # 内部パッケージ（sqlc生成物）
│   │   ├── migrations/     # データベースマイグレーション
│   │   ├── queries/        # SQLクエリファイル
│   │   ├── client.go       # データベースクライアント
│   │   ├── migrator.go     # マイグレーション管理
│   │   ├── types.go        # 型定義の再エクスポート
│   │   └── sqlc.yaml       # sqlc設定
│   ├── go.mod
│   ├── go.sum
│   └── Makefile
├── apps/
│   └── backend/            # API サーバー
│       ├── cmd/            # メインアプリケーション
│       └── internal/       # 内部API実装
└── doc/                    # ドキュメント
```

## 機能

### package-go/database

- **データベースマイグレーション**: golang-migrate を使用
- **SQL クエリ生成**: sqlc を使用してタイプセーフなクエリを生成
- **データベースクライアント**: PostgreSQL 接続とクエリ実行
- **CLI ツール**: マイグレーション、シーディング、テスト用コマンド

### apps/backend

- **REST API**: Gin フレームワークを使用
- **CRUD 操作**: Users, Products, Orders
- **マイクロサービス連携**: package-go/database ライブラリを使用
- **管理機能**: データベース管理用 API

## セットアップ

### 前提条件

- Go 1.21+
- PostgreSQL
- sqlc
- migrate

### 1. 依存関係のインストール

```bash
# SQLCとmigrateツールのインストール
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# package-goの依存関係
cd package-go
go mod tidy

# apps/backendの依存関係
cd ../apps/backend
go mod tidy
```

### 2. 環境変数の設定

```bash
# package-go/.env
DATABASE_URL=postgres://user:password@localhost:5432/dbname?sslmode=disable
MIGRATION_DIR=./migrations

# apps/backend/.env
DATABASE_URL=postgres://user:password@localhost:5432/dbname?sslmode=disable
PORT=8080
GIN_MODE=debug
CORS_ALLOW_ORIGINS=http://localhost:3000,http://localhost:8080
CORS_ALLOW_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOW_HEADERS=Content-Type,Authorization
LOG_LEVEL=debug
```

### 3. データベースのセットアップ

```bash
# マイグレーション実行
cd package-go
make migrate-up

# サンプルデータ投入
cd database && go run cmd/main.go -seed-db
```

### 4. API サーバーの起動

```bash
cd apps/backend
go run cmd/main.go
```

## API エンドポイント

### ヘルスチェック

- `GET /health` - サーバーの状態確認

### ユーザー管理

- `GET /api/v1/users` - ユーザー一覧取得
- `POST /api/v1/users` - ユーザー作成
- `GET /api/v1/users/:id` - ユーザー詳細取得
- `PUT /api/v1/users/:id` - ユーザー更新
- `DELETE /api/v1/users/:id` - ユーザー削除

### 商品管理

- `GET /api/v1/products` - 商品一覧取得
- `POST /api/v1/products` - 商品作成
- `GET /api/v1/products/:id` - 商品詳細取得
- `PUT /api/v1/products/:id` - 商品更新
- `DELETE /api/v1/products/:id` - 商品削除

### 注文管理

- `GET /api/v1/orders` - 注文一覧取得
- `POST /api/v1/orders` - 注文作成
- `GET /api/v1/orders/:id` - 注文詳細取得
- `PUT /api/v1/orders/:id` - 注文ステータス更新
- `DELETE /api/v1/orders/:id` - 注文削除

### 管理機能

- `POST /api/v1/admin/migrate` - マイグレーション実行
- `POST /api/v1/admin/reset` - データベースリセット
- `POST /api/v1/admin/seed` - サンプルデータ投入

## 使用技術

### Backend (Go)

- **Web Framework**: Gin
- **Database**: PostgreSQL
- **ORM/Query Builder**: sqlc
- **Migration**: golang-migrate
- **Config**: godotenv

### Development Tools

- **Code Generation**: sqlc
- **Database Migration**: migrate
- **Dependency Management**: Go modules

## 開発コマンド

### package-go

```bash
cd package-go

# ヘルプ表示
make help

# 依存関係インストール
make install

# SQLからGoコード生成
make sqlc

# マイグレーション実行
make migrate-up

# マイグレーションロールバック
make migrate-down

# 新しいマイグレーション作成
make migrate-create name=add_new_table

# データベースツールビルド
make build

# データベースツール実行
make run

# テスト実行
make test
```

### apps/backend

```bash
cd apps/backend

# 開発サーバー起動
go run cmd/main.go

# ビルド
go build -o bin/server cmd/main.go

# テスト実行
go test -v ./...
```

## アーキテクチャの特徴

### マイクロサービス設計

- **package-go/database**: データベース操作専用ライブラリ
- **apps/backend**: API サーバー（package-go/database を依存関係として使用）
- **疎結合**: 各サービスが独立して開発・デプロイ可能

### モジュラー構造

- **機能別ディレクトリ**: `database/`、将来的に`auth/`、`logging/`なども追加可能
- **明確な責任分離**: データベース操作と API サーバーが分離
- **再利用性**: package-go ライブラリは複数のマイクロサービスで共有可能

### データベース管理

- **マイグレーション**: バージョン管理されたスキーマ変更
- **タイプセーフ**: sqlc による自動生成されたタイプセーフなクエリ
- **CLI ツール**: 開発・運用を支援するコマンドラインツール

### 拡張性

- 新しいマイクロサービスを追加する際は package-go ライブラリを再利用可能
- package-go 内に新しい機能パッケージを追加可能（例：`auth/`、`config/`、`logging/`）
- API バージョニング対応
- CORS 設定によるフロントエンド連携対応

## ライセンス

MIT License
