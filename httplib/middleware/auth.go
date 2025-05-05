package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/kodenkai-labs/go-lib/errlib"
	"github.com/kodenkai-labs/go-lib/httplib"
	"github.com/kodenkai-labs/go-lib/jwt"
)

// Authorization header
const (
	AuthorizationHeaderName = "Authorization"
	bearerPrefix            = "Bearer"
	tokenParts              = 2
)

// Cookies
const (
	CookieClientIDKey     = "client_id"
	CookieRefreshTokenKey = "refresh_token"
)

func AuthMiddleware(accessTokenSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.GetHeader(AuthorizationHeaderName)
		if bearerToken == "" {
			httplib.HandleError(c, errlib.NewAppError(
				nil, errlib.UnauthorizedCode, errlib.SlugUserUnauthorized))
			c.Abort()

			return
		}

		splitToken := strings.Split(bearerToken, " ")
		if len(splitToken) != tokenParts && strings.EqualFold(splitToken[0], bearerPrefix) {
			httplib.HandleError(c, errlib.NewAppError(
				nil, errlib.UnauthorizedCode, errlib.SlugInvalidAccessToken))
			c.Abort()

			return
		}

		claims := jwt.DefaultClaims{}
		err := jwt.ParseTokenWithClaims(splitToken[1], accessTokenSecret, &claims)
		if err != nil {
			httplib.HandleError(c, errlib.NewAppError(err, errlib.UnauthorizedCode, errlib.SlugInvalidAccessToken))
			c.Abort()

			return
		}

		clientID, err := getClientIDFromCookie(c)
		if err != nil && !errors.Is(err, http.ErrNoCookie) {
			httplib.HandleError(c, errlib.NewAppError(err, errlib.UnauthorizedCode, errlib.SlugEmptyClientID))
			c.Abort()

			return
		}

		refreshToken, err := getRefreshTokenFromCookie(c)
		if err != nil && !errors.Is(err, http.ErrNoCookie) {
			httplib.HandleError(c, errlib.NewAppError(err, errlib.UnauthorizedCode, errlib.SlugEmptyRefreshToken))
			c.Abort()

			return
		}

		c.Set(httplib.SessionDataKey, claims.Data)
		c.Set(httplib.ClientIDKey, clientID)
		c.Set(httplib.RefreshTokenKey, refreshToken)

		c.Next()
	}
}

func CookiesMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID, err := getClientIDFromCookie(c)
		if err != nil && !errors.Is(err, http.ErrNoCookie) {
			httplib.HandleError(c, errlib.NewAppError(err, errlib.InvalidInputCode, errlib.SlugEmptyClientID))
			c.Abort()

			return
		}

		refreshToken, err := getRefreshTokenFromCookie(c)
		if err != nil && !errors.Is(err, http.ErrNoCookie) {
			httplib.HandleError(c, errlib.NewAppError(err, errlib.InvalidInputCode, errlib.SlugEmptyRefreshToken))
			c.Abort()

			return
		}

		c.Set(httplib.ClientIDKey, clientID)
		c.Set(httplib.RefreshTokenKey, refreshToken)

		c.Next()
	}
}

func getClientIDFromCookie(c *gin.Context) (string, error) {
	clientIDCookie, err := c.Request.Cookie(CookieClientIDKey)
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		return "", fmt.Errorf("get client id cookie: %w", err)
	}

	var clientID string
	if clientIDCookie != nil {
		clientID = clientIDCookie.Value
	}

	return clientID, nil
}

func getRefreshTokenFromCookie(c *gin.Context) (string, error) {
	refreshTokenCookie, err := c.Request.Cookie(CookieRefreshTokenKey)
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		return "", fmt.Errorf("get refresh token cookie: %w", err)
	}

	var refreshToken string
	if refreshTokenCookie != nil {
		refreshToken = refreshTokenCookie.Value
	}

	return refreshToken, nil
}
