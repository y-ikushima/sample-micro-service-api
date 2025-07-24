package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"sample-micro-service-api/apps/backend/app-service/internal"
	"sample-micro-service-api/package-go/database"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Initialize database client
	dbClient, err := database.NewClient()
	if err != nil {
		log.Fatalf("Failed to initialize database client: %v", err)
	}
	defer dbClient.Close()

	// Test database connection
	if !dbClient.IsConnected() {
		log.Fatal("Database connection failed")
	}
	log.Println("✅ Database connection successful")

	// Initialize and start API server
	server := internal.NewServer(dbClient)
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Starting API server on port %s", port)
	if err := server.Start(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
} 