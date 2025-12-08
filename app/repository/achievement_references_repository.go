package repository

import (
	"database/sql"
	"project-uas/app/model"
	"time"

	"github.com/google/uuid"
)

func GetAllAchievementReferences(db *sql.DB) ([]model.AchievementReference, error) {
	rows, err := db.Query(`
        SELECT id, student_id, mongo_achievement_id, status, submitted_at, 
               verified_at, verified_by, rejection_note, created_at, updated_at
        FROM achievement_references ORDER BY created_at DESC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refs []model.AchievementReference
	for rows.Next() {
		var r model.AchievementReference
		if err := rows.Scan(
			&r.ID, &r.StudentID, &r.MongoAchievementID, &r.Status, &r.SubmittedAt,
			&r.VerifiedAt, &r.VerifiedBy, &r.RejectionNote, &r.CreatedAt, &r.UpdatedAt,
		); err != nil {
			return nil, err
		}
		refs = append(refs, r)
	}
	return refs, nil
}

func GetAchievementReferenceByID(db *sql.DB, id uuid.UUID) (*model.AchievementReference, error) {
	var r model.AchievementReference

	row := db.QueryRow(`
        SELECT id, student_id, mongo_achievement_id, status, submitted_at, 
               verified_at, verified_by, rejection_note, created_at, updated_at
        FROM achievement_references WHERE id = $1
    `, id)

	err := row.Scan(
		&r.ID, &r.StudentID, &r.MongoAchievementID, &r.Status, &r.SubmittedAt,
		&r.VerifiedAt, &r.VerifiedBy, &r.RejectionNote, &r.CreatedAt, &r.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func CreateAchievementReference(db *sql.DB, r *model.AchievementReference) error {
	r.ID = uuid.New()
	now := time.Now()
	if r.CreatedAt.IsZero() {
		r.CreatedAt = now
	}
	r.UpdatedAt = now

	_, err := db.Exec(`
        INSERT INTO achievement_references 
        (id, student_id, mongo_achievement_id, status, submitted_at, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `, r.ID, r.StudentID, r.MongoAchievementID, r.Status, r.SubmittedAt, r.CreatedAt, r.UpdatedAt)

	return err
}

func UpdateAchievementReference(db *sql.DB, r *model.AchievementReference) error {
	r.UpdatedAt = time.Now()

	_, err := db.Exec(`
        UPDATE achievement_references 
        SET status = $1, verified_at = $2, verified_by = $3, rejection_note = $4, updated_at = $5
        WHERE id = $6
    `, r.Status, r.VerifiedAt, r.VerifiedBy, r.RejectionNote, r.UpdatedAt, r.ID)

	return err
}

func DeleteAchievementReference(db *sql.DB, id uuid.UUID) error {
	_, err := db.Exec("DELETE FROM achievement_references WHERE id = $1", id)
	return err
}