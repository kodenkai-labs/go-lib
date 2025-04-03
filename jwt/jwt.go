package jwt

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken            = errors.New("invalid token")
	ErrTokenExpired            = errors.New("token expired")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
)

// GenerateToken generates a JWT token based on the given claims.
func GenerateToken(secretKey string, claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("get signed string: %w", err)
	}

	return tokenString, nil
}

// ParseToken parses and validates the token string, returning the claims if valid.
func ParseTokenWithClaims(secretKey, tokenString string, claims jwt.Claims) error {
	// Parse the token with the secret key and validate it
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token's signing method matches the expected method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
		}

		return []byte(secretKey), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return ErrTokenExpired
		}

		return fmt.Errorf("parse with claims: %w", err)
	}

	if token.Valid {
		return nil
	}

	return ErrInvalidToken
}
