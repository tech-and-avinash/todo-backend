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

	// if err := AutoMigrate(db); err != nil {
	// 	log.Printf("⚠️ Migration warning: %v", err)
	// }

	// Run database migrations...
	err = Migrate(db)
	if err != nil {
		log.Fatalf("could not migrate database: %v", err)
		return
	}

	// Setup and run the server
	r := routes.SetupRoutes(db)
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

//	func AutoMigrate(db *gorm.DB) error {
//		modelsToMigrate := []interface{}{
func Migrate(DB *gorm.DB) error {
	err := DB.AutoMigrate(&models.User{}, &models.Note{}, &models.ChecklistItem{}, &models.Reminder{})
	if err != nil {
		log.Fatalf("could not migrate database: %v", err)
	}

	log.Println("✅ All migrations completed successfully.")
	return nil
}
