package main

import (
	"testing"

	"github.com/gofiber/fiber/v2"
)

/*
   Test ini memastikan aplikasi bisa dibuat
   tanpa panic / error (basic sanity test)
*/

func TestAppInitialization(t *testing.T) {
	app := fiber.New()

	if app == nil {
		t.Fatalf("fiber app should not be nil")
	}
}
