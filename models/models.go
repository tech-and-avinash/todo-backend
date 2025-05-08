package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ClerkID   string    `json:"clerkId" gorm:"uniqueIndex"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	ImageURL  string    `json:"imageUrl"`
	Password  string    `gorm:"size:255" json:"password"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Note struct {
	ID             uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Title          string          `gorm:"size:255" json:"title"`
	Description    string          `gorm:"type:text" json:"description"`
	IsPinned       bool            `json:"isPinned"`
	IsArchived     bool            `json:"isArchived"`
	IsChecklist    bool            `json:"isChecklist"`
	ChecklistItems []ChecklistItem `gorm:"foreignKey:NoteID" json:"checklistItems"`
	Reminders      []Reminder      `gorm:"foreignKey:NoteID" json:"reminders"`
	CreatedBy      string          `gorm:"not null" json:"created_by"`
	CreatedAt      time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedBy      string          `gorm:"not null" json:"updated_by"`
	UpdatedAt      time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt      gorm.DeletedAt  `gorm:"index" json:"deleted_at"`
}

type ChecklistItem struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	NoteID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Note      Note      `gorm:"foreignKey:NoteID;references:ID"`
	Text      string    `json:"text"`
	IsChecked bool      `json:"isChecked"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Reminder struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	NoteID uuid.UUID `gorm:"type:uuid;not null;index"`
	Note   Note      `gorm:"foreignKey:NoteID;references:ID"`
	Time   time.Time `json:"time"`
}
