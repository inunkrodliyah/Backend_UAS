package helper

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword membuat hash dari password (SUDAH ADA SEBELUMNYA)
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

// CheckPasswordHash membandingkan password input dengan hash database (BARU)
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateToken membuat JWT Token (BARU)
func GenerateToken(userID string, roleID string) (string, error) {
	// Buat Claims (Isi token)
	claims := jwt.MapClaims{
		"user_id": userID,
		"role_id": roleID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // Token berlaku 72 jam
	}

	// Buat token dengan algoritma HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token dengan Secret Key (Harusnya dari .env, tapi hardcode dulu buat belajar)
	// Ganti "RAHASIA_SUPER" dengan string acak yang panjang
	return token.SignedString([]byte("RAHASIA_SUPER")) 
}