package helper

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword generates bcrypt hash
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

// CheckPasswordHash compares password with hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateToken creates JWT with Role AND Permissions (FR-001)
func GenerateToken(userID string, roleID string, permissions []string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":     userID,
		"role_id":     roleID,
		"permissions": permissions, // Menyimpan permissions di token
		"exp":         time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Ganti "RAHASIA_SUPER" dengan env variable di production
	return token.SignedString([]byte("RAHASIA_SUPER"))
}