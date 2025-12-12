package route

import (
	"project-uas/middleware"
	"project-uas/app/service"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(api fiber.Router) {
	auth := api.Group("/auth")

	// Endpoint Public (5.1)
	auth.Post("/login", service.Login)
	auth.Post("/refresh", service.RefreshToken)

	// Endpoint Protected (Butuh Token)
	auth.Use(middleware.AuthProtected)
	auth.Post("/logout", service.Logout)
	auth.Get("/profile", service.GetProfile)
}