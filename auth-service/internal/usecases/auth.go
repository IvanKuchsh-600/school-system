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
) (*AuthUseCase, error) {
	if userRepo == nil {
		return nil, fmt.Errorf("userRepository is required")
	}
	if jwt == nil {
		return nil, fmt.Errorf("jwtManager is required")
	}
	if hasher == nil {
		return nil, fmt.Errorf("passwordHasher is required")
	}

	return &AuthUseCase{
		userRepo: userRepo,
		jwt:      jwt,
		hasher:   hasher,
	}, nil
}

func (s *AuthUseCase) Register(email, password, role string) (string, error) {
	err := entities.ValidateEmail(email)
	if err != nil {
		return "", fmt.Errorf("validate email '%s': %w", email, err)
	}

	err = entities.ValidatePassword(password)
	if err != nil {
		return "", fmt.Errorf("validate password for '%s': %w", email, err)
	}

	err = entities.ValidateRole(role)
	if err != nil {
		return "", fmt.Errorf("validate role '%s' for user '%s': %w", role, email, err)
	}

	existing, err := s.userRepo.GetByEmail(email)
	if err != nil && errors.Is(err, entities.ErrDatabaseOperation) {
		return "", fmt.Errorf("check user existence for %s: %w", email, err)
	}
	if existing != nil {
		return "", fmt.Errorf("user with email '%s' already exists: %w", email, entities.ErrUserAlreadyExists)
	}

	hashedPassword, err := s.hasher.Hash(password)
	if err != nil {
		return "", fmt.Errorf("hash password for user '%s': %w", email, err)
	}

	user, err := entities.NewUser(email, hashedPassword, role)
	if err != nil {
		return "", fmt.Errorf("create user entity for '%s': %w", email, err)
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return "", fmt.Errorf("save user '%s' to database: %w", email, entities.ErrDatabaseOperation)
	}

	token, err := s.jwt.Generate(user.ID, user.Role)
	if err != nil {
		return "", fmt.Errorf("generate token for user '%s' (id=%d): %w", email, user.ID, err)
	}

	return token, nil
}

func (s *AuthUseCase) Login(email, password string) (string, error) {
	err := entities.ValidateEmail(email)
	if err != nil {
		return "", fmt.Errorf("validate email '%s': %w", email, err)
	}

	err = entities.ValidatePassword(password)
	if err != nil {
		return "", fmt.Errorf("validate password for '%s': %w", email, err)
	}

	user, err := s.userRepo.GetByEmail(email)
	if err != nil && errors.Is(err, entities.ErrDatabaseOperation) {
		return "", fmt.Errorf("check user existence for %s: %w", email, err)
	}
	if user == nil {
		return "", fmt.Errorf("user with email '%s' not found: %w", email, entities.ErrInvalidCredentials)
	}

	err = s.hasher.Verify(user.PasswordHash, password)
	if err != nil {
		if errors.Is(err, entities.ErrInvalidCredentials) {
			return "", entities.ErrInvalidCredentials
		}
		return "", fmt.Errorf("bcrypt compare failed: %w", entities.ErrInternalError)
	}

	token, err := s.jwt.Generate(user.ID, user.Role)
	if err != nil {
		return "", fmt.Errorf("generate token for user %d: %w", user.ID, err)
	}

	return token, nil
}
