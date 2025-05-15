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

type UpdateUserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	ImageURL  string `json:"imageUrl"`
}

type CreateNoteRequest struct {
	Title          string                 `json:"title" binding:"required"`
	Description    string                 `json:"description"`
	IsPinned       bool                   `json:"isPinned"`
	IsArchived     bool                   `json:"isArchived"`
	IsChecklist    bool                   `json:"isChecklist"`
	Attachments    []NoteAttachmentInput  `json:"attachments"`
	ChecklistItems []ChecklistItemRequest `json:"checklistItems"`
	Reminders      []ReminderRequest      `json:"reminders"`
}

type ChecklistItemRequest struct {
	Text      string `json:"text" binding:"required"`
	IsChecked bool   `json:"isChecked"`
}

type ReminderRequest struct {
	Time time.Time `json:"time"`
}
type NoteAttachmentInput struct {
	FileName    string `json:"fileName"`
	URL         string `json:"url"`
	ContentType string `json:"contentType"`
}

type NoteResponse struct {
	ID             uuid.UUID                `json:"id"`
	Title          string                   `json:"title"`
	Description    string                   `json:"description"`
	FirstName      *string                  `json:"user_firstname"`
	IsPinned       bool                     `json:"isPinned"`
	IsArchived     bool                     `json:"isArchived"`
	IsChecklist    bool                     `json:"isChecklist"`
	Attachments    []NoteAttachmentResponse `json:"attachments,omitempty"`
	ChecklistItems []ChecklistItemResponse  `json:"checklist_items,omitempty"`
	Reminders      []ReminderResponse       `json:"reminders,omitempty"`
	CreatedAt      time.Time                `json:"created_at"`
	UpdatedAt      time.Time                `json:"updated_at"`
	CreatedBy      string                   `json:"created_by"`
	UpdatedBy      string                   `json:"updated_by"`
}

type NoteAttachmentResponse struct {
	FileName    string `json:"file_name"`
	URL         string `json:"url"`
	ContentType string `json:"content_type"`
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

type ReminderResponse struct {
	Time time.Time `json:"time"`
}

type CreateContactRequest struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
}

type CreateContactResponse struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateContactRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
}
type UpdateContactResponse struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
