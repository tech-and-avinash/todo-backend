package models

import (
	"time"

	"gorm.io/gorm"
)

// User model
type User struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:100;not null"`
	Email     string `gorm:"unique;not null"`
	Password  string `gorm:"size:255;not null"`
	Age       int    `gorm:"column:age"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Todo model
type Note struct {
	ID          uint           `gorm:"primaryKey;column:id" json:"id"`
	Title       string         `gorm:"size:255;column:title" json:"title"`
	Description string         `gorm:"type:text;column:description" json:"description"`
	CreatedBy   uint           `gorm:"not null;column:created_by" json:"created_by"`
	CreatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP;column:created_at" json:"created_at"`
	UpdatedBy   uint           `gorm:"not null;column:updated_by" json:"updated_by"`
	UpdatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP;column:updated_at" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index;column:deleted_at" json:"deleted_at"`
}

// MigrationVersion tracks the applied migrations
type MigrationVersion struct {
	Version   string    `gorm:"primaryKey"`
	AppliedAt time.Time `gorm:"not null"`
}
