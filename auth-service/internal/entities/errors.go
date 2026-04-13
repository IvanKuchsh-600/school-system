package entities

import "errors"

// User errors
var (
	ErrUserAlreadyExists = errors.New("User already exists")
	ErrUserNotFound      = errors.New("User not found")
)

// Email errors
var (
	ErrEmailRequired      = errors.New("Email is required")
	ErrEmailInvalid       = errors.New("Invalid email format")
	ErrInvalidCredentials = errors.New("Invalid credentials")
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

// JWT errors
var (
	ErrTokenGeneration = errors.New("failed to generate token")
	ErrInvalidToken    = errors.New("invalid or expired token")
)

// Internal
var (
	ErrInternalError = errors.New("Internal error")
)

var (
	ErrInvalidParams = errors.New("Invalid params")
)
