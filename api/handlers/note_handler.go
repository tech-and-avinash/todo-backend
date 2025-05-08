package handlers

import (
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
	user, err := extractUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req models.CreateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid note input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	noteID := uuid.New()
	note := models.Note{
		ID:          noteID,
		Title:       req.Title,
		Description: req.Description,
		IsPinned:    req.IsPinned,
		IsArchived:  req.IsArchived,
		IsChecklist: req.IsChecklist,
		CreatedBy:   user.ID,
		UpdatedBy:   user.ID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save note
	if err := h.repo.Create(&note); err != nil {
		log.Printf("Failed to create note: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create note"})
		return
	}

	// Save checklist items
	for _, item := range req.ChecklistItems {
		itemID := uuid.New()
		newItem := models.ChecklistItem{
			ID:        itemID,
			NoteID:    noteID,
			Text:      item.Text,
			IsChecked: item.IsChecked,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := h.repo.CreateChecklistItem(&newItem); err != nil {
			log.Printf("Checklist create error: %v", err)
		}
	}

	// Save reminders
	for _, r := range req.Reminders {
		reminder := models.Reminder{
			ID:     uuid.New(),
			NoteID: noteID,
			Time:   r.Time, // this now works because r is time.Time
		}
		if err := h.repo.CreateReminder(&reminder); err != nil {
			log.Printf("Reminder create error: %v", err)
		}
	}

	// Fetch the checklist and reminders again to include in response
	checklistModel, _ := h.repo.GetChecklistItemsByNoteID(noteID)
	var checklist []models.ChecklistItemResponse
	for _, item := range checklistModel {
		checklist = append(checklist, models.ChecklistItemResponse{
			ID:        item.ID,
			Text:      item.Text,
			IsChecked: item.IsChecked,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		})
	}

	remindersModel, _ := h.repo.GetRemindersByNoteID(noteID)

	var reminders []models.ReminderResponse
	for _, r := range remindersModel {
		reminders = append(reminders, models.ReminderResponse{
			Time: r.Time,
		})
	}

	// Build response
	response := models.NoteResponse{
		ID:             note.ID,
		Title:          note.Title,
		Description:    note.Description,
		FirstName:      user.FirstName,
		CreatedAt:      note.CreatedAt,
		UpdatedAt:      note.UpdatedAt,
		CreatedBy:      note.CreatedBy,
		UpdatedBy:      note.UpdatedBy,
		ChecklistItems: checklist,
		Reminders:      reminders,
	}
	c.JSON(http.StatusCreated, response)
}

func (h *NoteHandler) GetAllNotes(c *gin.Context) {
	user, err := extractUserFromHeader(c)
	if err != nil {
		log.Printf("Auth error: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	log.Printf("Fetching notes for user ID: %s", user.ID)

	notes, err := h.repo.GetAllByUser(user.ID)
	if err != nil {
		log.Printf("Failed to fetch notes for user %s: %v", user.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch notes"})
		return
	}

	var response []models.NoteResponse
	for _, n := range notes {
		items, _ := h.repo.GetChecklistItemsByNoteID(n.ID)
		var checklist []models.ChecklistItemResponse
		for _, item := range items {
			checklist = append(checklist, models.ChecklistItemResponse{
				ID:        item.ID,
				Text:      item.Text,
				IsChecked: item.IsChecked,
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
			})
		}

		remindersModel, _ := h.repo.GetRemindersByNoteID(n.ID)
		var reminders []models.ReminderResponse
		for _, r := range remindersModel {
			reminders = append(reminders, models.ReminderResponse{
				Time: r.Time,
			})
		}

		response = append(response, models.NoteResponse{
			ID:             n.ID,
			Title:          n.Title,
			Description:    n.Description,
			IsPinned:       n.IsPinned,
			IsArchived:     n.IsArchived,
			IsChecklist:    n.IsChecklist,
			CreatedAt:      n.CreatedAt,
			UpdatedAt:      n.UpdatedAt,
			CreatedBy:      n.CreatedBy,
			UpdatedBy:      n.UpdatedBy,
			FirstName:      user.FirstName,
			ChecklistItems: checklist,
			Reminders:      reminders,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"Notes": response,
		"Total": len(response),
	})
}

func (h *NoteHandler) GetNoteByID(c *gin.Context) {
	user, err := extractUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	noteID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	note, err := h.repo.GetByID(noteID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	if note.CreatedBy != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not allowed to access this note"})
		return
	}

	items, _ := h.repo.GetChecklistItemsByNoteID(noteID)
	var checklist []models.ChecklistItemResponse
	for _, item := range items {
		checklist = append(checklist, models.ChecklistItemResponse{
			ID:        item.ID,
			Text:      item.Text,
			IsChecked: item.IsChecked,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		})
	}

	remindersModel, _ := h.repo.GetRemindersByNoteID(noteID)
	var reminders []models.ReminderResponse
	for _, r := range remindersModel {
		reminders = append(reminders, models.ReminderResponse{
			Time: r.Time,
		})
	}

	response := models.NoteResponse{
		ID:             note.ID,
		Title:          note.Title,
		Description:    note.Description,
		IsPinned:       note.IsPinned,
		IsArchived:     note.IsArchived,
		IsChecklist:    note.IsChecklist,
		CreatedAt:      note.CreatedAt,
		UpdatedAt:      note.UpdatedAt,
		CreatedBy:      note.CreatedBy,
		UpdatedBy:      note.UpdatedBy,
		FirstName:      user.FirstName,
		ChecklistItems: checklist,
		Reminders:      reminders,
	}

	c.JSON(http.StatusOK, response)
}

func (h *NoteHandler) UpdateNote(c *gin.Context) {
	user, err := extractUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	noteID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	var req models.CreateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update note
	note := models.Note{
		ID:          noteID,
		Title:       req.Title,
		Description: req.Description,
		IsPinned:    req.IsPinned,
		IsArchived:  req.IsArchived,
		IsChecklist: req.IsChecklist,
		UpdatedBy:   user.ID,
		UpdatedAt:   time.Now(),
	}

	if err := h.repo.Update(&note); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update note"})
		return
	}

	// Clear old checklist/reminders
	h.repo.DeleteChecklistItemsByNote(noteID)
	h.repo.DeleteRemindersByNote(noteID)

	// Add updated checklist items
	for _, item := range req.ChecklistItems {
		newItem := models.ChecklistItem{
			ID:        uuid.New(),
			NoteID:    noteID,
			Text:      item.Text,
			IsChecked: item.IsChecked,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		h.repo.CreateChecklistItem(&newItem)
	}

	// Add updated reminders
	for _, r := range req.Reminders {
		reminder := models.Reminder{
			ID:     uuid.New(),
			NoteID: noteID,
			Time:   r.Time,
		}
		h.repo.CreateReminder(&reminder)
	}

	c.JSON(http.StatusOK, note)
}

func (h *NoteHandler) DeleteNote(c *gin.Context) {
	noteID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	// Delete related checklist and reminders first
	h.repo.DeleteChecklistItemsByNote(noteID)
	h.repo.DeleteRemindersByNote(noteID)

	// Then delete the note
	if err := h.repo.Delete(noteID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete note"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}
