package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/joho/godotenv"

	"go-migrate-example/api/routes"
	"go-migrate-example/migrations"
	"go-migrate-example/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// MigrationFunc represents a migration function signature
type MigrationFunc func(*gorm.DB) error

// Migrations is a map of migration functions
var Migrations = map[string]MigrationFunc{
	"000001_create_users_table": migrations.Up_000001_create_users_table,
	"000004_create_notes_table": migrations.Up_000002_create_notes_table,
}

// DownMigrations is a map of down migration functions
var DownMigrations = map[string]MigrationFunc{
	"000004_create_notes_table": migrations.Down_000002_create_notes_table,
	"000001_create_users_table": migrations.Down_000001_create_users_table,
}

// RunMigrations executes all registered database migrations in order
func RunMigrations(db *gorm.DB) error {
	// Create migration_versions table if not exists
	err := db.Exec(`
        CREATE TABLE IF NOT EXISTS migration_versions (
            version VARCHAR(255) PRIMARY KEY,
            applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `).Error
	if err != nil {
		return fmt.Errorf("failed to create migration_versions table: %v", err)
	}

	// Get sorted migration keys to ensure consistent order
	migrationKeys := make([]string, 0, len(Migrations))
	for k := range Migrations {
		migrationKeys = append(migrationKeys, k)
	}
	sort.Strings(migrationKeys)

	// Run migrations in order
	for _, migrationKey := range migrationKeys {
		// Check if migration has already been applied
		var count int64
		db.Raw("SELECT COUNT(*) FROM migration_versions WHERE version = ?", migrationKey).Scan(&count)

		if count > 0 {
			log.Printf("Migration %s already applied, skipping", migrationKey)
			continue
		}

		migrationFunc := Migrations[migrationKey]
		log.Printf("Applying migration: %s", migrationKey)

		if err := migrationFunc(db); err != nil {
			return fmt.Errorf("failed to apply migration %s: %v", migrationKey, err)
		}

		// Record the applied migration
		err := db.Exec("INSERT INTO migration_versions (version) VALUES (?)", migrationKey).Error
		if err != nil {
			return fmt.Errorf("failed to record migration %s: %v", migrationKey, err)
		}
	}

	log.Println("All migrations applied successfully")
	return nil
}

// RunDownMigrations executes down migrations in reverse order
func RunDownMigrations(db *gorm.DB) error {
	// Get sorted migration keys to ensure consistent order
	migrationKeys := make([]string, 0, len(DownMigrations))
	for k := range DownMigrations {
		migrationKeys = append(migrationKeys, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(migrationKeys)))

	// Run down migrations in reverse order
	for _, migrationKey := range migrationKeys {
		// Check if migration is in the versions table
		var count int64
		db.Raw("SELECT COUNT(*) FROM migration_versions WHERE version = ?", migrationKey).Scan(&count)

		if count == 0 {
			log.Printf("Migration %s not found, skipping", migrationKey)
			continue
		}

		migrationFunc := DownMigrations[migrationKey]
		log.Printf("Rolling back migration: %s", migrationKey)

		if err := migrationFunc(db); err != nil {
			return fmt.Errorf("failed to rollback migration %s: %v", migrationKey, err)
		}

		// Remove the migration from versions table
		err := db.Exec("DELETE FROM migration_versions WHERE version = ?", migrationKey).Error
		if err != nil {
			return fmt.Errorf("failed to remove migration record %s: %v", migrationKey, err)
		}
	}

	log.Println("All migrations rolled back successfully")
	return nil
}

func MigrateLastDown(db *gorm.DB) error {
	// Get sorted migration keys to find the last migration
	migrationKeys := make([]string, 0, len(DownMigrations))
	for k := range DownMigrations {
		migrationKeys = append(migrationKeys, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(migrationKeys)))

	// If there are migrations, rollback only the last one
	if len(migrationKeys) > 0 {
		lastMigrationKey := migrationKeys[0]
		migrationFunc := DownMigrations[lastMigrationKey]
		log.Printf("Rolling back last migration: %s", lastMigrationKey)

		// Add detailed logging before and after migration
		log.Println("Checking existing migration versions before rollback:")
		var existingVersions []models.MigrationVersion
		db.Find(&existingVersions)
		for _, v := range existingVersions {
			log.Printf("Existing version: %s, Applied at: %v", v.Version, v.AppliedAt)
		}

		if err := migrationFunc(db); err != nil {
			return fmt.Errorf("failed to rollback last migration %s: %v", lastMigrationKey, err)
		}

		// Add logging to check migration versions after rollback
		log.Println("Checking migration versions after rollback:")
		existingVersions = []models.MigrationVersion{}
		db.Find(&existingVersions)
		if len(existingVersions) == 0 {
			log.Println("No migration versions remain after rollback")
		} else {
			for _, v := range existingVersions {
				log.Printf("Remaining version: %s, Applied at: %v", v.Version, v.AppliedAt)
			}
		}

		// Explicitly delete the migration version
		result := db.Where("version = ?", lastMigrationKey).Delete(&models.MigrationVersion{})
		log.Printf("Deletion result - Rows affected: %d, Error: %v", result.RowsAffected, result.Error)

		log.Println("Last migration rolled back successfully")
		return nil
	}

	log.Println("No migrations to roll back")
	return nil
}

func MigrateSpecificDown(db *gorm.DB, migrationName string) error {
	// Check if the specific migration exists
	migrationFunc, exists := DownMigrations[migrationName]
	if !exists {
		return fmt.Errorf("migration %s not found", migrationName)
	}

	log.Printf("Rolling back specific migration: %s", migrationName)
	if err := migrationFunc(db); err != nil {
		return fmt.Errorf("failed to rollback migration %s: %v", migrationName, err)
	}

	// Remove the migration version from the database
	result := db.Where("version = ?", migrationName).Delete(&models.MigrationVersion{})
	if result.Error != nil {
		log.Printf("Error removing migration version for %s: %v", migrationName, result.Error)
		return result.Error
	}
	log.Printf("Migration version for %s removed. Rows affected: %d", migrationName, result.RowsAffected)

	log.Println("Specific migration rolled back successfully")
	return nil
}

// MigrateUp applies all migrations
func MigrateUp(db *gorm.DB) error {
	return RunMigrations(db)
}

// MigrateDown rolls back all migrations
func MigrateDown(db *gorm.DB) error {
	return RunDownMigrations(db)
}

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

	// Connect to the database using GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Parse command-line flags
	migrateUp := flag.Bool("up", false, "Apply all migrations")
	migrateDown := flag.Bool("down", false, "Rollback all migrations")
	migrateLastDown := flag.Bool("last-down", false, "Rollback the last migration")
	migrateSpecificDown := flag.String("specific-down", "", "Rollback a specific migration by name")
	flag.Parse()

	// Perform migrations based on flags
	if *migrateUp {
		if err := MigrateUp(db); err != nil {
			log.Fatalf("Migration up failed: %v", err)
		}
		log.Println("Migrations applied successfully")
	} else if *migrateDown {
		if err := MigrateDown(db); err != nil {
			log.Fatalf("Migration down failed: %v", err)
		}
		log.Println("Migrations rolled back successfully")
	} else if *migrateLastDown {
		if err := MigrateLastDown(db); err != nil {
			log.Fatalf("Last migration rollback failed: %v", err)
		}
	} else if *migrateSpecificDown != "" {
		if err := MigrateSpecificDown(db, *migrateSpecificDown); err != nil {
			log.Fatalf("Specific migration rollback failed: %v", err)
		}
	} else {
		log.Println("No migration action specified. Use -up or -down flag, -last-down, or -specific-down flag.")

		// Setup and run the server
		r := routes.SetupRoutes(db)
		if err := r.Run(":8080"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}
}
