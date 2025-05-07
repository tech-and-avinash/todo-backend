// config/database.go
package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using environment variables")
	}

	// PostgreSQL connection parameters
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Log environment variables for debugging
	log.Printf("DB Connection Details:")
	log.Printf("Host: %s", host)
	log.Printf("Port: %s", port)
	log.Printf("User: %s", user)
	log.Printf("Database: %s", dbname)

	// Validate required environment variables
	if host == "" {
		return nil, fmt.Errorf("DB_HOST is not set")
	}
	if port == "" {
		return nil, fmt.Errorf("DB_PORT is not set")
	}
	if user == "" {
		return nil, fmt.Errorf("DB_USER is not set")
	}
	if dbname == "" {
		return nil, fmt.Errorf("DB_NAME is not set")
	}

	// Connection string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v\nConnection string: %s", err, dsn)
	}
	// Enable uuid-ossp extension
	err = db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error
	if err != nil {
		log.Printf("⚠️ Failed to enable uuid-ossp extension: %v", err)
	}

	return db, nil
}
