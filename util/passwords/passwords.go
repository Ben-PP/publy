package passwords

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// GenerateHash generates a bcrypt hash of the given password.
func GenerateHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", fmt.Errorf("Failed to generate hash: %w", err)
	}

	return string(bytes), nil
}

// CompareToHash compares a plaintext password to a bcrypt hash.
func CompareToHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}
