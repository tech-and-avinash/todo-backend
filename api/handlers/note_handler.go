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

	if err := h.repo.Create(&note); err != nil {
		log.Printf("Failed to create note: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create note"})
		return
	}

	for _, item := range req.ChecklistItems {
		newItem := models.ChecklistItem{
			ID:        uuid.New(),
			NoteID:    noteID, // ✅ link to the new note
			Text:      item.Text,
			IsChecked: item.IsChecked,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		h.repo.CreateChecklistItem(&newItem)
	}

	for _, r := range req.Reminders {
		reminder := models.Reminder{
			ID:     uuid.New(),
			NoteID: noteID,
			Time:   r.Time,
		}
		h.repo.CreateReminder(&reminder)
	}

	for _, a := range req.Attachments {
		attachment := models.NoteAttachment{
			ID:          uuid.New(),
			NoteID:      noteID,
			FileName:    a.FileName,
			URL:         a.URL,
			ContentType: a.ContentType,
			CreatedAt:   time.Now(),
		}
		h.repo.CreateAttachment(&attachment)
	}

	attachments, _ := h.repo.GetAttachmentsByNoteID(noteID)
	var attachmentResponses []models.NoteAttachmentResponse
	for _, a := range attachments {
		attachmentResponses = append(attachmentResponses, models.NoteAttachmentResponse{
			FileName:    a.FileName,
			URL:         a.URL,
			ContentType: a.ContentType,
		})
	}

	response := models.NoteResponse{
		ID:             note.ID,
		Title:          note.Title,
		Description:    note.Description,
		FirstName:      user.FirstName,
		IsChecklist:    note.IsChecklist,
		CreatedAt:      note.CreatedAt,
		UpdatedAt:      note.UpdatedAt,
		CreatedBy:      note.CreatedBy,
		UpdatedBy:      note.UpdatedBy,
		ChecklistItems: []models.ChecklistItemResponse{},
		Reminders:      []models.ReminderResponse{},
		Attachments:    attachmentResponses,
	}

	c.JSON(http.StatusCreated, response)
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

	h.repo.DeleteChecklistItemsByNote(noteID)
	h.repo.DeleteRemindersByNote(noteID)
	h.repo.DeleteAttachmentsByNoteID(noteID)

	for _, item := range req.ChecklistItems {
		newItem := models.ChecklistItem{
			ID:        uuid.New(),
			NoteID:    noteID, // ✅ link to the new note
			Text:      item.Text,
			IsChecked: item.IsChecked,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		h.repo.CreateChecklistItem(&newItem)
	}

	for _, r := range req.Reminders {
		reminder := models.Reminder{
			ID:     uuid.New(),
			NoteID: noteID,
			Time:   r.Time,
		}
		h.repo.CreateReminder(&reminder)
	}

	for _, a := range req.Attachments {
		attachment := models.NoteAttachment{
			ID:          uuid.New(),
			NoteID:      noteID,
			FileName:    a.FileName,
			URL:         a.URL,
			ContentType: a.ContentType,
			CreatedAt:   time.Now(),
		}
		h.repo.CreateAttachment(&attachment)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note updated successfully"})
}

func (h *NoteHandler) GetAllNotes(c *gin.Context) {
	user, err := extractUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	notes, err := h.repo.GetAllByUser(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notes"})
		return
	}

	var response []models.NoteResponse
	for _, note := range notes {
		checklistItems, _ := h.repo.GetChecklistItemsByNoteID(note.ID)
		reminders, _ := h.repo.GetRemindersByNoteID(note.ID)
		attachments, _ := h.repo.GetAttachmentsByNoteID(note.ID)

		var checklist []models.ChecklistItemResponse
		for _, item := range checklistItems {
			checklist = append(checklist, models.ChecklistItemResponse{
				ID:        item.ID,
				Text:      item.Text,
				IsChecked: item.IsChecked,
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
			})
		}

		var reminderResponses []models.ReminderResponse
		for _, r := range reminders {
			reminderResponses = append(reminderResponses, models.ReminderResponse{Time: r.Time})
		}

		var attachmentResponses []models.NoteAttachmentResponse
		for _, a := range attachments {
			attachmentResponses = append(attachmentResponses, models.NoteAttachmentResponse{
				FileName:    a.FileName,
				URL:         a.URL,
				ContentType: a.ContentType,
			})
		}

		response = append(response, models.NoteResponse{
			ID:             note.ID,
			Title:          note.Title,
			Description:    note.Description,
			FirstName:      user.FirstName,
			CreatedAt:      note.CreatedAt,
			UpdatedAt:      note.UpdatedAt,
			CreatedBy:      note.CreatedBy,
			UpdatedBy:      note.UpdatedBy,
			ChecklistItems: checklist,
			Reminders:      reminderResponses,
			Attachments:    attachmentResponses,
		})
	}

	c.JSON(http.StatusOK, gin.H{"notes": response, "total": len(response)})
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
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	checklistItems, _ := h.repo.GetChecklistItemsByNoteID(note.ID)
	reminders, _ := h.repo.GetRemindersByNoteID(note.ID)
	attachments, _ := h.repo.GetAttachmentsByNoteID(note.ID)

	var checklist []models.ChecklistItemResponse
	for _, item := range checklistItems {
		checklist = append(checklist, models.ChecklistItemResponse{
			ID:        item.ID,
			Text:      item.Text,
			IsChecked: item.IsChecked,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		})
	}

	var reminderResponses []models.ReminderResponse
	for _, r := range reminders {
		reminderResponses = append(reminderResponses, models.ReminderResponse{Time: r.Time})
	}

	var attachmentResponses []models.NoteAttachmentResponse
	for _, a := range attachments {
		attachmentResponses = append(attachmentResponses, models.NoteAttachmentResponse{
			FileName:    a.FileName,
			URL:         a.URL,
			ContentType: a.ContentType,
		})
	}

	c.JSON(http.StatusOK, models.NoteResponse{
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
		Reminders:      reminderResponses,
		Attachments:    attachmentResponses,
	})
}

func (h *NoteHandler) DeleteNote(c *gin.Context) {
	noteID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	h.repo.DeleteChecklistItemsByNote(noteID)
	h.repo.DeleteRemindersByNote(noteID)
	h.repo.DeleteAttachmentsByNoteID(noteID)

	if err := h.repo.Delete(noteID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete note"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}
