package database

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Client struct {
	GormDB *gorm.DB
}

// NewClient creates a new GORM database client
func NewClient() (*Client, error) {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		// .envファイルが見つからない場合は警告のみ
	}

	databaseURL := os.Getenv("POSTGRES_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("POSTGRES_URL environment variable is required")
	}

	// GORM database connection
	gormDB, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database with GORM: %w", err)
	}

	return &Client{
		GormDB: gormDB,
	}, nil
}

// NewGormClient creates a new GORM-only database client
func NewGormClient() (*gorm.DB, error) {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		// .envファイルが見つからない場合は警告のみ
	}

	databaseURL := os.Getenv("POSTGRES_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("POSTGRES_URL environment variable is required")
	}

	// GORM database connection
	gormDB, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database with GORM: %w", err)
	}

	return gormDB, nil
}

// Close closes the database connection
func (c *Client) Close() error {
	sqlDB, err := c.GormDB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// IsConnected checks if the database connection is alive
func (c *Client) IsConnected() bool {
	sqlDB, err := c.GormDB.DB()
	if err != nil {
		return false
	}
	return sqlDB.Ping() == nil
} 