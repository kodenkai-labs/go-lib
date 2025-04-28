package errlib

type Slug string

// Slug is an error slug to be used in the localization service. It can be extended in the application to include more specific error slugs.

// General
const (
	SlugInternal           Slug = "internal"
	SlugBadGateway         Slug = "bad_gateway"
	SlugInvalidBodyRequest Slug = "invalid_body_request"
)

// Sessions
const (
	SlugUserUnauthorized Slug = "user_unauthorized"

	SlugEmptyClientID Slug = "empty_client_id"

	SlugEmptyRefreshToken     Slug = "empty_refresh_token"
	SlugInvalidRefreshToken   Slug = "invalid_refresh_token"
	SlugRefreshTokenDifferent Slug = "refresh_token_different"

	SlugInvalidAccessToken Slug = "invalid_access_token"
)
