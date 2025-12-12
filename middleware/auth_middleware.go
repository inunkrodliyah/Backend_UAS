package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// AuthProtected adalah middleware untuk memproteksi route
func AuthProtected(c *fiber.Ctx) error {
	// 1. Ambil Header Authorization
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Token tidak ditemukan (Unauthorized)",
		})
	}

	// 2. Format harus "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Format token salah (Gunakan: Bearer <token>)",
		})
	}

	tokenString := parts[1]

	// 3. Parse dan Validasi Token
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// Pastikan algoritma signing sesuai
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.ErrUnauthorized
		}
		// PENTING: Secret Key harus SAMA PERSIS dengan yang di helper/password.go
		return []byte("RAHASIA_SUPER"), nil 
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Token tidak valid atau kadaluwarsa",
		})
	}

	// 4. Ambil Data Claims (User ID & Role ID) dari Token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Gagal membaca data token",
		})
	}

	// 5. Simpan data user ke Context (c.Locals) agar bisa dipakai di Controller/Service nanti
	c.Locals("user_id", claims["user_id"])
	c.Locals("role_id", claims["role_id"])

	// 6. Lanjut ke proses berikutnya
	return c.Next()
}