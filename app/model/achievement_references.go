package model

import (
	"time"

	"github.com/google/uuid"
)

// Enum status
type AchievementStatus string

const (
	StatusDraft     AchievementStatus = "draft"
	StatusSubmitted AchievementStatus = "submitted"
	StatusVerified  AchievementStatus = "verified"
	StatusRejected  AchievementStatus = "rejected"
)

// AchievementReference (Sesuai Tabel PostgreSQL)
type AchievementReference struct {
	ID                 uuid.UUID         `json:"id"`
	StudentID          uuid.UUID         `json:"student_id"`
	MongoAchievementID string            `json:"mongo_achievement_id"`
	Status             AchievementStatus `json:"status"`
	SubmittedAt        *time.Time        `json:"submitted_at"`
	VerifiedAt         *time.Time        `json:"verified_at"`
	VerifiedBy         *uuid.UUID        `json:"verified_by"`
	RejectionNote      *string           `json:"rejection_note"`
	CreatedAt          time.Time         `json:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at"`
}

// Request: Create / Submit Awal (Draft)
type CreateAchievementRequest struct {
	StudentID       uuid.UUID              `json:"student_id"`
	AchievementType string                 `json:"achievement_type"`
	Title           string                 `json:"title"`
	Description     string                 `json:"description"`
	Details         map[string]interface{} `json:"details"`
	Tags            []string               `json:"tags"`
	Points          int                    `json:"points"`
}

// Request: Update Prestasi (Hanya bisa saat Draft)
type UpdateAchievementRequest struct {
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Details     map[string]interface{} `json:"details"`
	Tags        []string               `json:"tags"`
	Points      int                    `json:"points"`
}

// Request: Reject
type RejectAchievementRequest struct {
	RejectionNote string `json:"rejection_note"`
}

// Response: History
type AchievementHistoryResponse struct {
	Status    AchievementStatus `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Actor     string            `json:"actor"` // "Mahasiswa" atau "Dosen"
	Note      *string           `json:"note,omitempty"`
}