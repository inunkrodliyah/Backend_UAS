package route

import (
	"project-uas/middleware"
	"project-uas/app/service"

	"github.com/gofiber/fiber/v2"
)

func SetupStudentRoutes(api fiber.Router) {
	students := api.Group("/students")

	// Pasang Middleware Auth
	students.Use(middleware.AuthProtected)

	students.Post("/", service.CreateStudent)
	students.Get("/", service.GetAllStudents)
	students.Get("/:id", service.GetStudentByUserID)
	students.Put("/:id", service.UpdateStudent)
	
	// Endpoint Khusus SRS
	students.Get("/:id/achievements", service.GetStudentAchievements)
	students.Put("/:id/advisor", service.UpdateStudentAdvisor)
}