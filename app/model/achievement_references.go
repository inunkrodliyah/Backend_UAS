package model

import (
	"time"

	"github.com/google/uuid"
)

// Definisikan tipe data custom untuk status (sesuai ENUM 'achievement_status' di DB)
type AchievementStatus string

const (
	StatusPending  AchievementStatus = "pending"
	StatusApproved AchievementStatus = "approved"
	StatusRejected AchievementStatus = "rejected"
)

// AchievementReference struct ini sekarang mencerminkan tabel 'achievement_references' Anda
type AchievementReference struct {
	ID                   uuid.UUID          `json:"id"`
	StudentID            uuid.UUID          `json:"student_id"`
	MongoAchievementID   string             `json:"mongo_achievement_id"`
	Status               AchievementStatus  `json:"status"`
	SubmittedAt          *time.Time         `json:"submitted_at"`     // Pointer untuk handle NULL
	VerifiedAt           *time.Time         `json:"verified_at"`      // Pointer untuk handle NULL
	VerifiedBy           *uuid.UUID         `json:"verified_by"`      // Pointer untuk handle NULL
	RejectionNote        *string            `json:"rejection_note"`   // Pointer untuk handle NULL
	CreatedAt            time.Time          `json:"created_at"`
	UpdatedAt            time.Time          `json:"updated_at"`
}

// --- STRUCT REQUEST BARU UNTUK CREATE (SUBMIT) ---
// Ini adalah data yang dikirim Mahasiswa saat submit prestasi
type CreateAchievementReferenceRequest struct {
	StudentID          uuid.UUID `json:"student_id"`
	MongoAchievementID string    `json:"mongo_achievement_id"`
}

// --- STRUCT REQUEST BARU UNTUK UPDATE (VERIFY/REJECT) ---
// Ini adalah data yang dikirim Dosen/Admin saat memverifikasi
type UpdateAchievementStatusRequest struct {
	Status        AchievementStatus `json:"status"`         // "approved" or "rejected"
	VerifiedBy    uuid.UUID         `json:"verified_by"`    // ID Dosen/Admin yang memverifikasi
	RejectionNote *string           `json:"rejection_note"` // Opsional, wajib diisi jika status "rejected"
}