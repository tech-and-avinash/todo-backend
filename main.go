package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"go-migrate-example/api/routes"
	"go-migrate-example/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default or system environment variables")
	}

	// Database connection parameters from environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Construct DSN with fallback values
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)

	// Print the DSN
	fmt.Println(dsn)

	// Connect to the database using GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// Run migrations
	if err := db.AutoMigrate(&models.User{}, &models.Note{}); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)

		// Setup and run the server
		r := routes.SetupRoutes(db)
		if err := r.Run(":8080"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}
}
