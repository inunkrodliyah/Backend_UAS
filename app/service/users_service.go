package service

import (
	"project-uas/app/model"
	"project-uas/app/repository"
	"project-uas/database"
	"project-uas/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GET /api/v1/users
func GetAllUsers(c *fiber.Ctx) error {
	users, err := repository.GetAllUsers(database.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": users})
}

// GET /api/v1/users/:id
func GetUserByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid ID"})
	}
	user, err := repository.GetUserByID(database.DB, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "User not found"})
	}
	return c.JSON(fiber.Map{"success": true, "data": user})
}

// POST /api/v1/users
func CreateUser(c *fiber.Ctx) error {
	var req model.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid Body"})
	}

	hash, err := helper.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Hash failed"})
	}

	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hash,
		FullName:     req.FullName,
		RoleID:       req.RoleID,
	}

	if err := repository.CreateUser(database.DB, user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	user.PasswordHash = ""
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": user})
}

// PUT /api/v1/users/:id
func UpdateUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid ID"})
	}
	var req model.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid Body"})
	}

	user, err := repository.GetUserByID(database.DB, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "User not found"})
	}

	if req.Username != "" { user.Username = req.Username }
	if req.Email != "" { user.Email = req.Email }
	if req.FullName != "" { user.FullName = req.FullName }
	if req.RoleID != uuid.Nil { user.RoleID = req.RoleID }
	if req.IsActive != nil { user.IsActive = *req.IsActive }

	if err := repository.UpdateUser(database.DB, user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": user})
}

// DELETE /api/v1/users/:id
func DeleteUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid ID"})
	}
	if err := repository.DeleteUser(database.DB, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "User deleted"})
}

// PUT /api/v1/users/:id/role
func UpdateUserRole(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid ID"})
	}
	var req model.UpdateUserRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid Body"})
	}

	if err := repository.UpdateUserRole(database.DB, id, req.RoleID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Role updated"})
}