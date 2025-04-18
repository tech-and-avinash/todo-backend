package migrations

import (
	"gorm.io/gorm"
)

func Down_000002_create_notes_table(db *gorm.DB) error {
	return db.Migrator().DropTable(&Notes{})
}
