package models

import (
	"time"

	"github.com/google/uuid"
)

type CreateUserRequest struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	ClerkID   string `json:"clerkId"`
	Email     string `json:"email" binding:"required,email"`
	ImageUrl  string `json:"image_url"`
	Password  string `json:"password"`
}

// request model
type CreateNoteRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

// response model
type NoteResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedBy   uint      `json:"created_by"`
	UpdatedBy   uint      `json:"updated_by"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
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
