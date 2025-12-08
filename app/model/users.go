package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` 
	FullName     string    `json:"full_name"`
	RoleID       uuid.UUID `json:"role_id"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	FullName string    `json:"full_name"`
	RoleID   uuid.UUID `json:"role_id"`
}

type UpdateUserRequest struct {
	Username string    `json:"username"`
	Email    string    `json:"email"`
	FullName string    `json:"full_name"`
	RoleID   uuid.UUID `json:"role_id"`
	IsActive *bool     `json:"is_active"` 
}

// LoginRequest: Input dari user
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserLoginData: Objek 'user' di dalam 'data'
type UserLoginData struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	FullName string    `json:"fullName"` // SRS biasanya pakai camelCase
	RoleID   uuid.UUID `json:"roleId"`   // Sementara pakai RoleID dulu
	// Permissions []string `json:"permissions"` // Opsional: Tambahkan jika sudah ada logic permission
}

// LoginResponseData: Objek 'data'
type LoginResponseData struct {
	Token        string        `json:"token"`
	RefreshToken string        `json:"refreshToken"` // SRS minta refreshToken
	User         UserLoginData `json:"user"`         // Nested Object User
}

// LoginResponse: Wrapper Utama (Root)
type LoginResponse struct {
	Status string            `json:"status"` // "success" atau "fail"
	Data   LoginResponseData `json:"data"`
}