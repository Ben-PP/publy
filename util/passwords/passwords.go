package passwords

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func GenerateHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", fmt.Errorf("Failed to generate hash: %w", err)
	}

	return string(bytes), nil
}

func CompareToHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}
