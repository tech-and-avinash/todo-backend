package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbInstance *gorm.DB

func InitDB() (*gorm.DB, error) {
	if dbInstance != nil {
		return dbInstance, nil
	}

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

	log.Printf("DB Connection Details:")
	log.Printf("Host: %s", host)
	log.Printf("Port: %s", port)
	log.Printf("User: %s", user)
	log.Printf("Database: %s", dbname)

	if host == "" || port == "" || user == "" || dbname == "" {
		return nil, fmt.Errorf("missing required database environment variables")
	}

	// ✅ Disable statement caching at driver level using preferSimpleProtocol
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable application_name=nomadule_app preferSimpleProtocol=true",
		host, port, user, password, dbname)

	// ✅ Disable GORM-level prepared statements
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		PrepareStmt: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v\nConnection string: %s", err, dsn)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("cannot get *sql.DB from GORM: %w", err)
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	dbInstance = db
	return dbInstance, nil
}
