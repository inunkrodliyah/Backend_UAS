package route

import "github.com/gofiber/fiber/v2"

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1") 


	SetupRoleRoutes(api)
	SetupPermissionRoutes(api)
	SetupRolePermissionRoutes(api)
	SetupUserRoutes(api)
	SetupStudentRoutes(api)
	SetupLecturerRoutes(api)
	SetupAchievementReferenceRoutes(api)

}