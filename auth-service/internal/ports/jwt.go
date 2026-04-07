package ports

// JWTManager - порт для работы с JWT токенами
type JWTManager interface {
	Generate(userID int64, role string) (string, error)
	Verify(tokenString string) (*JWTClaims, error)
}

// JWTClaims - данные, которые хранятся в токене
type JWTClaims struct {
	UserID int64
	Role   string
}
