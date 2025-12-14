package service

import (
	"project-uas/app/model"
	"project-uas/app/repository"
	"project-uas/database"
	"strings" // Di-import untuk error handling

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetAllPermissions godoc
// @Summary      Lihat Semua Permission
// @Description  Mendapatkan daftar semua hak akses (permission) yang tersedia di sistem
// @Tags         Permissions
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  fiber.Map{data=[]model.Permission}
// @Failure      500  {object}  fiber.Map
// @Router       /permissions [get]
func GetAllPermissions(c *fiber.Ctx) error {
	permissions, err := repository.GetAllPermissions(database.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal mengambil data permissions", "error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{"success": true, "data": permissions})
}

// GetPermissionByID godoc
// @Summary      Detail Permission
// @Description  Mendapatkan detail permission berdasarkan ID
// @Tags         Permissions
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Permission ID (UUID)"
// @Success      200  {object}  fiber.Map{data=model.Permission}
// @Failure      400  {object}  fiber.Map
// @Failure      404  {object}  fiber.Map
// @Router       /permissions/{id} [get]
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

// CreatePermission godoc
// @Summary      Buat Permission Baru
// @Description  Menambahkan permission baru (Resource + Action) ke dalam sistem
// @Tags         Permissions
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body model.CreatePermissionRequest true "Data Permission"
// @Success      201  {object}  fiber.Map{data=model.Permission}
// @Failure      400  {object}  fiber.Map
// @Failure      409  {object}  fiber.Map
// @Failure      500  {object}  fiber.Map
// @Router       /permissions [post]
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
		Resource:    req.Resource,    
		Action:      req.Action,      
		Description: req.Description,
	}

	// 4. Panggil repository 
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

// UpdatePermission godoc
// @Summary      Update Permission
// @Description  Mengubah data permission (Name, Resource, Action, Description)
// @Tags         Permissions
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Permission ID (UUID)"
// @Param        request body model.UpdatePermissionRequest true "Data Update"
// @Success      200  {object}  fiber.Map{data=model.Permission}
// @Failure      400  {object}  fiber.Map
// @Failure      404  {object}  fiber.Map
// @Failure      409  {object}  fiber.Map
// @Failure      500  {object}  fiber.Map
// @Router       /permissions/{id} [put]
func UpdatePermission(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "ID tidak valid"})
	}

	// 1. Menggunakan struct request 
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
	permission.Resource = req.Resource 
	permission.Action = req.Action     
	permission.Description = req.Description

	// 5. Panggil repository 
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

// DeletePermission godoc
// @Summary      Hapus Permission
// @Description  Menghapus permission dari sistem berdasarkan ID
// @Tags         Permissions
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Permission ID (UUID)"
// @Success      200  {object}  fiber.Map
// @Failure      400  {object}  fiber.Map
// @Failure      500  {object}  fiber.Map
// @Router       /permissions/{id} [delete]
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