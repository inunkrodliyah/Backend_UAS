package route

import "github.com/gofiber/fiber/v2"

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1") 


	SetupAuthRoutes(api)   // <-- Daftarkan Auth
	SetupRoleRoutes(api)
	SetupPermissionRoutes(api)
	SetupRolePermissionRoutes(api)
	SetupUserRoutes(api)
	SetupStudentRoutes(api)
	SetupLecturerRoutes(api)
	SetupAchievementRoutes(api)
	SetupReportRoutes(api)

}