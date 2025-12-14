package service

import (
	"project-uas/app/model"
	"project-uas/app/repository"
	"project-uas/database"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetAllRoles godoc
// @Summary      Lihat Semua Role
// @Description  Mendapatkan daftar semua role yang tersedia di sistem
// @Tags         Roles
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  fiber.Map{data=[]model.Role}
// @Failure      500  {object}  fiber.Map
// @Router       /roles [get]
func GetAllRoles(c *fiber.Ctx) error {
	roles, err := repository.GetAllRoles(database.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal mengambil data roles", "error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{"success": true, "data": roles})
}

// GetRoleByID godoc
// @Summary      Detail Role
// @Description  Mendapatkan detail role berdasarkan ID
// @Tags         Roles
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Role ID (UUID)"
// @Success      200  {object}  fiber.Map{data=model.Role}
// @Failure      400  {object}  fiber.Map
// @Failure      404  {object}  fiber.Map
// @Router       /roles/{id} [get]
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

// CreateRole godoc
// @Summary      Buat Role Baru
// @Description  Menambahkan role baru ke dalam sistem
// @Tags         Roles
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body model.CreateRoleRequest true "Data Role"
// @Success      201  {object}  fiber.Map{data=model.Role}
// @Failure      400  {object}  fiber.Map
// @Failure      500  {object}  fiber.Map
// @Router       /roles [post]
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

// UpdateRole godoc
// @Summary      Update Role
// @Description  Mengubah data role (Nama dan Deskripsi)
// @Tags         Roles
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Role ID (UUID)"
// @Param        request body model.UpdateRoleRequest true "Data Update"
// @Success      200  {object}  fiber.Map{data=model.Role}
// @Failure      400  {object}  fiber.Map
// @Failure      404  {object}  fiber.Map
// @Failure      500  {object}  fiber.Map
// @Router       /roles/{id} [put]
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

// DeleteRole godoc
// @Summary      Hapus Role
// @Description  Menghapus role dari sistem berdasarkan ID
// @Tags         Roles
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Role ID (UUID)"
// @Success      200  {object}  fiber.Map
// @Failure      400  {object}  fiber.Map
// @Failure      500  {object}  fiber.Map
// @Router       /roles/{id} [delete]
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