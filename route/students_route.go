package route

import (
	"project-uas/app/service"

	"github.com/gofiber/fiber/v2"
)

func SetupStudentRoutes(api fiber.Router) {
	students := api.Group("/students")

	students.Post("/", service.CreateStudent)
	students.Get("/:id", service.GetStudentByUserID)
	students.Put("/:id", service.UpdateStudent)
	students.Get("/", service.GetAllStudents)
	students.Get("/:id/achievements", service.GetStudentAchievements)
	students.Put("/:id/advisor", service.UpdateStudentAdvisor)
}
