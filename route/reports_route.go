package route

import (
	"project-uas/middleware"
	"project-uas/app/service"

	"github.com/gofiber/fiber/v2"
)

func SetupReportRoutes(api fiber.Router) {
	reports := api.Group("/reports")

	// Pasang Middleware (Wajib Login)
	reports.Use(middleware.AuthProtected)

	// 5.8 Reports & Analytics
	reports.Get("/statistics", service.GetSystemStatistics) // Dashboard Admin/Dosen
	reports.Get("/student/:id", service.GetStudentReport)   // Raport Mahasiswa (By User ID)
}