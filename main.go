package main

import (
	"log"

	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"todo-backend/api/routes"
	"todo-backend/config"
	"todo-backend/models"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default or system environment variables")
	}

	// Initialize database
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := AutoMigrate(db); err != nil {
		log.Printf("⚠️ Migration warning: %v", err)
	}

	// Setup and run the server
	r := routes.SetupRoutes(db)
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func AutoMigrate(db *gorm.DB) error {
	for _, model := range []interface{}{
		&models.User{},
		&models.Note{},
		&models.ChecklistItem{},
		&models.Reminder{},
	} {
		log.Printf("Migrating: %T", model)
		if err := db.Migrator().AutoMigrate(model); err != nil {
			log.Printf("⚠️ Failed to migrate %T: %v", model, err)
			return err
		}
	}
	log.Println("✅ All migrations attempted.")
	return nil
}
