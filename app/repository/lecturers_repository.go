package repository

import (
	"database/sql"
	"project-uas/app/model"
	"time"

	"github.com/google/uuid"
)

func GetAllLecturers(db *sql.DB) ([]model.Lecturer, error) {
	rows, err := db.Query(`
        SELECT id, user_id, lecturer_id, department, created_at
        FROM lecturers
        ORDER BY created_at DESC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Lecturer
	for rows.Next() {
		var l model.Lecturer
		if err := rows.Scan(&l.ID, &l.UserID, &l.LecturerID, &l.Department, &l.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, l)
	}

	return list, nil
}

// GetLecturerByID (Sebenarnya GetByUserID)
// Mengambil data lecturer berdasarkan user_id (FK)
func GetLecturerByID(db *sql.DB, userID uuid.UUID) (*model.Lecturer, error) {
	var l model.Lecturer
	row := db.QueryRow(`
        SELECT id, user_id, lecturer_id, department, created_at 
        FROM lecturers WHERE user_id = $1
    `, userID)
	
	err := row.Scan(&l.ID, &l.UserID, &l.LecturerID, &l.Department, &l.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func GetAdviseesByLecturerID(db *sql.DB, lecturerID uuid.UUID) ([]model.Student, error) {
	rows, err := db.Query(`
        SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
        FROM students
        WHERE advisor_id = $1
        ORDER BY created_at DESC
    `, lecturerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []model.Student
	for rows.Next() {
		var s model.Student
		if err := rows.Scan(
			&s.ID, &s.UserID, &s.StudentID,
			&s.ProgramStudy, &s.AcademicYear,
			&s.AdvisorID, &s.CreatedAt,
		); err != nil {
			return nil, err
		}
		students = append(students, s)
	}

	return students, nil
}

// CreateLecturer
// HANYA membuat data di tabel 'lecturers'
func CreateLecturer(db *sql.DB, l *model.Lecturer) error {
	l.ID = uuid.New()
	l.CreatedAt = time.Now()
	
	_, err := db.Exec(`
        INSERT INTO lecturers (id, user_id, lecturer_id, department, created_at)
        VALUES ($1, $2, $3, $4, $5)
    `, l.ID, l.UserID, l.LecturerID, l.Department, l.CreatedAt)
	
	return err
}

// UpdateLecturer
// HANYA mengupdate data di tabel 'lecturers'
func UpdateLecturer(db *sql.DB, l *model.Lecturer) error {
	_, err := db.Exec(`
        UPDATE lecturers SET lecturer_id = $1, department = $2
        WHERE user_id = $3
    `, l.LecturerID, l.Department, l.UserID)
	
	return err
}