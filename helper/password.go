package helper

import "golang.org/x/crypto/bcrypt"

// HashPassword menghasilkan hash bcrypt dari password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}
