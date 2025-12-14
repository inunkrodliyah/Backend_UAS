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

// GetAllStudents godoc
// @Summary      Lihat Semua Mahasiswa
// @Description  Mendapatkan daftar semua mahasiswa yang terdaftar
// @Tags         Students
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  fiber.Map{data=[]model.Student}
// @Failure      500  {object}  fiber.Map
// @Router       /students [get]
func GetAllStudents(c *fiber.Ctx) error {
	students, err := repository.GetAllStudents(database.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal mengambil data students",
		})
	}
	return c.JSON(fiber.Map{"success": true, "data": students})
}

// GetStudentByUserID godoc
// @Summary      Detail Mahasiswa by ID
// @Description  Melihat detail profil mahasiswa beserta data usernya berdasarkan User ID
// @Tags         Students
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "User ID (UUID)"
// @Success      200  {object}  fiber.Map
// @Failure      400  {object}  fiber.Map
// @Failure      404  {object}  fiber.Map
// @Router       /students/{id} [get]
func GetStudentByUserID(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "ID tidak valid"})
	}

	user, err := repository.GetUserByID(database.DB, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "User tidak ditemukan"})
	}

	student, err := repository.GetStudentByID(database.DB, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "Data student tidak ditemukan"})
	}

	response := fiber.Map{
		"user_id":       user.ID,
		"username":      user.Username,
		"email":         user.Email,
		"full_name":     user.FullName,
		"role_id":       user.RoleID,
		"student_id":    student.StudentID,
		"program_study": student.ProgramStudy,
		"academic_year": student.AcademicYear,
		"advisor_id":    student.AdvisorID,
		"student_db_id": student.ID,
	}

	return c.JSON(fiber.Map{"success": true, "data": response})
}

// Mengambil prestasi berdasarkan ID User Mahasiswa
// GetStudentAchievements godoc
// @Summary      Lihat Prestasi Mahasiswa
// @Description  Mendapatkan daftar prestasi milik mahasiswa tertentu berdasarkan User ID
// @Tags         Students
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "User ID (UUID)"
// @Success      200  {object}  fiber.Map{data=[]model.AchievementReference}
// @Failure      400  {object}  fiber.Map
// @Failure      404  {object}  fiber.Map
// @Router       /students/{id}/achievements [get]
func GetStudentAchievements(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "ID tidak valid"})
	}

	// 1. Ambil data student dulu untuk dapat StudentPK (ID tabel student)
	student, err := repository.GetStudentByID(database.DB, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "Mahasiswa tidak ditemukan"})
	}

	// 2. Ambil achievement berdasarkan StudentPK
	achievements, err := repository.GetAchievementsByStudentID(database.DB, student.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Gagal mengambil data achievements"})
	}

	return c.JSON(fiber.Map{"success": true, "data": achievements})
}

// CreateStudent godoc
// @Summary      Tambah Data Mahasiswa
// @Description  Menambahkan profil mahasiswa ke user yang sudah ada (Admin)
// @Tags         Students
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body model.CreateStudentRequest true "Data Mahasiswa"
// @Success      201  {object}  fiber.Map{data=model.Student}
// @Failure      400  {object}  fiber.Map
// @Failure      409  {object}  fiber.Map
// @Router       /students [post]
func CreateStudent(c *fiber.Ctx) error {
	var req model.CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Request body tidak valid"})
	}

	if req.UserID == uuid.Nil || req.StudentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Field wajib: 'user_id' dan 'student_id'"})
	}

	student := &model.Student{
		UserID:       req.UserID,
		StudentID:    req.StudentID,
		ProgramStudy: req.ProgramStudy,
		AcademicYear: req.AcademicYear,
		AdvisorID:    req.AdvisorID,
	}

	if err := repository.CreateStudent(database.DB, student); err != nil {
		log.Println("Error membuat student:", err)
		if strings.Contains(err.Error(), "duplicate key") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"success": false, "message": "Duplikat data", "error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Gagal membuat student", "error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "message": "Data mahasiswa berhasil ditambahkan", "data": student})
}

// UpdateStudent godoc
// @Summary      Update Profil Mahasiswa
// @Description  Mengubah data NIM, Prodi, atau Dosen Wali
// @Tags         Students
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID (UUID)"
// @Param        request body model.UpdateStudentRequest true "Data Update"
// @Success      200  {object}  fiber.Map{data=model.Student}
// @Failure      400  {object}  fiber.Map
// @Failure      404  {object}  fiber.Map
// @Router       /students/{id} [put]
func UpdateStudent(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "ID tidak valid"})
	}

	var req model.UpdateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Request body tidak valid"})
	}

	student, err := repository.GetStudentByID(database.DB, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "Data student tidak ditemukan"})
	}

	if req.StudentID != "" { student.StudentID = req.StudentID }
	if req.ProgramStudy != "" { student.ProgramStudy = req.ProgramStudy }
	if req.AcademicYear != "" { student.AcademicYear = req.AcademicYear }
	
	student.AdvisorID = req.AdvisorID // Bisa null

	if err := repository.UpdateStudent(database.DB, student); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Gagal update", "error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Data student berhasil diupdate", "data": student})
}

// Update Advisor Only
// UpdateStudentAdvisor godoc
// @Summary      Update Dosen Wali
// @Description  Mengubah dosen wali mahasiswa tertentu
// @Tags         Students
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID (UUID)"
// @Param        request body model.UpdateAdvisorRequest true "Dosen Wali ID (Lecturer UUID)"
// @Success      200  {object}  fiber.Map
// @Failure      400  {object}  fiber.Map
// @Failure      500  {object}  fiber.Map
// @Router       /students/{id}/advisor [put]
func UpdateStudentAdvisor(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id")) // Ini adalah User ID dari URL
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "ID tidak valid"})
	}

	var req model.UpdateAdvisorRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Request body tidak valid"})
	}

	// Update menggunakan UserID
	err = repository.UpdateAdvisor(database.DB, userID, req.AdvisorID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Gagal mengubah advisor"})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Advisor berhasil diperbarui"})
}