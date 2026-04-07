package entities

import (
	"regexp"
)

func ValidateEmail(email string) error {
	if email == "" {
		return ErrEmailRequired
	}

	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return ErrEmailInvalid
	}

	return nil
}

func ValidatePassword(password string) error {
	if password == "" {
		return ErrPasswordRequired
	}

	return nil
}

func ValidateRole(role string) error {
	if role == "" {
		return ErrRoleRequired
	}

	validRoles := map[string]bool{"admin": true, "teacher": true, "parent": true, "student": true}

	if !validRoles[role] {
		return ErrRoleInvalid
	}

	return nil
}

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	Role         string // admin, teacher, parent, student
}

// NewUser создает нового пользователя с валидацией
func NewUser(email, passwordHash, role string) (*User, error) {

	err := ValidateEmail(email)
	if err != nil {
		return nil, err
	}

	err = ValidatePassword(passwordHash)
	if err != nil {
		return nil, err
	}

	err = ValidateRole(role)
	if err != nil {
		return nil, err
	}

	return &User{
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
	}, nil
}
