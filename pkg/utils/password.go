package utils

import "golang.org/x/crypto/bcrypt"

const passwordCost = 14

func HashPassword(pass string) string {
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(pass), passwordCost)

	return string(hashedPass)
}

func ValidatePassword(hashedPass, pass string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(pass))
}
