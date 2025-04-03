package jwt_test

import (
	"testing"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kodenkai-labs/go-lib/jwt"
)

type claims struct {
	Data string `json:"data"`
	jwtv5.RegisteredClaims
}

func newClaims(ttl time.Duration, data string) *claims {
	return &claims{
		Data: data,
		RegisteredClaims: jwtv5.RegisteredClaims{
			ExpiresAt: jwtv5.NewNumericDate(time.Now().Add(ttl)),
		},
	}
}

func Test_GenerateToken(t *testing.T) {
	expiresAt := time.Date(2023, 10, 10, 0, 0, 0, 0, time.UTC)

	type args struct {
		secret string
		claims *claims
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Test #1: Success",
			args: args{
				secret: "secret",
				claims: &claims{
					Data: "some_data",
					RegisteredClaims: jwtv5.RegisteredClaims{
						Issuer:    "app_name",
						ExpiresAt: jwtv5.NewNumericDate(expiresAt),
					},
				},
			},
			want:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjoic29tZV9kYXRhIiwiaXNzIjoiYXBwX25hbWUiLCJleHAiOjE2OTY4OTYwMDB9.dDCKy4TwxN5OvtvJN49_Nwi26ORafc9k2d-Ky1PvQB4", //nolint:lll
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jwt.GenerateToken(tt.args.secret, tt.args.claims)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_ParseTokenWithClaims(t *testing.T) {
	validClaims := newClaims(time.Hour, "some_data")
	invalidClaims := newClaims(-time.Hour, "some_data2")

	type args struct {
		secret string
		claims *claims
	}
	tests := []struct {
		name    string
		args    args
		want    *claims
		wantErr bool
		err     error
	}{
		{
			name: "Test #1: Success",
			args: args{
				secret: "secret",
				claims: validClaims,
			},
			want:    validClaims,
			wantErr: false,
		},
		{
			name: "Test #2: Expired Token",
			args: args{
				secret: "secret",
				claims: invalidClaims,
			},
			want:    invalidClaims,
			wantErr: true,
			err:     jwt.ErrTokenExpired,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := jwt.GenerateToken(tt.args.secret, tt.args.claims)
			require.NoError(t, err)

			claims := &claims{}
			err = jwt.ParseTokenWithClaims(tt.args.secret, token, claims)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want.Data, claims.Data)
		})
	}
}
