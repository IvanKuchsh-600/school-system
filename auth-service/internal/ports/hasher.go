package ports

// PasswordHasher - порт для хеширования и проверки паролей
type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(hashedPassword, plainPassword string) error
}
