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

func GetAllStudents(c *fiber.Ctx) error {
	students, err := repository.GetAllStudents(database.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal mengambil data students",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    students,
	})
}

func GetStudentByUserID(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "ID [user_id] tidak valid",
		})
	}

	user, err := repository.GetUserByID(database.DB, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false, "message": "User tidak ditemukan",
		})
	}

	student, err := repository.GetStudentByID(database.DB, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false, "message": "Data student tidak ditemukan",
		})
	}

	// Response map disesuaikan dengan field yang benar
	response := fiber.Map{
		"user_id":        user.ID,
		"username":       user.Username,
		"email":          user.Email,
		"full_name":      user.FullName,
		"role_id":        user.RoleID,
		"is_active":      user.IsActive,
		"student_id":     student.StudentID,     
		"program_study":  student.ProgramStudy,  
		"academic_year":  student.AcademicYear,  
		"advisor_id":     student.AdvisorID,     
		"student_db_id":  student.ID,            
	}

	return c.JSON(fiber.Map{"success": true, "data": response})
}

func GetStudentAchievements(c *fiber.Ctx) error {
	studentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "ID student tidak valid",
		})
	}

	achievements, err := repository.GetAchievementsByStudentID(database.DB, studentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal mengambil data achievements",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    achievements,
	})
}

// CREATE STUDENT 
// HANYA membuat data student dan menautkannya ke user_id yang ada
func CreateStudent(c *fiber.Ctx) error {
	// 1. Menggunakan struct request yang baru
	var req model.CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "Request body tidak valid",
		})
	}

	// 2. Validasi field wajib
	if req.UserID == uuid.Nil || req.StudentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "Field wajib: 'user_id' dan 'student_id'",
		})
	}


	// 3. Siapkan struct model.Student untuk disimpan
	student := &model.Student{
		UserID:       req.UserID,
		StudentID:    req.StudentID,
		ProgramStudy: req.ProgramStudy,
		AcademicYear: req.AcademicYear,
		AdvisorID:    req.AdvisorID, 
	}

	// 5. Panggil repository HANYA untuk CreateStudent
	if err := repository.CreateStudent(database.DB, student); err != nil {
		log.Println("Error membuat student:", err)
		// Cek jika error karena duplikat user_id atau student_id
		if strings.Contains(err.Error(), "duplicate key") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"success": false, "message": "Gagal: User ini sudah terdaftar sebagai mahasiswa atau ID mahasiswa sudah ada.", "error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal membuat data student", "error": err.Error(),
		})
	}

	// 6. Kembalikan data student yang baru dibuat
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Data mahasiswa berhasil ditambahkan dan ditautkan ke user",
		"data":    student,
	})
}

func UpdateStudentAdvisor(c *fiber.Ctx) error {
	studentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "ID student tidak valid",
		})
	}

	var req model.UpdateAdvisorRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "Request body tidak valid",
		})
	}

	// Boleh null → untuk menghapus advisor
	err = repository.UpdateAdvisor(database.DB, studentID, req.AdvisorID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal mengubah advisor student",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Advisor student berhasil diperbarui",
		"data": fiber.Map{
			"student_id": studentID,
			"advisor_id": req.AdvisorID,
		},
	})
}

// UPDATE STUDENT
func UpdateStudent(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "ID [user_id] tidak valid",
		})
	}

	// 1. Menggunakan struct request 
	var req model.UpdateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "Request body tidak valid",
		})
	}

	// 2. Ambil data student yang ada
	student, err := repository.GetStudentByID(database.DB, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false, "message": "Data student tidak ditemukan",
		})
	}

	// 3. Update field jika diisi
	if req.StudentID != "" {
		student.StudentID = req.StudentID
	}
	if req.ProgramStudy != "" {
		student.ProgramStudy = req.ProgramStudy
	}
	if req.AcademicYear != "" {
		student.AcademicYear = req.AcademicYear
	}
	
	// Selalu update AdvisorID (bisa jadi di-set ke null)
	student.AdvisorID = req.AdvisorID

	// 4. Simpan perubahan
	if err := repository.UpdateStudent(database.DB, student); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal mengupdate data student", "error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Data student berhasil diupdate",
		"data":    student,
	})
}