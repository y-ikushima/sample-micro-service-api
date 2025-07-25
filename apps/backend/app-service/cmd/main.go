package main

import (
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"sample-micro-service-api/apps/backend/app-service/internal/wire"
	"sample-micro-service-api/package-go/logging"
)

func main() {
	// ログの初期化（環境変数から自動設定）
	if err := logging.InitFromEnv(); err != nil {
		logging.Fatal("Failed to initialize logger", zap.Error(err))
	}
	defer logging.Sync()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logging.Warn("Warning: .env file not found", zap.Error(err))
	}

	// Initialize app using Wire (データベース、サービス、ハンドラーのみ)
	server, cleanup, err := wire.InitializeApp()
	if err != nil {
		logging.Fatal("Failed to initialize app", zap.Error(err))
	}
	defer cleanup()

	logging.Info("Database connection successful")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	logging.Info("Starting API server", 
		zap.String("port", port),
		zap.String("environment", env),
	)
	if err := server.Start(":" + port); err != nil {
		logging.Fatal("Failed to start server", zap.Error(err))
	}
} 