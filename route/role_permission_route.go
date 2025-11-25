package route

import (
	"project-uas/app/service"

	"github.com/gofiber/fiber/v2"
)

func SetupRolePermissionRoutes(api fiber.Router) {
	rp := api.Group("/role-permissions")

	// Mendapat semua permission yang dimiliki role
	rp.Get("/:role_id", service.GetPermissionsByRoleID)
	// Memberi permission ke role
	rp.Post("/", service.AssignPermissionToRole)
	// Menghapus permission dari role
	rp.Delete("/:role_id/:permission_id", service.RevokePermissionFromRole)
}