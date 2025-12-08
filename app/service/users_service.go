package service

import (
	"project-uas/app/model"
	"project-uas/app/repository"
	"project-uas/database"
	"project-uas/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetAllUsers(c *fiber.Ctx) error {
	users, err := repository.GetAllUsers(database.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal mengambil data users", "error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{"success": true, "data": users})
}

func GetUserByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "ID tidak valid"})
	}

	user, err := repository.GetUserByID(database.DB, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "User tidak ditemukan"})
	}
	return c.JSON(fiber.Map{"success": true, "data": user})
}

func CreateUser(c *fiber.Ctx) error {
	var req model.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Request body tidak valid"})
	}

	if req.Username == "" || req.Email == "" || req.Password == "" || req.FullName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Field wajib: username, email, password, full_name"})
	}

	// Hash password
	passwordHash, err := helper.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Gagal hashing password"})
	}

	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
		FullName:     req.FullName,
		RoleID:       req.RoleID,
	}

	if err := repository.CreateUser(database.DB, user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal menambah user", "error": err.Error(),
		})
	}

	// Sembunyikan password hash dari response
	user.PasswordHash = ""

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "message": "User berhasil ditambahkan", "data": user})
}

func UpdateUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "ID tidak valid"})
	}

	var req model.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Request body tidak valid"})
	}

	user, err := repository.GetUserByID(database.DB, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "User tidak ditemukan"})
	}

	// Update field
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.RoleID != uuid.Nil {
		user.RoleID = req.RoleID
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := repository.UpdateUser(database.DB, user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal mengupdate user", "error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{"success": true, "message": "User berhasil diupdate", "data": user})
}

func DeleteUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "ID tidak valid"})
	}

	if err := repository.DeleteUser(database.DB, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal menghapus user", "error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{"success": true, "message": "User berhasil dihapus"})
}
//yang bagian put id/role
func UpdateUserRole(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "ID tidak valid",
		})
	}

	var req model.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "Request body tidak valid",
		})
	}

	if req.RoleID == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "role_id wajib diisi",
		})
	}

	// Cek apakah user ada
	_, err = repository.GetUserByID(database.DB, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false, "message": "User tidak ditemukan",
		})
	}

	// Update role
	if err := repository.UpdateUserRole(database.DB, id, req.RoleID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal mengubah role user", "error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Role user berhasil diubah",
		"data": fiber.Map{
			"id":      id,
			"role_id": req.RoleID,
		},
	})
}

func Login(c *fiber.Ctx) error {
	// 1. Parsing Input
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail", 
			"message": "Request body tidak valid",
		})
	}

	// 2. Validasi Input
	if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail", 
			"message": "Username dan Password wajib diisi",
		})
	}

	// 3. Cari User di DB
	user, err := repository.GetUserByUsername(database.DB, req.Username)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "fail", 
			"message": "Username atau password salah",
		})
	}

	// 4. Cek Status Aktif
	if !user.IsActive {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "fail", 
			"message": "Akun dinonaktifkan",
		})
	}

	// 5. Cek Password
	match := helper.CheckPasswordHash(req.Password, user.PasswordHash)
	if !match {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "fail", 
			"message": "Username atau password salah",
		})
	}

	// 6. Generate Token
	token, err := helper.GenerateToken(user.ID.String(), user.RoleID.String())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error", 
			"message": "Gagal generate token",
		})
	}

	// 7. SIAPKAN RESPONSE (SESUAI STRUKTUR SRS)
	response := model.LoginResponse{
		Status: "success",
		Data: model.LoginResponseData{
			Token:        token,
			RefreshToken: "", // Kosongkan dulu (Dummy) karena tabel refresh token belum ada
			User: model.UserLoginData{
				ID:       user.ID,
				Username: user.Username,
				FullName: user.FullName,
				RoleID:   user.RoleID,
			},
		},
	}

	// 8. Return JSON
	return c.JSON(response)
}