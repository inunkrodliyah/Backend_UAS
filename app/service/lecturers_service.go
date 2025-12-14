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

// GetAllLecturers godoc
// @Summary      Lihat Semua Dosen
// @Description  Mendapatkan daftar semua dosen yang terdaftar
// @Tags         Lecturers
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  fiber.Map{data=[]model.Lecturer}
// @Failure      500  {object}  fiber.Map
// @Router       /lecturers [get]
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
// GetLecturerByUserID godoc
// @Summary      Detail Dosen by ID
// @Description  Melihat detail data dosen beserta info user-nya berdasarkan User ID
// @Tags         Lecturers
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "User ID (UUID)"
// @Success      200  {object}  fiber.Map
// @Failure      400  {object}  fiber.Map
// @Failure      404  {object}  fiber.Map
// @Router       /lecturers/{id} [get]
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

// GetLecturerAdvisees godoc
// @Summary      Lihat Mahasiswa Bimbingan
// @Description  Melihat daftar mahasiswa yang dibimbing oleh dosen tertentu
// @Tags         Lecturers
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Lecturer/User UUID"
// @Success      200  {object}  fiber.Map{data=[]model.Student}
// @Failure      400  {object}  fiber.Map
// @Failure      404  {object}  fiber.Map
// @Router       /lecturers/{id}/advisees [get]
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
// CreateLecturer godoc
// @Summary      Tambah Data Dosen
// @Description  Menambahkan profil dosen ke user yang sudah ada (Admin)
// @Tags         Lecturers
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body model.CreateLecturerRequest true "Data Dosen"
// @Success      201  {object}  fiber.Map{data=model.Lecturer}
// @Failure      400  {object}  fiber.Map
// @Failure      409  {object}  fiber.Map
// @Router       /lecturers [post]
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
// UpdateLecturer godoc
// @Summary      Update Profil Dosen
// @Description  Mengubah data NIP/NIDN atau Departemen Dosen
// @Tags         Lecturers
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID (UUID)"
// @Param        request body model.UpdateLecturerRequest true "Data Update"
// @Success      200  {object}  fiber.Map{data=model.Lecturer}
// @Failure      400  {object}  fiber.Map
// @Failure      404  {object}  fiber.Map
// @Router       /lecturers/{id} [put]
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