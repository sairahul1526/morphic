package auth

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashPassword hashes a password using pure SHA256
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// CheckPasswordHash compares a hashed password with its plaintext version
func CheckPasswordHash(password, hashedPassword string) bool {
	hash := HashPassword(password)
	return hash == hashedPassword
}
