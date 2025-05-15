package repositories

import (
	"fmt"
	"nomadule-backend/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ContactRepository struct {
	db *gorm.DB
}

func NewContactRepository(db *gorm.DB) *ContactRepository {
	return &ContactRepository{db: db}
}

// Create a new contact
func (r *ContactRepository) Create(contact *models.Contact) error {
	fmt.Printf("Saving to DB: %+v\n", contact)
	err := r.db.Create(contact).Error
	fmt.Printf("Saved Contact with ID: %s, Error: %v\n", contact.ID, err)
	return err
}

// Get all contacts created by a specific user (string user ID)
func (r *ContactRepository) GetAllByUser(userID string) ([]models.Contact, error) {
	var contacts []models.Contact
	err := r.db.Where("created_by = ?", userID).Find(&contacts).Error
	return contacts, err
}

// Get contact by UUID
func (r *ContactRepository) GetByID(id uuid.UUID) (*models.Contact, error) {
	var contact models.Contact
	err := r.db.First(&contact, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &contact, nil
}

// Update a contact
func (r *ContactRepository) Update(contact *models.Contact) error {
	return r.db.Save(contact).Error
}

// Delete a contact by UUID
func (r *ContactRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Contact{}, "id = ?", id).Error
}
