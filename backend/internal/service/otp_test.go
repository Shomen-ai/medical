package service_test

import (
	"context"
	"testing"

	"beautymed/internal/service"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestRedis(t *testing.T) *redis.Client {
	t.Helper()
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379", DB: 15})
	require.NoError(t, rdb.Ping(context.Background()).Err())
	t.Cleanup(func() { rdb.FlushDB(context.Background()) })
	return rdb
}

func TestOTPService_StoreAndVerify(t *testing.T) {
	rdb := newTestRedis(t)
	svc := service.NewOTPService(rdb)
	ctx := context.Background()

	code, err := svc.Generate(ctx, "+79001234567", "patient")
	require.NoError(t, err)
	assert.Len(t, code, 6)

	ok, err := svc.Verify(ctx, "+79001234567", "patient", code)
	require.NoError(t, err)
	assert.True(t, ok)

	// Code is consumed — second verify fails
	ok2, err := svc.Verify(ctx, "+79001234567", "patient", code)
	require.NoError(t, err)
	assert.False(t, ok2)
}

func TestOTPService_WrongCode(t *testing.T) {
	rdb := newTestRedis(t)
	svc := service.NewOTPService(rdb)
	ctx := context.Background()

	_, err := svc.Generate(ctx, "+79009999999", "patient")
	require.NoError(t, err)

	ok, err := svc.Verify(ctx, "+79009999999", "patient", "000000")
	require.NoError(t, err)
	assert.False(t, ok)
}
