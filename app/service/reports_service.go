package service

import (
	"project-uas/app/model"
	"project-uas/app/repository"
	"project-uas/database"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GET /api/v1/reports/statistics
// GetSystemStatistics godoc
// @Summary      Statistik Sistem
// @Description  Menampilkan ringkasan total mahasiswa, dosen, dan statistik prestasi (berdasarkan status & tipe)
// @Tags         Reports
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  fiber.Map{data=model.SystemStatisticsResponse}
// @Failure      500  {object}  fiber.Map
// @Router       /reports/statistics [get]
func GetSystemStatistics(c *fiber.Ctx) error {
	// 1. Ambil Data dari PostgreSQL
	totalStudents, _ := repository.CountTotalUsersByRole(database.DB, "student")
	totalLecturers, _ := repository.CountTotalUsersByRole(database.DB, "lecturer")
	statsByStatus, _ := repository.CountAchievementsByStatus(database.DB)

	// 2. Ambil Data dari MongoDB (Agregasi Tipe)
	statsByType, _ := repository.AggregateAchievementsByType(database.MongoDB)

	// Hitung total prestasi
	totalAchievements := 0
	for _, count := range statsByStatus {
		totalAchievements += count
	}

	// 3. Gabungkan Response
	response := model.SystemStatisticsResponse{
		TotalStudents:     totalStudents,
		TotalLecturers:    totalLecturers,
		TotalAchievements: totalAchievements,
		ByStatus:          statsByStatus,
		ByType:            statsByType,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GET /api/v1/reports/student/:id
// GetStudentReport godoc
// @Summary      Laporan Prestasi Mahasiswa
// @Description  Mendapatkan detail profil, total poin, dan daftar prestasi mahasiswa tertentu berdasarkan User ID
// @Tags         Reports
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "User ID (UUID)"
// @Success      200  {object}  fiber.Map{data=model.StudentReportResponse}
// @Failure      400  {object}  fiber.Map
// @Failure      404  {object}  fiber.Map
// @Router       /reports/student/{id} [get]
func GetStudentReport(c *fiber.Ctx) error {
	// Param ID bisa berupa StudentID (UUID table students) atau UserID.
	// Asumsi disini adalah UserID (dari login).
	paramID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "ID tidak valid"})
	}

	// 1. Ambil Profil Mahasiswa (Postgres)
	student, err := repository.GetStudentByID(database.DB, paramID) // Menggunakan func yg ada di students_repository
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Mahasiswa tidak ditemukan"})
	}
	// Ambil nama user juga
	user, _ := repository.GetUserByID(database.DB, student.UserID)

	// 2. Ambil List Prestasi Mahasiswa (Postgres)
	refs, err := repository.GetAchievementsByStudentID(database.DB, student.ID)
	if err != nil {
		refs = []model.AchievementReference{}
	}

	// 3. Ambil Total Poin (Mongo) - Menggunakan UUID student (string)
	totalPoints, _ := repository.SumStudentPoints(database.MongoDB, student.ID.String())

	// 4. Buat List Prestasi Ringkas (Gabung Status Postgres + Detail Mongo)
	var simpleList []model.SimpleAchievementView
	for _, ref := range refs {
		// Ambil detail judul & points dari mongo satu per satu (bisa dioptimasi dengan query $in)
		detail, _ := repository.GetAchievementMongoByID(database.MongoDB, ref.MongoAchievementID)
		
		title := "Unknown Title"
		aType := "Unknown"
		points := 0
		if detail != nil {
			title = detail.Title
			aType = detail.AchievementType
			points = detail.Points
		}

		simpleList = append(simpleList, model.SimpleAchievementView{
			ID:              ref.ID,
			Title:           title,
			AchievementType: aType,
			Status:          string(ref.Status),
			Points:          points,
		})
	}

	// 5. Response Final
	response := model.StudentReportResponse{
		StudentID:         student.StudentID, // NIM
		FullName:          user.FullName,
		ProgramStudy:      student.ProgramStudy,
		TotalAchievements: len(refs),
		TotalPoints:       totalPoints,
		Achievements:      simpleList,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}