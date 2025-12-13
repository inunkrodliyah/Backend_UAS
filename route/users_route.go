package route

import (
	"project-uas/middleware"
	"project-uas/app/service"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(api fiber.Router) {
	users := api.Group("/users")

	// 1. Semua endpoint users butuh Login
	users.Use(middleware.AuthProtected)

	// 2. Pasang Permission Spesifik per Endpoint (Sesuai FR-009)

	// GET /api/v1/users (Butuh: user:read)
	users.Get("/", middleware.RequirePermission("user:read"), service.GetAllUsers)

	// GET /api/v1/users/:id (Butuh: user:read)
	users.Get("/:id", middleware.RequirePermission("user:read"), service.GetUserByID)

	// POST /api/v1/users (Butuh: user:create)
	users.Post("/", middleware.RequirePermission("user:create"), service.CreateUser)

	// PUT /api/v1/users/:id (Butuh: user:update)
	users.Put("/:id", middleware.RequirePermission("user:update"), service.UpdateUser)

	// DELETE /api/v1/users/:id (Butuh: user:delete)
	users.Delete("/:id", middleware.RequirePermission("user:delete"), service.DeleteUser)

	// PUT /api/v1/users/:id/role (Butuh: user:assign_role)
	users.Put("/:id/role", middleware.RequirePermission("user:assign_role"), service.UpdateUserRole)
}