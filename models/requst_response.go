package models

import "time"

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Age      int    `json:"age"`
}

// request model
type CreateNoteRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

// response model
type NoteResponse struct {
	ID          uint      `json:"id"`
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
