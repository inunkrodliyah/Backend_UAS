package service

import (
	"log"
	"project-uas/app/model"
	"project-uas/app/repository"
	"project-uas/database"
	"strings" 

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetAllLecturers(c *fiber.Ctx) error {
	lecturers, err := repository.GetAllLecturers(database.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data lecturers",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    lecturers,
	})
}

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

func GetLecturerAdvisees(c *fiber.Ctx) error {
	lecturerID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID lecturer tidak valid",
		})
	}

	// Pastikan lecturer ada
	_, err = repository.GetLecturerByID(database.DB, lecturerID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Lecturer tidak ditemukan",
		})
	}

	advisees, err := repository.GetAdviseesByLecturerID(database.DB, lecturerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil advisees",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    advisees,
	})
}

// CreateLecturer
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
		if strings.Contains(err.Error(), "duplicate key") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"success": false, "message": "Gagal: User ini sudah terdaftar sebagai dosen atau ID dosen sudah ada.", "error": err.Error(),
			})
		}

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