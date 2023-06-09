// Package hasher provide functions for coverting
// given plain-text password into a hash
// and validating plain-text password against given hash
package hasher

import "golang.org/x/crypto/bcrypt"

// Receives string password as input and returns its string hash or error if any problems occur
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Receives string password and hash and returns true if hashed password is equal to provided hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
