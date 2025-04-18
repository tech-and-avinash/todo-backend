package repositories

import (
	"fmt"
	"go-migrate-example/models"

	"gorm.io/gorm"
)

type NoteRepository struct {
	db *gorm.DB
}

func NewNoteRepository(db *gorm.DB) *NoteRepository {
	return &NoteRepository{db: db}
}

// Create a new note
func (r *NoteRepository) Create(note *models.Note) error {
	fmt.Printf("Saving to DB: %+v\n", note)
	err := r.db.Create(note).Error
	fmt.Printf("Saved Note with ID: %d, Error: %v\n", note.ID, err)
	return err
}

func (r *NoteRepository) GetAllByUser(userID uint) ([]models.Note, error) {
	var notes []models.Note
	err := r.db.Where("created_by = ?", userID).Find(&notes).Error
	return notes, err
}

func (r *NoteRepository) GetByID(id uint) (*models.Note, error) {
	var note models.Note
	err := r.db.First(&note, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &note, nil
}

// Update a note
func (r *NoteRepository) Update(note *models.Note) error {
	return r.db.Save(note).Error
}

// Delete a note (soft delete)
func (r *NoteRepository) Delete(id uint) error {
	return r.db.Delete(&models.Note{}, id).Error
}

// List all notes with associated users
func (r *NoteRepository) GetAll() ([]models.Note, error) {
	var notes []models.Note
	err := r.db.Preload("User").Find(&notes).Error
	return notes, err
}
