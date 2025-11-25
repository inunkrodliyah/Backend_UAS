package repository

import (
	"database/sql"
	"project-uas/app/model"
	"time"

	"github.com/google/uuid"
)

func GetAllRoles(db *sql.DB) ([]model.Role, error) {
	rows, err := db.Query(`
        SELECT id, name, description, created_at FROM roles ORDER BY name ASC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []model.Role
	for rows.Next() {
		var r model.Role
		if err := rows.Scan(&r.ID, &r.Name, &r.Description, &r.CreatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, r)
	}
	return roles, nil
}

func GetRoleByID(db *sql.DB, id uuid.UUID) (*model.Role, error) {
	var r model.Role
	row := db.QueryRow(`
        SELECT id, name, description, created_at FROM roles WHERE id = $1
    `, id)
	err := row.Scan(&r.ID, &r.Name, &r.Description, &r.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func CreateRole(db *sql.DB, r *model.Role) error {
	r.ID = uuid.New()
	r.CreatedAt = time.Now()
	_, err := db.Exec(`
        INSERT INTO roles (id, name, description, created_at)
        VALUES ($1, $2, $3, $4)
    `, r.ID, r.Name, r.Description, r.CreatedAt)
	return err
}

func UpdateRole(db *sql.DB, r *model.Role) error {
	_, err := db.Exec(`
        UPDATE roles SET name = $1, description = $2
        WHERE id = $3
    `, r.Name, r.Description, r.ID)
	return err
}

func DeleteRole(db *sql.DB, id uuid.UUID) error {
	_, err := db.Exec("DELETE FROM roles WHERE id = $1", id)
	return err
}