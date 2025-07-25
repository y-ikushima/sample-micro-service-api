package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"sample-micro-service-api/package-go/database/seed"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables from parent directory
	if err := godotenv.Load("../.env"); err != nil {
		if err := godotenv.Load("../../.env"); err != nil {
			log.Printf("Warning: .env file not found: %v", err)
		}
	}

	// Define CLI commands
	var (
		migrateUp    = flag.Bool("migrate-up", false, "Run database migrations up")
		migrateDown  = flag.Bool("migrate-down", false, "Run database migrations down")
		migrateReset = flag.Bool("migrate-reset", false, "Reset database (down then up)")
		testDB       = flag.Bool("test-db", false, "Test database connection")
		seedDB       = flag.Bool("seed-db", false, "Seed database with sample data")
	)
	flag.Parse()

	databaseURL := os.Getenv("POSTGRES_URL")
	if databaseURL == "" {
		log.Fatal("POSTGRES_URL environment variable is required")
	}

	// Connect to database
	database, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Test database connection
	if err := database.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	switch {
	case *migrateUp:
		if err := runMigrationsUp(database); err != nil {
			log.Fatalf("Failed to run migrations up: %v", err)
		}
		fmt.Println("Migrations up completed successfully")

	case *migrateDown:
		if err := runMigrationsDown(database); err != nil {
			log.Fatalf("Failed to run migrations down: %v", err)
		}
		fmt.Println("Migrations down completed successfully")

	case *migrateReset:
		if err := runMigrationsDown(database); err != nil {
			log.Printf("Warning: migrations down failed: %v", err)
		}
		if err := runMigrationsUp(database); err != nil {
			log.Fatalf("Failed to run migrations up: %v", err)
		}
		fmt.Println("Database reset completed successfully")

	case *testDB:
		fmt.Println("Database connection successful!")

	case *seedDB:
		if err := seedDatabase(database); err != nil {
			log.Fatalf("Failed to seed database: %v", err)
		}
		fmt.Println("Database seeding completed successfully")

	default:
		fmt.Println("Database Utility Tool")
		fmt.Println("Usage:")
		fmt.Println("  -migrate-up    Run migrations up")
		fmt.Println("  -migrate-down  Run migrations down") 
		fmt.Println("  -migrate-reset Reset database (down then up)")
		fmt.Println("  -test-db       Test database connection")
		fmt.Println("  -seed-db       Seed database with sample data")
	}
}

func runMigrationsUp(database *sql.DB) error {
	driver, err := postgres.WithInstance(database, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migrate driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func runMigrationsDown(database *sql.DB) error {
	driver, err := postgres.WithInstance(database, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migrate driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations down: %w", err)
	}

	return nil
}

func seedDatabase(database *sql.DB) error {
	fmt.Println("Starting database seeding...")
	
	// Get database URL for GORM connection
	databaseURL := os.Getenv("POSTGRES_URL")
	gormDB, err := gorm.Open(gormpostgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to create GORM connection: %w", err)
	}
	
	// Seed systems data
	if err := seed.SeedSystems(gormDB); err != nil {
		return fmt.Errorf("failed to seed systems: %w", err)
	}
	
	fmt.Println("All seeding completed successfully!")
	return nil
}

func stringPtr(s string) *string {
	return &s
} 