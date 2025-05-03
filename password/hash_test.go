package password_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kodenkai-labs/go-lib/password"
)

func Test_HashPassword(t *testing.T) {
	rawPassword := "password"
	wrongRawPassword := "wrong_password"

	hashedPassword, err := password.HashPassword(rawPassword)
	require.NoError(t, err)

	isValid := password.CheckPassword(hashedPassword, wrongRawPassword)
	assert.False(t, isValid)

	isValid = password.CheckPassword(hashedPassword, rawPassword)
	assert.True(t, isValid)
}
