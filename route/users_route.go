package route

import (
	"project-uas/app/service"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(api fiber.Router) {
	users := api.Group("/users")

	users.Get("/", service.GetAllUsers)
	users.Get("/:id", service.GetUserByID)
	users.Post("/", service.CreateUser) // Ini user general
	users.Put("/:id", service.UpdateUser)
	users.Delete("/:id", service.DeleteUser)
}