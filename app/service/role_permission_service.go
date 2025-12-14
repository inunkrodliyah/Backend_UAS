package service

import (
	"project-uas/app/model"
	"project-uas/app/repository"
	"project-uas/database"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetPermissionsByRoleID godoc
// @Summary      Lihat Permission milik Role
// @Description  Mendapatkan daftar permission yang dimiliki oleh Role tertentu
// @Tags         RolePermissions
// @Security     BearerAuth
// @Produce      json
// @Param        role_id   path      string  true  "Role ID (UUID)"
// @Success      200  {object}  fiber.Map{data=[]model.Permission}
// @Failure      400  {object}  fiber.Map
// @Failure      500  {object}  fiber.Map
// @Router       /roles/{role_id}/permissions [get]
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

// AssignPermissionToRole godoc
// @Summary      Assign Permission ke Role
// @Description  Menambahkan hak akses (permission) ke sebuah role
// @Tags         RolePermissions
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body model.AssignPermissionRequest true "Data Assign"
// @Success      201  {object}  fiber.Map
// @Failure      400  {object}  fiber.Map
// @Failure      500  {object}  fiber.Map
// @Router       /roles/permissions [post]
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

// RevokePermissionFromRole godoc
// @Summary      Revoke Permission dari Role
// @Description  Mencabut hak akses (permission) dari sebuah role
// @Tags         RolePermissions
// @Security     BearerAuth
// @Produce      json
// @Param        role_id        path      string  true  "Role ID (UUID)"
// @Param        permission_id  path      string  true  "Permission ID (UUID)"
// @Success      200  {object}  fiber.Map
// @Failure      400  {object}  fiber.Map
// @Failure      500  {object}  fiber.Map
// @Router       /roles/{role_id}/permissions/{permission_id} [delete]
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