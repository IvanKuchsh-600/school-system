package entities

import "errors"

// User errors
var (
	ErrUserAlreadyExists = errors.New("user already exists")
)

// Email errors
var (
	ErrEmailRequired = errors.New("Email is required")
	ErrEmailInvalid  = errors.New("Invalid email format")
)

// Password errors
var (
	ErrPasswordRequired = errors.New("Password is required")
	ErrHashFailed       = errors.New("Failed to hash password")
)

// Role errors
var (
	ErrRoleRequired = errors.New("Role is required")
	ErrRoleInvalid  = errors.New("Invalid role")
)

// Database errors
var (
	ErrDatabaseOperation = errors.New("Database operation failed")
)
