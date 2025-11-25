package service

import (
	"log"
	"project-uas/app/model"
	"project-uas/app/repository"
	"project-uas/database"
	"strings" // <-- 1. IMPORT PACKAGE 'strings'

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetLecturerByUserID
// (Fungsi ini sudah benar, tidak ada perubahan)
func GetLecturerByUserID(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "ID [user_id] tidak valid"})
	}

	user, err := repository.GetUserByID(database.DB, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "User tidak ditemukan"})
	}

	lecturer, err := repository.GetLecturerByID(database.DB, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "Data lecturer tidak ditemukan untuk user ini"})
	}

	response := fiber.Map{
		"user_id":        user.ID,
		"username":       user.Username,
		"email":          user.Email,
		"full_name":      user.FullName,
		"role_id":        user.RoleID,
		"is_active":      user.IsActive,
		"lecturer_id":    lecturer.LecturerID,
		"department":     lecturer.Department,
		"lecturer_db_id": lecturer.ID,
	}

	return c.JSON(fiber.Map{"success": true, "data": response})
}

// CreateLecturer
// (Fungsi ini diperbaiki)
func CreateLecturer(c *fiber.Ctx) error {
	var req model.CreateLecturerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Request body tidak valid"})
	}

	if req.UserID == uuid.Nil || req.LecturerID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Field wajib: 'user_id' dan 'lecturer_id'",
		})
	}

	lecturer := &model.Lecturer{
		UserID:     req.UserID,
		LecturerID: req.LecturerID,
		Department: req.Department,
	}

	if err := repository.CreateLecturer(database.DB, lecturer); err != nil {
		log.Println("Error membuat lecturer:", err)

		// --- 2. PERBAIKAN DI SINI ---
		// Diubah dari err.Error().Contains(...) menjadi strings.Contains(err.Error(), ...)
		if strings.Contains(err.Error(), "duplicate key") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"success": false, "message": "Gagal: User ini sudah terdaftar sebagai dosen atau ID dosen sudah ada.", "error": err.Error(),
			})
		}
		// --- AKHIR PERBAIKAN ---

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal membuat data lecturer", "error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Data dosen berhasil ditambahkan dan ditautkan ke user",
		"data":    lecturer,
	})
}

// UpdateLecturer
// (Fungsi ini sudah benar, tidak ada perubahan)
func UpdateLecturer(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "ID [user_id] tidak valid"})
	}

	var req model.UpdateLecturerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Request body tidak valid"})
	}

	lecturer, err := repository.GetLecturerByID(database.DB, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "Data lecturer tidak ditemukan"})
	}

	if req.LecturerID != "" {
		lecturer.LecturerID = req.LecturerID
	}
	if req.Department != "" {
		lecturer.Department = req.Department
	}

	if err := repository.UpdateLecturer(database.DB, lecturer); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal mengupdate data lecturer", "error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Data lecturer berhasil diupdate", "data": lecturer})
}