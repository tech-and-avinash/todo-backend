package models

import (
	"time"

	"github.com/google/uuid"
)

type CreateUserRequest struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName"`
	ClerkID   string `json:"clerkId"`
	Email     string `json:"email" binding:"required,email"`
	ImageUrl  string `json:"image_url"`
	Password  string `json:"password"`
}

// request model
type CreateNoteRequest struct {
	Title          string            `json:"title" binding:"required"`
	Description    string            `json:"description"`
	IsPinned       bool              `json:"isPinned"`
	IsArchived     bool              `json:"isArchived"`
	IsChecklist    bool              `json:"isChecklist"`
	ChecklistItems []ChecklistItem   `json:"checklistItems"`
	Reminders      []ReminderRequest `json:"reminders"`
}

// response model
type NoteResponse struct {
	ID             uuid.UUID               `json:"id"`
	Title          string                  `json:"title"`
	Description    string                  `json:"description"`
	FirstName      *string                 `json:"user_firstname"`
	IsPinned       bool                    `json:"isPinned"`
	IsArchived     bool                    `json:"isArchived"`
	IsChecklist    bool                    `json:"isChecklist"`
	CreatedAt      time.Time               `json:"created_at"`
	UpdatedAt      time.Time               `json:"updated_at"`
	CreatedBy      string                  `json:"created_by"`
	UpdatedBy      string                  `json:"updated_by"`
	ChecklistItems []ChecklistItemResponse `json:"checklist_items,omitempty"`
	Reminders      []ReminderResponse      `json:"reminders,omitempty"`
}

type NoteListResponse struct {
	Total int            `json:"total"`
	Notes []NoteResponse `json:"notes"`
}
type ChecklistItemResponse struct {
	ID        uuid.UUID `json:"id"`
	Text      string    `json:"text"`
	IsChecked bool      `json:"isChecked"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ReminderRequest struct {
	Time time.Time `json:"time"`
}

type ReminderResponse struct {
	Time time.Time `json:"time"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LogoutRequest struct {
	Token string `json:"token" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}
