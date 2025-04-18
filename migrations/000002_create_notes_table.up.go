package migrations

import (
	"time"

	"gorm.io/gorm"
)

func Up_000002_create_notes_table(db *gorm.DB) error {
	return db.AutoMigrate(&Notes{})
}

type Notes struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"size:255" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedBy   uint           `gorm:"not null" json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedBy   uint           `gorm:"not null" json:"updated_by"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
