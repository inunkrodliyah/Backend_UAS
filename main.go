package main

import (
	"fmt"
	"log"
	"os"

	"project-uas/database"
	"project-uas/route"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"

	_ "project-uas/docs" // WAJIB: side-effect import Swagger docs
)

// @title           Sistem Pencatatan Prestasi Mahasiswa API
// @version         1.0
// @description     Dokumentasi API untuk Project UAS Backend Lanjutan
// @termsOfService  http://swagger.io/terms/

// @contact.name    Tim Pengembang
// @contact.email   admin@univ.ac.id

// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

// @host            localhost:1063
// @BasePath        /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Masukkan token dengan format: Bearer <your_token>
func main() {

	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not loaded")
	}

	// Connect Database
	database.ConnectDB()

	// Init Fiber
	app := fiber.New()
	app.Use(logger.New())

	// Swagger endpoint
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Setup API routes
	route.SetupRoutes(app)

	// Port
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "1063"
	}

	listenAddr := fmt.Sprintf(":%s", port)
	log.Printf("Server running on %s\n", listenAddr)

	log.Fatal(app.Listen(listenAddr))
}
