package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"sample-micro-service-api/apps/backend/app-service/internal/wire"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Initialize app using Wire
	server, cleanup, err := wire.InitializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}
	defer cleanup()

	log.Println("✅ Database connection successful")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Starting API server on port %s", port)
	if err := server.Start(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
} 