package service

import (
	"project-uas/app/model"
	"project-uas/app/repository"
	"project-uas/database"
	"project-uas/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// POST /api/v1/auth/login
func Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Request body invalid"})
	}

	// 1. Validasi User
	user, err := repository.GetUserByUsername(database.DB, req.Username)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "Invalid credentials"})
	}

	// 2. Cek Password
	if !helper.CheckPasswordHash(req.Password, user.PasswordHash) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "Invalid credentials"})
	}

	// 3. Cek Active Status
	if !user.IsActive {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "User inactive"})
	}

	// 4. Ambil Permissions (FR-001 Flow 4)
	permissions, _ := repository.GetPermissionNamesByRoleID(database.DB, user.RoleID)
	if permissions == nil {
		permissions = []string{}
	}

	// 5. Generate Token dengan Permissions
	token, err := helper.GenerateToken(user.ID.String(), user.RoleID.String(), permissions)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Token generation failed"})
	}

	// 6. Return Response
	return c.JSON(model.AuthResponse{
		Status: "success",
		Data: model.LoginResponseData{
			Token:        token,
			RefreshToken: "dummy-refresh-token", // Placeholder
			User: model.UserLoginData{
				ID:          user.ID,
				Username:    user.Username,
				FullName:    user.FullName,
				RoleID:      user.RoleID,
				Permissions: permissions,
			},
		},
	})
}

// POST /api/v1/auth/refresh
func RefreshToken(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "success", "message": "Token refreshed (Logic Placeholder)"})
}

// POST /api/v1/auth/logout
func Logout(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "success", "message": "Logged out successfully"})
}

// GET /api/v1/auth/profile
func GetProfile(c *fiber.Ctx) error {
	// Ambil user_id dari middleware
	userIDStr := c.Locals("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	user, err := repository.GetUserByID(database.DB, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "User not found"})
	}

	return c.JSON(model.AuthResponse{
		Status: "success",
		Data:   user,
	})
}