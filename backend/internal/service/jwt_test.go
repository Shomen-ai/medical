package service_test

import (
	"testing"

	"beautymed/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTService_IssueAndParse(t *testing.T) {
	svc := service.NewJWTService("test-secret-32-chars-minimum-ok!")

	token, err := svc.IssueAccess("user-123", "patient")
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := svc.Parse(token)
	require.NoError(t, err)
	assert.Equal(t, "user-123", claims.UserID)
	assert.Equal(t, "patient", claims.Role)
	assert.Equal(t, "access", claims.TokenType)
}

func TestJWTService_InvalidToken(t *testing.T) {
	svc := service.NewJWTService("test-secret-32-chars-minimum-ok!")
	_, err := svc.Parse("not.a.token")
	assert.Error(t, err)
}

func TestJWTService_ParseRefresh_Success(t *testing.T) {
	svc := service.NewJWTService("test-secret-32-chars-minimum-ok!")

	token, err := svc.IssueRefresh("user-456", "staff")
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := svc.ParseRefresh(token)
	require.NoError(t, err)
	assert.Equal(t, "user-456", claims.UserID)
	assert.Equal(t, "staff", claims.Role)
	assert.Equal(t, "refresh", claims.TokenType)
}

func TestJWTService_ParseRefresh_RejectsAccessToken(t *testing.T) {
	svc := service.NewJWTService("test-secret-32-chars-minimum-ok!")

	accessToken, err := svc.IssueAccess("user-789", "patient")
	require.NoError(t, err)

	_, err = svc.ParseRefresh(accessToken)
	assert.Error(t, err)
	assert.EqualError(t, err, "not a refresh token")
}
