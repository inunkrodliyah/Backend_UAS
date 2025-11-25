package service

import (
	"project-uas/app/model"
	"project-uas/app/repository"
	"project-uas/database"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetAllRoles(c *fiber.Ctx) error {
	roles, err := repository.GetAllRoles(database.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal mengambil data roles", "error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{"success": true, "data": roles})
}

func GetRoleByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "ID tidak valid"})
	}

	role, err := repository.GetRoleByID(database.DB, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "Role tidak ditemukan"})
	}
	return c.JSON(fiber.Map{"success": true, "data": role})
}

func CreateRole(c *fiber.Ctx) error {
	var req model.CreateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Request body tidak valid"})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Field 'name' wajib diisi"})
	}

	role := &model.Role{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := repository.CreateRole(database.DB, role); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal menambah role", "error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "message": "Role berhasil ditambahkan", "data": role})
}

func UpdateRole(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "ID tidak valid"})
	}

	var req model.UpdateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Request body tidak valid"})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Field 'name' wajib diisi"})
	}

	// Cek apakah role ada
	role, err := repository.GetRoleByID(database.DB, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "Role tidak ditemukan"})
	}

	// Update data
	role.Name = req.Name
	role.Description = req.Description

	if err := repository.UpdateRole(database.DB, role); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal mengupdate role", "error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Role berhasil diupdate", "data": role})
}

func DeleteRole(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "ID tidak valid"})
	}

	if err := repository.DeleteRole(database.DB, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal menghapus role", "error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Role berhasil dihapus"})
}