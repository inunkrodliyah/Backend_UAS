package route

import (
	"project-uas/app/service"

	"github.com/gofiber/fiber/v2"
)

func SetupStudentRoutes(api fiber.Router) {
	students := api.Group("/students")

	// --- Rute POST diubah ---
	// Sebelumnya: service.CreateStudentUser
	// Sekarang: service.CreateStudent
	students.Post("/", service.CreateStudent)

	// Mendapat data gabungan user + student (via user_id)
	students.Get("/:id", service.GetStudentByUserID)
	
	// Mengupdate data student (via user_id)
	students.Put("/:id", service.UpdateStudent)
}