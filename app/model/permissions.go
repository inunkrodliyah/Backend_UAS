package model

import (

	"github.com/google/uuid"
)

type Permission struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Resource    string    `json:"resource"`    
	Action      string    `json:"action"`      
	Description *string   `json:"description"`
}

// Struct request Create 
type CreatePermissionRequest struct {
	Name        string  `json:"name"`
	Resource    string  `json:"resource"`    
	Action      string  `json:"action"`      
	Description *string `json:"description"`
}

// Struct request Update 
type UpdatePermissionRequest struct {
	Name        string  `json:"name"`
	Resource    string  `json:"resource"`    
	Action      string  `json:"action"`      
	Description *string `json:"description"`
}