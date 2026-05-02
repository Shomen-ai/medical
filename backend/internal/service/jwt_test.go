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
}

func TestJWTService_InvalidToken(t *testing.T) {
	svc := service.NewJWTService("test-secret-32-chars-minimum-ok!")
	_, err := svc.Parse("not.a.token")
	assert.Error(t, err)
}
