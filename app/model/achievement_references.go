package model

import (
	"time"

	"github.com/google/uuid"
)

// Enum status (HARUS sama dengan PostgreSQL)
type AchievementStatus string

const (
	StatusDraft     AchievementStatus = "draft"
	StatusSubmitted AchievementStatus = "submitted"
	StatusVerified  AchievementStatus = "verified"
	StatusRejected  AchievementStatus = "rejected"
)

// AchievementReference mencerminkan tabel PostgreSQL
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

// --- STRUCT REQUEST SUBMIT ---
type SubmitAchievementRequest struct {
	StudentID       uuid.UUID              `json:"student_id"`
	AchievementType string                 `json:"achievement_type"`
	Title           string                 `json:"title"`
	Description     string                 `json:"description"`
	Details         map[string]interface{} `json:"details"`
	Tags            []string               `json:"tags"`
	Points          int                    `json:"points"`
}

// --- STRUCT REQUEST UPDATE STATUS ---
type UpdateAchievementStatusRequest struct {
	Status        AchievementStatus `json:"status"`
	VerifiedBy    uuid.UUID         `json:"verified_by"`
	RejectionNote *string           `json:"rejection_note"`
}
