package repository

import (
	"database/sql"
	"project-uas/app/model"
	"time"

	"github.com/google/uuid"
)

// Ambil semua (Untuk Admin)
func GetAllAchievementReferences(db *sql.DB) ([]model.AchievementReference, error) {
	// TAMBAHAN: Filter status != 'deleted' dan select deleted_at
	return fetchAchievements(db, `
		SELECT id, student_id, mongo_achievement_id, status, submitted_at, 
		       verified_at, verified_by, rejection_note, created_at, updated_at, deleted_at
		FROM achievement_references 
		WHERE status != 'deleted' 
		ORDER BY created_at DESC
	`)
}

// Ambil berdasarkan Student ID
func GetAchievementReferencesByStudentID(db *sql.DB, studentID uuid.UUID) ([]model.AchievementReference, error) {
	return fetchAchievements(db, `
		SELECT id, student_id, mongo_achievement_id, status, submitted_at, 
		       verified_at, verified_by, rejection_note, created_at, updated_at, deleted_at
		FROM achievement_references 
		WHERE student_id = $1 AND status != 'deleted' 
		ORDER BY created_at DESC
	`, studentID)
}

// Ambil berdasarkan Advisor ID
func GetAchievementReferencesByAdvisorID(db *sql.DB, advisorID uuid.UUID) ([]model.AchievementReference, error) {
	return fetchAchievements(db, `
		SELECT ar.id, ar.student_id, ar.mongo_achievement_id, ar.status, ar.submitted_at, 
		       ar.verified_at, ar.verified_by, ar.rejection_note, ar.created_at, ar.updated_at, ar.deleted_at
		FROM achievement_references ar
		JOIN students s ON ar.student_id = s.id
		WHERE s.advisor_id = $1 AND ar.status != 'draft' AND ar.status != 'deleted'
		ORDER BY ar.created_at DESC
	`, advisorID)
}

// Helper function untuk scan rows
func fetchAchievements(db *sql.DB, query string, args ...interface{}) ([]model.AchievementReference, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refs []model.AchievementReference
	for rows.Next() {
		var r model.AchievementReference
		// TAMBAHAN: Scan juga r.DeletedAt
		if err := rows.Scan(&r.ID, &r.StudentID, &r.MongoAchievementID, &r.Status, &r.SubmittedAt, &r.VerifiedAt, &r.VerifiedBy, &r.RejectionNote, &r.CreatedAt, &r.UpdatedAt, &r.DeletedAt); err != nil {
			return nil, err
		}
		refs = append(refs, r)
	}
	return refs, nil
}

// Get By ID (Detail)
func GetAchievementReferenceByID(db *sql.DB, id uuid.UUID) (*model.AchievementReference, error) {
	var r model.AchievementReference
	// TAMBAHAN: Select deleted_at
	row := db.QueryRow(`
		SELECT id, student_id, mongo_achievement_id, status, submitted_at, 
		       verified_at, verified_by, rejection_note, created_at, updated_at, deleted_at
		FROM achievement_references WHERE id = $1
	`, id)
	// TAMBAHAN: Scan deleted_at
	err := row.Scan(&r.ID, &r.StudentID, &r.MongoAchievementID, &r.Status, &r.SubmittedAt, &r.VerifiedAt, &r.VerifiedBy, &r.RejectionNote, &r.CreatedAt, &r.UpdatedAt, &r.DeletedAt)
	if err != nil {
		return nil, err
	}
	// Opsional: Jika ingin return error kalau ternyata deleted
	if r.Status == model.StatusDeleted {
		return nil, sql.ErrNoRows 
	}
	return &r, nil
}

func CreateAchievementReference(db *sql.DB, r *model.AchievementReference) error {
	r.ID = uuid.New()
	now := time.Now()
	r.CreatedAt = now
	r.UpdatedAt = now
	_, err := db.Exec(`
		INSERT INTO achievement_references (id, student_id, mongo_achievement_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, r.ID, r.StudentID, r.MongoAchievementID, r.Status, r.CreatedAt, r.UpdatedAt)
	return err
}

// Update Status Generic
func UpdateAchievementStatus(db *sql.DB, r *model.AchievementReference) error {
	r.UpdatedAt = time.Now()
	_, err := db.Exec(`
		UPDATE achievement_references 
		SET status = $1, submitted_at = $2, verified_at = $3, verified_by = $4, rejection_note = $5, updated_at = $6
		WHERE id = $7
	`, r.Status, r.SubmittedAt, r.VerifiedAt, r.VerifiedBy, r.RejectionNote, r.UpdatedAt, r.ID)
	return err
}

func DeleteAchievementReference(db *sql.DB, id uuid.UUID) error {
	// MENGUBAH QUERY DELETE MENJADI UPDATE
	_, err := db.Exec(`
		UPDATE achievement_references 
		SET status = 'deleted', deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`, id)
	return err
}