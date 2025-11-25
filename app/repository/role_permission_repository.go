package repository

import (
	"database/sql"
	"project-uas/app/model"

	"github.com/google/uuid"
)

// GetPermissionsByRoleID mengambil data Permission (lengkap) berdasarkan Role
func GetPermissionsByRoleID(db *sql.DB, roleID uuid.UUID) ([]model.Permission, error) {
	
	// --- PERBAIKAN DI SINI ---
	// Query diubah: -created_at, +resource, +action
	rows, err := db.Query(`
        SELECT p.id, p.name, p.resource, p.action, p.description
        FROM permissions p
        JOIN role_permissions rp ON p.id = rp.permission_id
        WHERE rp.role_id = $1
    `, roleID)
	// --- AKHIR PERBAIKAN ---

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []model.Permission
	for rows.Next() {
		var p model.Permission // Ini adalah model.Permission, BUKAN model.RolePermission

		// --- PERBAIKAN DI SINI ---
		// Scan disesuaikan: -created_at, +resource, +action
		if err := rows.Scan(&p.ID, &p.Name, &p.Resource, &p.Action, &p.Description); err != nil {
			return nil, err
		}
		// --- AKHIR PERBAIKAN ---

		permissions = append(permissions, p)
	}
	return permissions, nil
}

// Fungsi AssignPermissionToRole sudah benar
func AssignPermissionToRole(db *sql.DB, rp *model.RolePermission) error {
	_, err := db.Exec(`
        INSERT INTO role_permissions (role_id, permission_id)
        VALUES ($1, $2)
        ON CONFLICT (role_id, permission_id) DO NOTHING
    `, rp.RoleID, rp.PermissionID)
	return err
}

// Fungsi RevokePermissionFromRole sudah benar
func RevokePermissionFromRole(db *sql.DB, roleID, permissionID uuid.UUID) error {
	_, err := db.Exec(`
        DELETE FROM role_permissions
        WHERE role_id = $1 AND permission_id = $2
    `, roleID, permissionID)
	return err
}