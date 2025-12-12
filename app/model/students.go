package model

import (
	"time"

	"github.com/google/uuid"
)

type Student struct {
	ID           uuid.UUID  `json:"id"`
	UserID       uuid.UUID  `json:"user_id"`
	StudentID    string     `json:"student_id"`
	ProgramStudy string     `json:"program_study"`
	AcademicYear string     `json:"academic_year"`
	AdvisorID    *uuid.UUID `json:"advisor_id"`
	CreatedAt    time.Time  `json:"created_at"`
}

// Struct Request Create
type CreateStudentRequest struct {
	UserID       uuid.UUID  `json:"user_id"`
	StudentID    string     `json:"student_id"`
	ProgramStudy string     `json:"program_study"`
	AcademicYear string     `json:"academic_year"`
	AdvisorID    *uuid.UUID `json:"advisor_id"`
}

// Struct Request Update
type UpdateStudentRequest struct {
	StudentID    string     `json:"student_id"`
	ProgramStudy string     `json:"program_study"`
	AcademicYear string     `json:"academic_year"`
	AdvisorID    *uuid.UUID `json:"advisor_id"`
}

// Request khusus update advisor (Endpoint: PUT /students/:id/advisor)
type UpdateAdvisorRequest struct {
	AdvisorID *uuid.UUID `json:"advisor_id"`
}