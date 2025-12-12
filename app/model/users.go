package model

import (
	"time"

	"github.com/google/uuid"
)

// Entity Database
type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Hidden
	FullName     string    `json:"full_name"`
	RoleID       uuid.UUID `json:"role_id"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Request Create User
type CreateUserRequest struct {
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	FullName string    `json:"full_name"`
	RoleID   uuid.UUID `json:"role_id"`
}

// Request Update User
type UpdateUserRequest struct {
	Username string    `json:"username"`
	Email    string    `json:"email"`
	FullName string    `json:"full_name"`
	RoleID   uuid.UUID `json:"role_id"`
	IsActive *bool     `json:"is_active"`
}

// Request Update Role Only
type UpdateUserRoleRequest struct {
	RoleID uuid.UUID `json:"role_id"`
}