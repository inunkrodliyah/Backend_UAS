package model

import (
	"time"

	"github.com/google/uuid"
)

type Student struct {
	ID            uuid.UUID  `json:"id"`
	UserID        uuid.UUID  `json:"user_id"`
	StudentID     string     `json:"student_id"`     
	ProgramStudy  string     `json:"program_study"`  
	AcademicYear  string     `json:"academic_year"`  
	AdvisorID     *uuid.UUID `json:"advisor_id"`     
	CreatedAt     time.Time  `json:"created_at"`
}

// --- STRUCT REQUEST UNTUK CREATE ---
type CreateStudentRequest struct {
	UserID        uuid.UUID  `json:"user_id"`
	StudentID     string     `json:"student_id"`
	ProgramStudy  string     `json:"program_study"`
	AcademicYear  string     `json:"academic_year"`
	AdvisorID     *uuid.UUID `json:"advisor_id"` // Opsional
}

// --- STRUCT REQUEST UNTUK UPDATE ---
type UpdateStudentRequest struct {
	StudentID     string     `json:"student_id"`
	ProgramStudy  string     `json:"program_study"`
	AcademicYear  string     `json:"academic_year"`
	AdvisorID     *uuid.UUID `json:"advisor_id"` 
}

// Achievement untuk respons API prestasi mahasiswa
type StudentAchievement struct {
	ID          uuid.UUID `json:"id"`
	StudentID   uuid.UUID `json:"student_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// Request khusus update advisor
type UpdateAdvisorRequest struct {
	AdvisorID *uuid.UUID `json:"advisor_id"`
}
