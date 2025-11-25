package model

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"` 
	CreatedAt   time.Time `json:"created_at"`
}

// Struct untuk request body saat Create
type CreateRoleRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

// Struct untuk request body saat Update
type UpdateRoleRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}