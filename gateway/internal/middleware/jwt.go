package middleware

import (
	"gateway/internal/client"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuth(authClient *client.AuthClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
			c.Abort()
			return
		}

		valid, role, userID := authClient.ValidateToken(parts[1])
		if !valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Сохраняем в контекст Gin
		c.Set("user_id", userID)
		c.Set("role", role)
		c.Next()
	}
}
