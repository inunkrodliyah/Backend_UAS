package route

import (
	"project-uas/middleware"
	"project-uas/app/service"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(api fiber.Router) {
	users := api.Group("/users")

	// Endpoint Users (5.2) - PROTECTED (Admin Only)
	users.Use(middleware.AuthProtected)
    // Optional: Tambahkan middleware.RequirePermission("user:manage") jika ingin lebih ketat

	users.Get("/", service.GetAllUsers)
	users.Get("/:id", service.GetUserByID)
	users.Post("/", service.CreateUser)
	users.Put("/:id", service.UpdateUser)
	users.Delete("/:id", service.DeleteUser)
	users.Put("/:id/role", service.UpdateUserRole)
}