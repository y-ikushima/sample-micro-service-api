package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"sample-micro-service-api/package-go/database/internal/db"
)

type Client struct {
	DB      *sql.DB
	Queries *db.Queries
}

// NewClient creates a new database client
func NewClient() (*Client, error) {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		// .envファイルが見つからない場合は警告のみ
	}

	databaseURL := os.Getenv("POSTGRES_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("POSTGRES_URL environment variable is required")
	}

	database, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test connection
	if err := database.Ping(); err != nil {
		database.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Client{
		DB:      database,
		Queries: db.New(database),
	}, nil
}

// Close closes the database connection
func (c *Client) Close() error {
	return c.DB.Close()
}

// IsConnected checks if the database connection is alive
func (c *Client) IsConnected() bool {
	return c.DB.Ping() == nil
} 