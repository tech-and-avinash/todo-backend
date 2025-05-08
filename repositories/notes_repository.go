package repositories

import (
	"fmt"
	"todo-backend/models"

	"github.com/google/uuid"
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
	fmt.Printf("Saved Note with ID: %s, Error: %v\n", note.ID, err)
	return err
}

// Get all notes created by a specific user (UUID)
func (r *NoteRepository) GetAllByUser(userID string) ([]models.Note, error) {
	var notes []models.Note
	err := r.db.Where("created_by = ?", userID).Find(&notes).Error
	return notes, err
}

// Get note by ID (UUID)
func (r *NoteRepository) GetByID(id uuid.UUID) (*models.Note, error) {
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
func (r *NoteRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Note{}, "id = ?", id).Error
}

// List all notes
func (r *NoteRepository) GetAll() ([]models.Note, error) {
	var notes []models.Note
	err := r.db.Find(&notes).Error
	return notes, err
}

// CreateChecklistItem saves a checklist item linked to a note
func (r *NoteRepository) CreateChecklistItem(item *models.ChecklistItem) error {
	return r.db.Create(item).Error
}

// CreateReminder saves a reminder linked to a note
func (r *NoteRepository) CreateReminder(reminder *models.Reminder) error {
	return r.db.Create(reminder).Error
}

// DeleteChecklistItemsByNote deletes all checklist items for a note (used in update)
func (r *NoteRepository) DeleteChecklistItemsByNote(noteID uuid.UUID) error {
	return r.db.Where("note_id = ?", noteID).Delete(&models.ChecklistItem{}).Error
}

// DeleteRemindersByNote deletes all reminders for a note (used in update)
func (r *NoteRepository) DeleteRemindersByNote(noteID uuid.UUID) error {
	return r.db.Where("note_id = ?", noteID).Delete(&models.Reminder{}).Error
}

func (r *NoteRepository) GetChecklistItemsByNoteID(noteID uuid.UUID) ([]models.ChecklistItem, error) {
	var items []models.ChecklistItem
	err := r.db.Where("note_id = ?", noteID).Find(&items).Error
	return items, err
}

func (r *NoteRepository) GetRemindersByNoteID(noteID uuid.UUID) ([]models.Reminder, error) {
	var reminders []models.Reminder
	err := r.db.Where("note_id = ?", noteID).Find(&reminders).Error
	return reminders, err
}
