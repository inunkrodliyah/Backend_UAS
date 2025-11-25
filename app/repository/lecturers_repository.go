package repository

import (
	"database/sql"
	"project-uas/app/model"
	"time"

	"github.com/google/uuid"
)

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