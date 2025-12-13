package route

import (
	"project-uas/middleware"
	"project-uas/app/service"

	"github.com/gofiber/fiber/v2"
)

func SetupLecturerRoutes(api fiber.Router) {
	lecturers := api.Group("/lecturers")

	// Pasang Middleware Auth
	lecturers.Use(middleware.AuthProtected)

	// HANYA ADMIN (Permission: lecturer:create)
	lecturers.Post("/", middleware.RequirePermission("lecturer:create"), service.CreateLecturer)
	lecturers.Get("/", service.GetAllLecturers)
	lecturers.Get("/:id", service.GetLecturerByUserID)
	lecturers.Put("/:id", service.UpdateLecturer)
	
	// Endpoint Khusus SRS (FR-006)
	lecturers.Get("/:id/advisees", service.GetLecturerAdvisees)
}