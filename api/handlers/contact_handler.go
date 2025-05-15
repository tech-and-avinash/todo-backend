package handlers

import (
	"log"
	"net/http"
	"nomadule-backend/models"
	"nomadule-backend/repositories"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ContactHandler struct {
	repo *repositories.ContactRepository
}

func NewContactHandler(db *gorm.DB) *ContactHandler {
	return &ContactHandler{
		repo: repositories.NewContactRepository(db),
	}
}

func (h *ContactHandler) CreateContact(c *gin.Context) {
	user, err := extractUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var req models.CreateContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid contact input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	contactID := uuid.New()
	contact := models.Contact{
		ID:        contactID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Address:   req.Address,
		CreatedBy: user.ID,
		UpdatedBy: user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := h.repo.Create(&contact); err != nil {
		log.Printf("Failed to create contact: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create contact"})
		return
	}
	c.JSON(http.StatusCreated, contact)
}

func (h *ContactHandler) GetAllContacts(c *gin.Context) {
	user, err := extractUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	contacts, err := h.repo.GetAllByUser(user.ID)
	if err != nil {
		log.Printf("Failed to get contacts: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve contacts"})
		return
	}
	c.JSON(http.StatusOK, contacts)
}

func (h *ContactHandler) GetContact(c *gin.Context) {
	user, err := extractUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID"})
		return
	}

	contact, err := h.repo.GetByID(id)
	if err != nil || contact.CreatedBy != user.ID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
		return
	}

	c.JSON(http.StatusOK, contact)
}

func (h *ContactHandler) UpdateContact(c *gin.Context) {
	user, err := extractUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID"})
		return
	}

	var req models.UpdateContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contact, err := h.repo.GetByID(id)
	if err != nil || contact.CreatedBy != user.ID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
		return
	}

	// Update fields
	contact.FirstName = req.FirstName
	contact.LastName = req.LastName
	contact.Email = req.Email
	contact.Phone = req.Phone
	contact.Address = req.Address
	contact.UpdatedBy = user.ID
	contact.UpdatedAt = time.Now()

	if err := h.repo.Update(contact); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update contact"})
		return
	}

	c.JSON(http.StatusOK, contact)
}

func (h *ContactHandler) DeleteContact(c *gin.Context) {
	user, err := extractUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID"})
		return
	}

	contact, err := h.repo.GetByID(id)
	if err != nil || contact.CreatedBy != user.ID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
		return
	}

	if err := h.repo.Delete(id); err != nil {
		log.Printf("Failed to delete contact: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete contact"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contact deleted successfully"})
}
