package main

import (
	
	"log"
	"os"

	"project-uas/database"
	"project-uas/route"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"

	_ "project-uas/docs"
)

// @title           Sistem Pencatatan Prestasi Mahasiswa API
// @version         1.0
// @description     Dokumentasi API untuk Project UAS Backend Lanjutan
// @host            localhost:1063
// @BasePath        /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {

	// Load env
	_ = godotenv.Load()

	// DB
	database.ConnectDB()

	// Fiber
	app := fiber.New()

	// ✅ CORS (WAJIB SEBELUM ROUTE)
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Logger
	app.Use(logger.New())

	// Swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Routes
	route.SetupRoutes(app)

	// Port
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "1063"
	}

	log.Fatal(app.Listen(":" + port))
}
