package migrations

import (
	"gorm.io/gorm"
)

func Down_000001_create_users_table(db *gorm.DB) error {
	return db.Migrator().DropTable(&User{})
}
