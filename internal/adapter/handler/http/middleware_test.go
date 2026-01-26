package handler // Ensure this matches the package name in middleware.go (e.g., 'http' or 'handler')

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateToken(t *testing.T) {
	email := "test_user@example.com"

	tokenString, err := GenerateToken(email)

	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	if tokenString == "" {
		t.Error("Expected a token string, got empty string")
	}

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil // We need access to the secret key from middleware.go
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Does the email inside the token match?
		if claims["email"] != email {
			t.Errorf("Expected email %v, got %v", email, claims["email"])
		}
	} else {
		t.Error("Token is invalid")
	}
}
