package model

import (
	"time"

	"github.com/google/uuid"
)

// Student struct ini sekarang mencerminkan tabel 'students' Anda
type Student struct {
	ID            uuid.UUID  `json:"id"`
	UserID        uuid.UUID  `json:"user_id"`
	StudentID     string     `json:"student_id"`     // Sebelumnya 'NIM'
	ProgramStudy  string     `json:"program_study"`  // Sebelumnya 'Major'
	AcademicYear  string     `json:"academic_year"`  // Sebelumnya 'EntryYear' (tipe data diubah ke string)
	AdvisorID     *uuid.UUID `json:"advisor_id"`     // Dosen Wali (pointer untuk handle NULL)
	CreatedAt     time.Time  `json:"created_at"`
}

// --- STRUCT REQUEST BARU UNTUK CREATE ---
// Ini adalah data yang dikirim client untuk menautkan User yang ADA ke data Student BARU.
type CreateStudentRequest struct {
	UserID        uuid.UUID  `json:"user_id"`
	StudentID     string     `json:"student_id"`
	ProgramStudy  string     `json:"program_study"`
	AcademicYear  string     `json:"academic_year"`
	AdvisorID     *uuid.UUID `json:"advisor_id"` // Opsional
}

// --- STRUCT REQUEST BARU UNTUK UPDATE ---
type UpdateStudentRequest struct {
	StudentID     string     `json:"student_id"`
	ProgramStudy  string     `json:"program_study"`
	AcademicYear  string     `json:"academic_year"`
	AdvisorID     *uuid.UUID `json:"advisor_id"` // Opsional
}

// Struct 'CreateStudentUserRequest' yang lama (gabungan) sudah dihapus.