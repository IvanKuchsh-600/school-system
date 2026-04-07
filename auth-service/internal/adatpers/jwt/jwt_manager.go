package jwt

import (
	"auth-service/internal/ports"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

// NewJWTManager - конструктор адаптера
func NewJWTManager(secretKey string, durationHours int) *JWTManager {
	return &JWTManager{
		secretKey:     secretKey,
		tokenDuration: time.Duration(durationHours) * time.Hour,
	}
}

// JWTClaims - внутренняя структура для jwt
type jwtClaims struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Generate реализует ports.JWTManager
func (m *JWTManager) Generate(userID int64, role string) (string, error) {
	claims := &jwtClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}

// Verify реализует ports.JWTManager
func (m *JWTManager) Verify(tokenString string) (*ports.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok {
		return nil, errors.New("cannot parse claims")
	}

	return &ports.JWTClaims{
		UserID: claims.UserID,
		Role:   claims.Role,
	}, nil
}
