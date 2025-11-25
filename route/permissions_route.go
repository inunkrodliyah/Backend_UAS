package route

import (
	"project-uas/app/service"

	"github.com/gofiber/fiber/v2"
)

func SetupPermissionRoutes(api fiber.Router) {
	permissions := api.Group("/permissions")

	permissions.Get("/", service.GetAllPermissions)
	permissions.Get("/:id", service.GetPermissionByID)
	permissions.Post("/", service.CreatePermission)
	permissions.Put("/:id", service.UpdatePermission)
	permissions.Delete("/:id", service.DeletePermission)
}