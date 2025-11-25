package service

import (
	"project-uas/app/model"
	"project-uas/app/repository"
	"project-uas/database"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetPermissionsByRoleID(c *fiber.Ctx) error {
	roleID, err := uuid.Parse(c.Params("role_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Role ID tidak valid"})
	}

	permissions, err := repository.GetPermissionsByRoleID(database.DB, roleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal mengambil permissions", "error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{"success": true, "data": permissions})
}

func AssignPermissionToRole(c *fiber.Ctx) error {
	var req model.AssignPermissionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Request body tidak valid"})
	}

	rp := &model.RolePermission{
		RoleID:       req.RoleID,
		PermissionID: req.PermissionID,
	}

	if err := repository.AssignPermissionToRole(database.DB, rp); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal assign permission", "error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "message": "Permission berhasil di-assign"})
}

func RevokePermissionFromRole(c *fiber.Ctx) error {
	roleID, err := uuid.Parse(c.Params("role_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Role ID tidak valid"})
	}

	permissionID, err := uuid.Parse(c.Params("permission_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Permission ID tidak valid"})
	}

	if err := repository.RevokePermissionFromRole(database.DB, roleID, permissionID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal revoke permission", "error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Permission berhasil di-revoke"})
}