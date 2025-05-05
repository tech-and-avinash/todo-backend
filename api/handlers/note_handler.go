package handlers

import (
	"errors"
	"log"
	"net/http"
	"time"

	"todo-backend/models"
	"todo-backend/repositories"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NoteHandler struct {
	repo     *repositories.NoteRepository
	userRepo *repositories.UserRepository
}

func NewNoteHandler(db *gorm.DB) *NoteHandler {
	return &NoteHandler{
		repo:     repositories.NewNoteRepository(db),
		userRepo: repositories.NewUserRepository(db),
	}
}

func (h *NoteHandler) CreateNote(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		log.Println("CreateNote: Missing Authorization header")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
		return
	}

	userID, err := extractUserIDFromToken(c)
	if err != nil {
		log.Printf("CreateNote: Failed to extract user ID - Error: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	var note models.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		log.Printf("CreateNote: Failed to bind JSON - Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note.CreatedBy = userID
	note.UpdatedBy = userID
	note.CreatedAt = time.Now()
	note.UpdatedAt = time.Now()

	if err := h.repo.Create(&note); err != nil {
		log.Printf("CreateNote: Failed to create note - Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create note"})
		return
	}

	c.JSON(http.StatusCreated, note)
}

func (h *NoteHandler) GetAllNotes(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
		return
	}

	userID, err := extractUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	notes, err := h.repo.GetAllByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notes"})
		return
	}

	c.JSON(http.StatusOK, notes)
}

func (h *NoteHandler) GetNoteByID(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
		return
	}

	userID, err := extractUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	noteIDStr := c.Param("id")
	noteID, err := uuid.Parse(noteIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	note, err := h.repo.GetByID(noteID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	if note.CreatedBy != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to access this note"})
		return
	}

	c.JSON(http.StatusOK, note)
}

func (h *NoteHandler) UpdateNote(c *gin.Context) {
	noteIDStr := c.Param("id")
	noteID, err := uuid.Parse(noteIDStr)
	if err != nil {
		log.Printf("UpdateNote: Failed to parse note ID - Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	var note models.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		log.Printf("UpdateNote: Failed to bind JSON - Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := extractUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	note.ID = noteID
	note.UpdatedBy = userID
	note.UpdatedAt = time.Now()

	if err := h.repo.Update(&note); err != nil {
		log.Printf("UpdateNote: Failed to update note - Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, note)
}

func (h *NoteHandler) DeleteNote(c *gin.Context) {
	noteIDStr := c.Param("id")
	noteID, err := uuid.Parse(noteIDStr)
	if err != nil {
		log.Printf("DeleteNote: Failed to parse note ID - Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	if err := h.repo.Delete(noteID); err != nil {
		log.Printf("DeleteNote: Failed to delete note - Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}

func extractUserIDFromToken(c *gin.Context) (uuid.UUID, error) {
	val, exists := c.Get("userID")
	if !exists {
		return uuid.Nil, errors.New("user ID not found in context")
	}

	userID, ok := val.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("invalid user ID type in context")
	}

	return userID, nil
}
