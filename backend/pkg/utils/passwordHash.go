package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a given password using bcrypt and returns the hashed password.
func HashPassword(password string) (string, error) {
	// Convert the password string to a byte slice
	passwordBytes := []byte(password)
	// Generate the bcrypt hash from the password with a cost of bcrypt.DefaultCost
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	// Convert the hashed password byte slice back to a string and return it
	return string(hashedPasswordBytes), nil
}

// CheckPasswordHash checks if the given plain-text password matches the hashed password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
