package model

import (
	"time"

	"github.com/google/uuid"
)

// Struct ini sudah SAMA DENGAN TABEL ANDA (Benar)
type Lecturer struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	LecturerID string    `json:"lecturer_id"`
	Department string    `json:"department"`
	CreatedAt  time.Time `json:"created_at"`
}

// --- STRUCT REQUEST BARU UNTUK CREATE ---
// Ini adalah data yang dikirim client (cth: Postman)
// untuk menautkan User yang ADA ke data Dosen BARU.
type CreateLecturerRequest struct {
	UserID     uuid.UUID `json:"user_id"`
	LecturerID string    `json:"lecturer_id"`
	Department string    `json:"department"`
}

// Struct Update (Tetap sama, sudah benar)
type UpdateLecturerRequest struct {
	LecturerID string `json:"lecturer_id"`
	Department string `json:"department"`
}

// Struct 'CreateLecturerUserRequest' yang lama (gabungan) sudah dihapus
// karena tidak sesuai dengan alur Anda.