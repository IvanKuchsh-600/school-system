package usecases

import (
	"auth-service/internal/entities"
	"auth-service/internal/ports"
	"errors"
	"fmt"
)

type AuthUseCase struct {
	userRepo ports.UserRepository
	jwt      ports.JWTManager
	hasher   ports.PasswordHasher
}

func NewAuthUseCase(
	userRepo ports.UserRepository,
	jwt ports.JWTManager,
	hasher ports.PasswordHasher,
) *AuthUseCase {
	return &AuthUseCase{
		userRepo: userRepo,
		jwt:      jwt,
		hasher:   hasher,
	}
}

func (s *AuthUseCase) Register(email, plainPassword, role string) (string, error) {
	err := entities.ValidateEmail(email)
	if err != nil {
		return "", fmt.Errorf("email validation failed: %w", err)
	}

	err = entities.ValidatePassword(plainPassword)
	if err != nil {
		return "", fmt.Errorf("password validation failed: %w", err)
	}

	err = entities.ValidateRole(role)
	if err != nil {
		return "", fmt.Errorf("role validation failed: %w", err)
	}

	existing, err := s.userRepo.FindByEmail(email)
	if err != nil && errors.Is(err, entities.ErrDatabaseOperation) {
		return "", fmt.Errorf("check user existence for %s: %w", email, err)
	}
	if existing != nil {
		return "", entities.ErrUserAlreadyExists
	}

	hashedPassword, err := s.hasher.Hash(plainPassword)
	if err != nil {
		return "", fmt.Errorf("hash password for user %s: %w", email, err)
	}

	user, err := entities.NewUser(email, hashedPassword, role)
	if err != nil {
		return "", fmt.Errorf("create user entity for %s: %w", email, err)
	}

	if err := s.userRepo.Create(user); err != nil {
		return "", fmt.Errorf("save user %s to database: %w", email, err)
	}

	token, err := s.jwt.Generate(user.ID, user.Role)
	if err != nil {
		return "", fmt.Errorf("generate token for user %d: %w", user.ID, err)
	}

	return token, nil
}

func (s *AuthUseCase) Login(email, plainPassword string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil || user == nil {
		return "", errors.New("invalid credentials")
	}

	if !s.hasher.Verify(user.PasswordHash, plainPassword) {
		return "", errors.New("invalid credentials")
	}

	token, err := s.jwt.Generate(user.ID, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthUseCase) ValidateToken(tokenString string) (*ports.JWTClaims, error) {
	return s.jwt.Verify(tokenString)
}
