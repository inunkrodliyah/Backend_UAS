package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// AuthProtected: Validasi Token & Ekstrak Data User
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
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.ErrUnauthorized
		}
		// PENTING: Secret Key harus SAMA dengan helper/password.go
		return []byte("RAHASIA_SUPER"), nil 
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Token tidak valid atau kadaluwarsa",
		})
	}

	// 4. Ambil Data Claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Gagal membaca data token",
		})
	}

	// 5. Simpan data user ke Context (Locals)
	c.Locals("user_id", claims["user_id"])
	c.Locals("role_id", claims["role_id"])

	// --- LOGIC TAMBAHAN UNTUK PERMISSION ---
	// Ambil permissions dari token (biasanya bentuknya []interface{})
	// Kita perlu konversi jadi []string agar mudah dicek
	var permissions []string
	if permsClaim, ok := claims["permissions"].([]interface{}); ok {
		for _, p := range permsClaim {
			if str, ok := p.(string); ok {
				permissions = append(permissions, str)
			}
		}
	}
	// Simpan list permission ke Locals
	c.Locals("permissions", permissions)

	return c.Next()
}

// RequirePermission: Middleware Cek Hak Akses (RBAC)
// Fungsi ini yang tadi error "undefined"
func RequirePermission(requiredPerm string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Ambil list permission user dari Locals (yang diset di AuthProtected)
		userPerms, ok := c.Locals("permissions").([]string)
		if !ok {
			// Jika tidak ada permission sama sekali
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "fail",
				"message": "Forbidden: No permissions found",
			})
		}

		// 2. Cek apakah user punya permission yang diminta
		hasPermission := false
		for _, p := range userPerms {
			if p == requiredPerm {
				hasPermission = true
				break
			}
		}

		// 3. Jika tidak punya, tolak request
		if !hasPermission {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "fail",
				"message": "Forbidden: You don't have permission '" + requiredPerm + "'",
			})
		}

		// 4. Jika punya, lanjut
		return c.Next()
	}
}