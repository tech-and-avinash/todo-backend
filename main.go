package main

import (
	"log"

	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"nomadule-backend/api/routes"
	"nomadule-backend/azure"
	"nomadule-backend/config"
	"nomadule-backend/models"
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

	// ✅ Initialize Azure Storage Client
	azureClient, err := azure.NewAzureStorageClient()
	if err != nil {
		log.Fatalf("Failed to initialize Azure client: %v", err)
	}
	// Run database migrations...
	// err = Migrate(db)
	// if err != nil {
	// 	log.Fatalf("could not migrate database: %v", err)
	// 	return
	// }

	// ✅ Pass Azure client into route setup
	r := routes.SetupRoutes(db, azureClient)

	if err := r.Run("0.0.0.0:8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

//	func AutoMigrate(db *gorm.DB) error {
//		modelsToMigrate := []interface{}{
func Migrate(DB *gorm.DB) error {
	err := DB.AutoMigrate(&models.User{}, &models.Note{}, &models.ChecklistItem{}, &models.Reminder{}, &models.NoteAttachment{}, &models.Contact{})
	if err != nil {
		log.Fatalf("could not migrate database: %v", err)
	}

	log.Println("✅ All migrations completed successfully.")
	return nil
}
