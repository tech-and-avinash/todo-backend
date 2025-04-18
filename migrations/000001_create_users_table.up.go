package migrations

import (
	"time"

	"gorm.io/gorm"
)

func Up_000001_create_users_table(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:100;not null"`
	Email     string `gorm:"unique;not null"`
	Age       int    `gorm:"column:age"`
	Password  string `gorm:"size:255;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
