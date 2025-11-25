package repository

import (
	"database/sql"
	"project-uas/app/model"
	// "time" // Dihapus
	"github.com/google/uuid"
)

func GetAllPermissions(db *sql.DB) ([]model.Permission, error) {
	// Query diubah: +resource, +action, -created_at
	rows, err := db.Query(`
        SELECT id, name, resource, action, description 
        FROM permissions ORDER BY resource, action, name ASC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []model.Permission
	for rows.Next() {
		var p model.Permission
		// Scan diubah: +resource, +action, -created_at
		if err := rows.Scan(&p.ID, &p.Name, &p.Resource, &p.Action, &p.Description); err != nil {
			return nil, err
		}
		permissions = append(permissions, p)
	}
	return permissions, nil
}

func GetPermissionByID(db *sql.DB, id uuid.UUID) (*model.Permission, error) {
	var p model.Permission
	// Query diubah: +resource, +action, -created_at
	row := db.QueryRow(`
        SELECT id, name, resource, action, description 
        FROM permissions WHERE id = $1
    `, id)
	// Scan diubah: +resource, +action, -created_at
	err := row.Scan(&p.ID, &p.Name, &p.Resource, &p.Action, &p.Description)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func CreatePermission(db *sql.DB, p *model.Permission) error {
	p.ID = uuid.New()
	// p.CreatedAt = time.Now() // <-- DIHAPUS

	// Query diubah: 5 kolom, -created_at
	_, err := db.Exec(`
        INSERT INTO permissions (id, name, resource, action, description)
        VALUES ($1, $2, $3, $4, $5)
    `, p.ID, p.Name, p.Resource, p.Action, p.Description)
	return err
}

func UpdatePermission(db *sql.DB, p *model.Permission) error {
	// Query diubah: menyertakan resource dan action
	_, err := db.Exec(`
        UPDATE permissions SET name = $1, resource = $2, action = $3, description = $4
        WHERE id = $5
    `, p.Name, p.Resource, p.Action, p.Description, p.ID)
	return err
}

func DeletePermission(db *sql.DB, id uuid.UUID) error {
	_, err := db.Exec("DELETE FROM permissions WHERE id = $1", id)
	return err
}