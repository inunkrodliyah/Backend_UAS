package service

import (
	"project-uas/app/model"
	"project-uas/app/repository"
	"project-uas/database"
	"strings" // Di-import untuk error handling

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetAllPermissions(c *fiber.Ctx) error {
	permissions, err := repository.GetAllPermissions(database.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal mengambil data permissions", "error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{"success": true, "data": permissions})
}

func GetPermissionByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "ID tidak valid"})
	}

	permission, err := repository.GetPermissionByID(database.DB, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "Permission tidak ditemukan"})
	}
	return c.JSON(fiber.Map{"success": true, "data": permission})
}

func CreatePermission(c *fiber.Ctx) error {
	// 1. Menggunakan struct request yang baru
	var req model.CreatePermissionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Request body tidak valid"})
	}

	// 2. Validasi diubah: +resource, +action (NOT NULL)
	if req.Name == "" || req.Resource == "" || req.Action == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "Field 'name', 'resource', dan 'action' wajib diisi",
		})
	}

	// 3. Model diisi dengan data baru
	permission := &model.Permission{
		Name:        req.Name,
		Resource:    req.Resource,    // <-- DITAMBAHKAN
		Action:      req.Action,      // <-- DITAMBAHKAN
		Description: req.Description,
	}

	// 4. Panggil repository (yang sudah diperbarui)
	if err := repository.CreatePermission(database.DB, permission); err != nil {
		// Cek error duplikat 'name'
		if strings.Contains(err.Error(), "duplicate key") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"success": false, "message": "Gagal: Nama permission sudah ada.", "error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal menambah permission", "error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "message": "Permission berhasil ditambahkan", "data": permission})
}

func UpdatePermission(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "ID tidak valid"})
	}

	// 1. Menggunakan struct request yang baru
	var req model.UpdatePermissionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Request body tidak valid"})
	}

	// 2. Validasi diubah: +resource, +action (NOT NULL)
	if req.Name == "" || req.Resource == "" || req.Action == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "Field 'name', 'resource', dan 'action' wajib diisi",
		})
	}

	// 3. Ambil data lama
	permission, err := repository.GetPermissionByID(database.DB, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "Permission tidak ditemukan"})
	}

	// 4. Update data
	permission.Name = req.Name
	permission.Resource = req.Resource // <-- DITAMBAHKAN
	permission.Action = req.Action     // <-- DITAMBAHKAN
	permission.Description = req.Description

	// 5. Panggil repository (yang sudah diperbarui)
	if err := repository.UpdatePermission(database.DB, permission); err != nil {
		// Cek error duplikat 'name'
		if strings.Contains(err.Error(), "duplicate key") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"success": false, "message": "Gagal: Nama permission sudah ada.", "error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal mengupdate permission", "error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Permission berhasil diupdate", "data": permission})
}

func DeletePermission(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "ID tidak valid"})
	}

	if err := repository.DeletePermission(database.DB, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal menghapus permission", "error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Permission berhasil dihapus"})
}