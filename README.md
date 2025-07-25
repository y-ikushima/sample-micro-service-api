# sample-micro-service-api

Sample Micro Service が提供する API サーバ

## 概要

このプロジェクトは、マイクロサービスアーキテクチャを採用したサンプル API サーバーです。
以下のサービスで構成されています：

- **app-service**: バックエンド API サーバー (Go + Gin)
- **app-web**: フロントエンド Web アプリケーション (Next.js + React)
- **postgres**: PostgreSQL データベース

## 前提条件

### 必須ソフトウェア

以下のソフトウェアがインストールされている必要があります：

- **Docker**: バージョン 24.0 以上
- **Docker Compose**: バージョン 2.0 以上
- **Node.js**: バージョン 20.x (LTS 推奨)
- **Go**: バージョン 1.24.x
- **pnpm**: パッケージマネージャー (npm install -g pnpm)

### インストール確認

```bash
# バージョン確認
docker --version
docker-compose --version
node --version
go version
pnpm --version
```

## 環境構築手順

### Step 1: リポジトリのクローン

```bash
# プロジェクトをクローン
git clone <repository-url>
cd sample-micro-service-api
```

### Step 2: 環境変数の設定

```bash
# 環境変数ファイルの作成
cp .env.example .env.local

# または手動で作成
touch .env.local
```

環境変数ファイル（`.env.local`）に以下の設定を追加：

```bash
# データベース設定
POSTGRES_DB=sample_db
POSTGRES_USER=sample_user
POSTGRES_PASSWORD=sample_password
POSTGRES_HOST=postgres
POSTGRES_PORT=5432

# アプリケーション設定
APP_ENV=development
APP_PORT=3003
WEB_PORT=3000

# 認証設定（必要に応じて）
JWT_SECRET=your-jwt-secret-key
AES_KEY=your-aes-encryption-key
```

### Step 3: 依存関係のインストール

```bash
# Node.js依存関係のインストール
pnpm install

# Go依存関係のダウンロード（必要に応じて）
cd apps/backend/app-service && go mod download
cd ../../../
```

### Step 4: Docker 環境の準備

```bash
# Dockerイメージのビルド
docker-compose build

# または高速ビルド
make build
```

### Step 5: データベースの初期化

```bash
# データベースとサービスの起動
make start

# データベースマイグレーション実行
make db-migrate

# 初期データの投入（オプション）
make db-seed
```

### Step 6: サービス起動の確認

```bash
# 全サービス起動
make start

# ログでサービス状態を確認
make logs
```

### Step 7: アクセス確認

各サービスが正常に起動していることを確認：

- **フロントエンド**: http://localhost:3000
- **バックエンド API**: http://localhost:3003
- **API 健全性チェック**: http://localhost:3003/health
- **PostgreSQL**: localhost:5432 (DB 接続ツールで確認可能)

### Step 8: 開発環境の確認

```bash
# API疎通テスト
curl http://localhost:3003/health

# システム一覧API テスト
curl http://localhost:3003/api/v1/systems
```

## 環境構築の完了確認

以下が全て成功すれば環境構築完了です：

- [ ] `make start` でエラーなく起動する
- [ ] http://localhost:3000 でフロントエンドにアクセスできる
- [ ] http://localhost:3003/health で `{"status": "ok"}` が返される
- [ ] `make logs` でエラーログが出ていない

## 開発時の注意点

### フロントエンドの依存関係を変更した場合

package.json を変更した際は、以下の手順で確実に反映してください：

```bash
# 方法1: コンテナを再ビルド（推奨）
make frontend-rebuild

# 方法2: 依存関係のみ再インストール
make frontend-clean

# 方法3: 手動で実行
docker-compose stop app-web
docker-compose build --no-cache app-web
docker-compose up app-web -d
```

### よくある問題と解決法

#### 1. 「Module not found」エラー

新しい npm パッケージが見つからない場合：

```bash
make frontend-clean
```

#### 2. フロントエンドが起動しない

```bash
make frontend-rebuild
```

#### 3. データベース接続エラー

```bash
make db-reset
```

## 開発用コマンド

```bash
# サービス管理
make start          # 全サービス起動
make stop           # 全サービス停止
make restart        # 全サービス再起動
make logs           # 全ログ表示

# フロントエンド
make frontend-rebuild  # フロントエンド再ビルド
make frontend-clean    # 依存関係をクリーンインストール

# バックエンド
make backend-restart   # バックエンド再起動

# データベース
make db-reset         # データベースリセット
make db-seed          # データベースシード実行

# 開発
make dev             # 開発モード起動（ログ表示）
make dev-logs        # 開発用ログ表示
```

## プロジェクト構成

```
sample-micro-service-api/
├── apps/
│   ├── backend/
│   │   └── app-service/          # Go APIサーバー
│   └── frontend/
│       └── app-web/              # Next.jsアプリ
├── package-go/                   # Go共通パッケージ
├── package-ts/                   # TypeScript共通パッケージ
├── doc/                          # API仕様書
└── docker-compose.yml
```

## API 仕様

### システム一覧取得

```
GET /api/v1/systems
```

クエリパラメータ：

- `systemName`: システム名で絞り込み
- `email`: メールアドレスで絞り込み
- `localGovernmentId`: 自治体 ID で絞り込み

## トラブルシューティング

### Docker キャッシュの問題

依存関係の変更が反映されない場合：

```bash
# すべてのキャッシュをクリア
docker system prune -a

# 特定のコンテナのみ再ビルド
docker-compose build --no-cache app-web
```

### ポート競合

他のアプリケーションがポートを使用している場合：

```bash
# ポート使用状況確認
lsof -i :3000
lsof -i :3003
lsof -i :5432

# 使用中のプロセスを停止
kill -9 <PID>
```

モデル生成

```
gentool -db postgres -dsn "postgres://postgres:password@localhost:5432/sample_micro_service?sslmode=disable" -onlyModel -outPath ./database/generated
```
