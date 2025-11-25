package model

import (
	// "time" // Dihapus karena tidak ada 'created_at' di tabel
	"github.com/google/uuid"
)

// Permission struct ini sekarang mencerminkan tabel 'permissions' Anda
type Permission struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Resource    string    `json:"resource"`    // <-- DITAMBAHKAN (sesuai tabel)
	Action      string    `json:"action"`      // <-- DITAMBAHKAN (sesuai tabel)
	Description *string   `json:"description"`
	// CreatedAt   time.Time `json:"created_at"` // <-- DIHAPUS (tidak ada di tabel)
}

// Struct request Create disesuaikan
type CreatePermissionRequest struct {
	Name        string  `json:"name"`
	Resource    string  `json:"resource"`    // <-- DITAMBAHKAN
	Action      string  `json:"action"`      // <-- DITAMBAHKAN
	Description *string `json:"description"`
}

// Struct request Update disesuaikan
type UpdatePermissionRequest struct {
	Name        string  `json:"name"`
	Resource    string  `json:"resource"`    // <-- DITAMBAHKAN
	Action      string  `json:"action"`      // <-- DITAMBAHKAN
	Description *string `json:"description"`
}