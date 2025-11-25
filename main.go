package main

import (
	"fmt" // 1. Import 'fmt' untuk memformat string port
	"log"
	"os"
	// Pastikan nama modul ini (project-uas) sesuai dengan file go.mod Anda
	"project-uas/database"
	"project-uas/route"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Could not load .env file")
	}

	// Connect to database
	database.ConnectDB() // Ini akan otomatis membaca DB_DSN dari .env

	// Setup Fiber app
	app := fiber.New()

	// Add logger middleware
	app.Use(logger.New())

	// Setup all routes
	route.SetupRoutes(app)

	// 2. Get port from env or default
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000" // Default jika APP_PORT tidak disetel
	}

	// 3. Buat alamat 'Listen' yang valid (cth: ":1063")
	listenAddr := fmt.Sprintf(":%s", port)

	log.Printf("Starting server on port %s\n", listenAddr)
	// 4. Gunakan listenAddr yang sudah diformat
	log.Fatal(app.Listen(listenAddr))
}