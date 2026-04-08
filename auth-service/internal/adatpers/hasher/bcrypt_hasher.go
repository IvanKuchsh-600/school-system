package hasher

import (
	"auth-service/internal/entities"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct{}

func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{}
}

func (h *BcryptHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", entities.ErrHashFailed
	}
	return string(hash), nil
}

func (h *BcryptHasher) Verify(hashedPassword, plainPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))

	if err == nil {
		return nil
	}

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return entities.ErrInvalidCredentials
	}

	return entities.ErrInternalError
}
