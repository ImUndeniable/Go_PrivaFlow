package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Secret Key (In production, this comes from ENV variables!)
var jwtSecret = []byte("super_secret_key_123")

// 1. The Bouncer Function
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// A. Get the Token from the Header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No Authorization header found"})
			return
		}

		// Header format is usually "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Header format. Use 'Bearer <token>'"})
			return
		}

		tokenString := parts[1]

		// B. Validate the Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or Expired Token"})
			return
		}

		// C. If Valid, Let them pass!
		c.Next()
	}
}

// 2. A Helper to Generate Tokens (For testing)
func GenerateToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"role":  "admin",
	})
	return token.SignedString(jwtSecret)
}
