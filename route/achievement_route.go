package route

import (
	"project-uas/middleware"
	"project-uas/app/service"

	"github.com/gofiber/fiber/v2"
)

func SetupAchievementRoutes(api fiber.Router) {
	achievements := api.Group("/achievements")

	// Middleware Auth Wajib
	achievements.Use(middleware.AuthProtected)

	// List & Detail
	achievements.Get("/", service.ListAchievements)
	achievements.Get("/:id", service.GetAchievementDetail)

	// Mahasiswa Actions
	achievements.Post("/", service.CreateAchievement)           // Create Draft
	achievements.Put("/:id", service.UpdateAchievement)         // Update Draft
	achievements.Delete("/:id", service.DeleteAchievement)      // Delete Draft
	achievements.Post("/:id/submit", service.SubmitAchievement) // Submit to Dosen
	achievements.Post("/:id/attachments", service.UploadAttachment) // Upload File

	// Dosen Actions
	achievements.Post("/:id/verify", service.VerifyAchievement) // Verify
	achievements.Post("/:id/reject", service.RejectAchievement) // Reject

	// History
	achievements.Get("/:id/history", service.GetAchievementHistory)
}