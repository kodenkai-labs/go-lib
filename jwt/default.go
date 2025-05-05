package jwt

import (
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

type DefaultClaims struct {
	Data string `json:"data"`

	jwtv5.RegisteredClaims
}

// NewClaims creates a new Claims object with the user-specific data.
func NewDefaultClaims(ttl time.Duration, data string) *DefaultClaims {
	return &DefaultClaims{
		Data: data,
		RegisteredClaims: jwtv5.RegisteredClaims{
			ExpiresAt: jwtv5.NewNumericDate(time.Now().Add(ttl)),
		},
	}
}
