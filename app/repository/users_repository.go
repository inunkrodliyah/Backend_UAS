package repository

import (
	"database/sql"
	"project-uas/app/model"
	"time"

	"github.com/google/uuid"
)

// Catatan: Transaksi (Tx) sangat direkomendasikan di sini,
// terutama saat membuat User + Student/Lecturer.
// Namun, untuk mengikuti pola, saya akan membuatnya terpisah.

func GetAllUsers(db *sql.DB) ([]model.User, error) {
	rows, err := db.Query(`
        SELECT id, username, email, full_name, role_id, is_active, created_at, updated_at
        FROM users ORDER BY full_name ASC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.FullName, &u.RoleID, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func GetUserByID(db *sql.DB, id uuid.UUID) (*model.User, error) {
	var u model.User
	row := db.QueryRow(`
        SELECT id, username, email, full_name, role_id, is_active, created_at, updated_at
        FROM users WHERE id = $1
    `, id)
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.FullName, &u.RoleID, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func CreateUser(db *sql.DB, u *model.User) error {
	u.ID = uuid.New()
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
	u.IsActive = true // Default

	_, err := db.Exec(`
        INSERT INTO users (id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, u.ID, u.Username, u.Email, u.PasswordHash, u.FullName, u.RoleID, u.IsActive, u.CreatedAt, u.UpdatedAt)
	return err
}

func UpdateUser(db *sql.DB, u *model.User) error {
	u.UpdatedAt = time.Now()
	_, err := db.Exec(`
        UPDATE users SET username = $1, email = $2, full_name = $3, role_id = $4, is_active = $5, updated_at = $6
        WHERE id = $7
    `, u.Username, u.Email, u.FullName, u.RoleID, u.IsActive, u.UpdatedAt, u.ID)
	return err
}

func DeleteUser(db *sql.DB, id uuid.UUID) error {
	// Menghapus user akan cascade ke student/lecturer jika FK di-setting ON DELETE CASCADE
	_, err := db.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}