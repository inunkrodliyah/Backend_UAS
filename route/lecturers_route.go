package route

import (
	"project-uas/app/service"

	"github.com/gofiber/fiber/v2"
)

func SetupLecturerRoutes(api fiber.Router) {
	lecturers := api.Group("/lecturers")

	// --- Rute POST diubah ---
	// Sebelumnya: service.CreateLecturerUser
	// Sekarang: service.CreateLecturer
	lecturers.Post("/", service.CreateLecturer)

	// Rute ini tetap sama (menggunakan user_id sebagai param 'id')
	lecturers.Get("/:id", service.GetLecturerByUserID)
	lecturers.Put("/:id", service.UpdateLecturer)
}