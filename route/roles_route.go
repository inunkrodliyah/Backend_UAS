package route

import (
	"project-uas/app/service"

	"github.com/gofiber/fiber/v2"
)

func SetupRoleRoutes(api fiber.Router) {
	roles := api.Group("/roles")

	roles.Get("/", service.GetAllRoles)
	roles.Get("/:id", service.GetRoleByID)
	roles.Post("/", service.CreateRole)
	roles.Put("/:id", service.UpdateRole)
	roles.Delete("/:id", service.DeleteRole)
}