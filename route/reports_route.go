package route

import (
	"project-uas/middleware"
	"project-uas/app/service"

	"github.com/gofiber/fiber/v2"
)

func SetupReportRoutes(api fiber.Router) {
    reports := api.Group("/reports")

    reports.Use(middleware.AuthProtected)
    
    // Tambahkan baris ini agar endpoint dicek permission-nya
    reports.Use(middleware.RequirePermission("report:read"))

    reports.Get("/statistics", service.GetSystemStatistics)
    reports.Get("/student/:id", service.GetStudentReport)
}