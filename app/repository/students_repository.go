package repository

import (
	"database/sql"
	"project-uas/app/model"
	"time"

	"github.com/google/uuid"
)

func GetAllStudents(db *sql.DB) ([]model.Student, error) {
	rows, err := db.Query(`
        SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at 
        FROM students ORDER BY created_at DESC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []model.Student
	for rows.Next() {
		var s model.Student
		if err := rows.Scan(
			&s.ID, &s.UserID, &s.StudentID, &s.ProgramStudy,
			&s.AcademicYear, &s.AdvisorID, &s.CreatedAt,
		); err != nil {
			return nil, err
		}
		students = append(students, s)
	}

	return students, nil
}

// GetStudentByID (Sebenarnya GetByUserID)
// Mengambil data student berdasarkan user_id (FK)
func GetStudentByID(db *sql.DB, userID uuid.UUID) (*model.Student, error) {
	var s model.Student
	row := db.QueryRow(`
        SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at 
        FROM students WHERE user_id = $1
    `, userID)

	err := row.Scan(
		&s.ID, &s.UserID, &s.StudentID, &s.ProgramStudy,
		&s.AcademicYear, &s.AdvisorID, &s.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func GetAchievementsByStudentID(db *sql.DB, studentID uuid.UUID) ([]model.StudentAchievement, error) {
	rows, err := db.Query(`
        SELECT id, student_id, title, description, created_at
        FROM achievements
        WHERE student_id = $1
        ORDER BY created_at DESC
    `, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.StudentAchievement
	for rows.Next() {
		var a model.StudentAchievement
		if err := rows.Scan(
			&a.ID, &a.StudentID, &a.Title,
			&a.Description, &a.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, nil
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

func UpdateAdvisor(db *sql.DB, studentID uuid.UUID, advisorID *uuid.UUID) error {
	_, err := db.Exec(`
        UPDATE students
        SET advisor_id = $1
        WHERE id = $2
    `, advisorID, studentID)
	return err
}
