package route

import (
	"project-uas/app/service"

	"github.com/gofiber/fiber/v2"
)

func SetupAchievementReferenceRoutes(api fiber.Router) {
	refs := api.Group("/achievement-references")

	refs.Get("/", service.GetAllAchievementReferences)
	refs.Get("/:id", service.GetAchievementReferenceByID)
	refs.Post("/", service.CreateAchievementReference)
	refs.Put("/:id", service.UpdateAchievementReference)
	refs.Delete("/:id", service.DeleteAchievementReference)
}