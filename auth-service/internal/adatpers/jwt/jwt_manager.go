package jwt

import (
	"auth-service/internal/entities"
	"auth-service/internal/ports"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

func NewJWTManager(secretKey string, durationHours int) (*JWTManager, error) {
	if secretKey == "" {
		return nil, fmt.Errorf("JWT secret key is required")
	}

	if durationHours <= 0 {
		return nil, fmt.Errorf("token duration must be positive, got %d", durationHours)
	}

	return &JWTManager{
		secretKey:     secretKey,
		tokenDuration: time.Duration(durationHours) * time.Hour,
	}, nil
}

type jwtClaims struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

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

	signedToken, err := token.SignedString([]byte(m.secretKey))
	if err != nil {
		return "", entities.ErrInternalError
	}
	return signedToken, nil
}

func (m *JWTManager) Verify(tokenString string) (*ports.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.secretKey), nil
	})

	if err != nil {
		return nil, entities.ErrInvalidToken
	}

	if !token.Valid {
		return nil, entities.ErrInvalidToken
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok {
		return nil, entities.ErrInvalidToken
	}

	return &ports.JWTClaims{
		UserID: claims.UserID,
		Role:   claims.Role,
	}, nil
}
