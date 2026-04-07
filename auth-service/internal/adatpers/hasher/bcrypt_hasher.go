package hasher

import (
	"golang.org/x/crypto/bcrypt"
)

// BcryptHasher - адаптер для порта PasswordHasher
type BcryptHasher struct{}

// NewBcryptHasher - конструктор адаптера
func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{}
}

// Hash реализует ports.PasswordHasher
func (h *BcryptHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Verify реализует ports.PasswordHasher
func (h *BcryptHasher) Verify(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
