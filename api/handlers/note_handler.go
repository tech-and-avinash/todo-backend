package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"go-migrate-example/models"
	"go-migrate-example/repositories"

	"github.com/gin-gonic/gin"
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
	// 1. Get token from header
	token := c.GetHeader("Authorization")
	if token == "" {
		log.Println("CreateNote: Missing Authorization header")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
		return
	}
	log.Printf("CreateNote: Received token - %s", token)

	// 2. Extract user ID from token (you need to implement this based on your auth strategy)
	userID, err := extractUserIDFromToken(c)
	if err != nil {
		log.Printf("CreateNote: Failed to extract user ID - Error: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	log.Printf("CreateNote: Extracted User ID - %d", userID)

	// 3. Bind JSON input to Note model
	var note models.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		log.Printf("CreateNote: Failed to bind JSON - Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("CreateNote: Received Note - Title: %s, Description: %s", note.Title, note.Description)

	// 4. Set user info and timestamps
	note.CreatedBy = userID
	note.UpdatedBy = userID
	note.CreatedAt = time.Now()
	note.UpdatedAt = time.Now()
	log.Printf("CreateNote: Prepared Note - CreatedBy: %d, UpdatedBy: %d", note.CreatedBy, note.UpdatedBy)

	// 5. Save note using repository
	if err := h.repo.Create(&note); err != nil {
		log.Printf("CreateNote: Failed to create note - Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create note"})
		return
	}

	log.Printf("CreateNote: Successfully created note - ID: %d", note.ID)
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

	noteIDParam := c.Param("id")
	noteID, err := strconv.ParseUint(noteIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	note, err := h.repo.GetByID(uint(noteID))
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
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
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

	note.ID = uint(id)
	if err := h.repo.Update(&note); err != nil {
		log.Printf("UpdateNote: Failed to update note - Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("UpdateNote: Successfully updated note - ID: %d", note.ID)
	c.JSON(http.StatusOK, note)
}

func (h *NoteHandler) DeleteNote(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Printf("DeleteNote: Failed to parse note ID - Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		log.Printf("DeleteNote: Failed to delete note - Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("DeleteNote: Successfully deleted note - ID: %d", id)
	c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}
