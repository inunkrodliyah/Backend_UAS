package model

import "github.com/google/uuid"

// Request untuk Login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Request untuk Refresh Token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" example:"a3f5d9c2-8e34-4f90-9d21-1cfa9f4b8e71"`
}


// --- STRUCTURE RESPONSE (SESUAI FR-001 & SRS) ---

// UserLoginData: Detail user dalam response login
type UserLoginData struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	FullName    string    `json:"fullName"`
	RoleID      uuid.UUID `json:"roleId"`
	Permissions []string  `json:"permissions"` // Sesuai FR-001
}

// LoginResponseData: Isi dari field 'data'
type LoginResponseData struct {
	Token        string        `json:"token"`
	RefreshToken string        `json:"refreshToken"`
	User         UserLoginData `json:"user"`
}

// AuthResponse: Wrapper utama
type AuthResponse struct {
	Status string      `json:"status"` // "success"
	Data   interface{} `json:"data"`   // Bisa LoginResponseData atau User Profile
}