package route

import (
	"project-uas/app/service"

	"github.com/gofiber/fiber/v2"
)

func SetupLecturerRoutes(api fiber.Router) {
	lecturers := api.Group("/lecturers")
	lecturers.Post("/", service.CreateLecturer)
	lecturers.Get("/:id", service.GetLecturerByUserID)
	lecturers.Put("/:id", service.UpdateLecturer)
	lecturers.Get("/", service.GetAllLecturers)
	lecturers.Get("/:id/advisees", service.GetLecturerAdvisees)
}