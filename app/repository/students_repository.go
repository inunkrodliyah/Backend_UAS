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

// GetStudentByID (Berdasarkan UserID sesuai endpoint param)
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

// CreateStudent
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
func UpdateStudent(db *sql.DB, s *model.Student) error {
	_, err := db.Exec(`
        UPDATE students SET student_id = $1, program_study = $2, academic_year = $3, advisor_id = $4
        WHERE user_id = $5
    `, s.StudentID, s.ProgramStudy, s.AcademicYear, s.AdvisorID, s.UserID)

	return err
}

// UpdateAdvisor Only
func UpdateAdvisor(db *sql.DB, studentID uuid.UUID, advisorID *uuid.UUID) error {
	// studentID disini adalah ID Primary Key dari tabel students (bukan UserID)
	// Kita sesuaikan querynya agar aman
	_, err := db.Exec(`
        UPDATE students
        SET advisor_id = $1
        WHERE user_id = $2 
    `, advisorID, studentID) // Asumsi ID di URL adalah UserID, jika ID tabel student, ubah query WHERE id=$2
	return err
}

// GetAchievementsByStudentID: Mengambil list referensi prestasi
// FIX: Menggunakan tabel achievement_references, bukan achievements
func GetAchievementsByStudentID(db *sql.DB, studentID uuid.UUID) ([]model.AchievementReference, error) {
	rows, err := db.Query(`
        SELECT id, student_id, mongo_achievement_id, status, submitted_at, 
               verified_at, verified_by, rejection_note, created_at, updated_at
        FROM achievement_references
        WHERE student_id = $1
        ORDER BY created_at DESC
    `, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.AchievementReference
	for rows.Next() {
		var r model.AchievementReference
		if err := rows.Scan(
			&r.ID, &r.StudentID, &r.MongoAchievementID, &r.Status, &r.SubmittedAt,
			&r.VerifiedAt, &r.VerifiedBy, &r.RejectionNote, &r.CreatedAt, &r.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, r)
	}
	return list, nil
}