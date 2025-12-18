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
// Login godoc
// @Summary      Login User
// @Description  Masuk sistem untuk mendapatkan Access Token & Refresh Token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body model.LoginRequest true "Email & Password"
// @Success      200  {object}  model.AuthResponse
// @Router /api/v1/auth/login [post]
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

	// Generate Access Token
accessToken, err := helper.GenerateToken(
	user.ID.String(),
	user.RoleID.String(),
	permissions,
)
if err != nil {
	return c.Status(500).JSON(fiber.Map{"message": "Token generation failed"})
}

// Generate Refresh Token
refreshToken, err := helper.GenerateRefreshToken(user.ID.String())
if err != nil {
	return c.Status(500).JSON(fiber.Map{"message": "Refresh token failed"})
}

return c.JSON(model.AuthResponse{
	Status: "success",
	Data: model.LoginResponseData{
		Token:        accessToken,
		RefreshToken: refreshToken,
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
// RefreshToken godoc
// @Summary      Refresh Access Token
// @Description  Mendapatkan token baru menggunakan Refresh Token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body model.RefreshTokenRequest true "Refresh Token"
// @Success      200  {object}  fiber.Map
// @Router /api/v1/auth/refresh [post]
func RefreshToken(c *fiber.Ctx) error {
	var req model.RefreshTokenRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid request",
		})
	}

	claims, err := helper.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid refresh token",
		})
	}

	userID := claims["user_id"].(string)

	user, err := repository.GetUserByID(database.DB, uuid.MustParse(userID))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "fail",
			"message": "User not found",
		})
	}

	permissions, _ := repository.GetPermissionNamesByRoleID(database.DB, user.RoleID)

	newAccessToken, err := helper.GenerateToken(
		user.ID.String(),
		user.RoleID.String(),
		permissions,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed generate token",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"access_token": newAccessToken,
		},
	})
}


// POST /api/v1/auth/logout
// Logout godoc
// @Summary      Logout
// @Description  Menghapus sesi login (Revoke Token)
// @Tags         Auth
// @Security     BearerAuth
// @Success      200  {object}  fiber.Map
// @Router /api/v1/auth/logout [post]
func Logout(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "success", "message": "Logged out successfully"})
}

// GET /api/v1/auth/profile
// GetProfile godoc
// @Summary      Get My Profile
// @Description  Melihat data diri user yang sedang login
// @Tags         Auth
// @Security     BearerAuth
// @Success      200  {object}  model.AuthResponse
// @Router /api/v1/auth/profile [get]
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