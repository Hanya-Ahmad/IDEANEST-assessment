package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword encrypts a password using bcrypt's hashing algorithm
func HashPassword(password []byte) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}

// CheckPasswordHash checks if a given password matches the hashed one
func CheckPasswordHash(hashedPassword []byte, password string) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	return err == nil
}
