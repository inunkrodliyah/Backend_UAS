package route

import "github.com/gofiber/fiber/v2"

// SetupRoutes mendaftarkan semua rute dari file-file terpisah
func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1") // Grup prefix /api/v1

	// Daftarkan semua rute
	SetupRoleRoutes(api)
	SetupPermissionRoutes(api)
	SetupRolePermissionRoutes(api)
	SetupUserRoutes(api)
	SetupStudentRoutes(api)
	SetupLecturerRoutes(api)
	SetupAchievementReferenceRoutes(api)

	// Tambahkan rute untuk Autentikasi (Login, Register, dll) di sini nanti
	// SetupAuthRoutes(api)
}