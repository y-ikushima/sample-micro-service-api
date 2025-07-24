# sample-micro-service-api

Sample Micro Service が提供する API サーバ

## 概要

このプロジェクトは、マイクロサービスアーキテクチャを採用したサンプル API サーバーです。
以下のサービスで構成されています：

- **app-service**: バックエンド API サーバー (Go + Gin)
- **app-web**: フロントエンド Web アプリケーション (Next.js + React)
- **postgres**: PostgreSQL データベース

## 前提条件

- Docker & Docker Compose
- Node.js 20.x
- Go 1.24.x

## セットアップ

### 1. 環境設定

```bash
# 環境変数ファイルの作成
cp ".env temp" .env.local
```

### 2. サービス起動

```bash
# 全サービス起動
make start

# または
docker-compose up -d
```

### 3. アクセス確認

- フロントエンド: http://localhost:3000
- バックエンド API: http://localhost:3003
- PostgreSQL: http://localhost:5432

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

## ライセンス

Apache 2.0

## サポート

問題が発生した場合は、まずこの README のトラブルシューティングセクションを確認してください。
