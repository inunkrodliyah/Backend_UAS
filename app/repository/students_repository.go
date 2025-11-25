package repository

import (
	"database/sql"
	"project-uas/app/model"
	"time"

	"github.com/google/uuid"
)

// GetStudentByID (Sebenarnya GetByUserID)
// Mengambil data student berdasarkan user_id (FK)
func GetStudentByID(db *sql.DB, userID uuid.UUID) (*model.Student, error) {
	var s model.Student
	row := db.QueryRow(`
        SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at 
        FROM students WHERE user_id = $1
    `, userID)
	
	// Scan disesuaikan dengan 7 kolom
	err := row.Scan(
		&s.ID, &s.UserID, &s.StudentID, &s.ProgramStudy,
		&s.AcademicYear, &s.AdvisorID, &s.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// CreateStudent
// HANYA membuat data di tabel 'students'
func CreateStudent(db *sql.DB, s *model.Student) error {
	s.ID = uuid.New()
	s.CreatedAt = time.Now()
	
	_, err := db.Exec(`
        INSERT INTO students (id, user_id, student_id, program_study, academic_year, advisor_id, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `, s.ID, s.UserID, s.StudentID, s.ProgramStudy, s.AcademicYear, s.AdvisorID, s.CreatedAt)
	
	return err
}

// UpdateStudent
// HANYA mengupdate data di tabel 'students'
func UpdateStudent(db *sql.DB, s *model.Student) error {
	_, err := db.Exec(`
        UPDATE students SET student_id = $1, program_study = $2, academic_year = $3, advisor_id = $4
        WHERE user_id = $5
    `, s.StudentID, s.ProgramStudy, s.AcademicYear, s.AdvisorID, s.UserID)
	
	return err
}

// DeleteStudent (Tidak diperlukan jika mengandalkan ON DELETE CASCADE dari tabel users)