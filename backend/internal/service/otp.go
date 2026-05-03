package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/redis/go-redis/v9"
)

const otpTTL = 5 * time.Minute

type OTPService struct{ rdb *redis.Client }

func NewOTPService(rdb *redis.Client) *OTPService { return &OTPService{rdb} }

func (s *OTPService) Generate(ctx context.Context, phone, role string) (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
	if err != nil {
		return "", fmt.Errorf("generate otp: %w", err)
	}
	code := fmt.Sprintf("%06d", n.Int64())
	key := fmt.Sprintf("otp:%s:%s", role, phone)
	return code, s.rdb.Set(ctx, key, code, otpTTL).Err()
}

func (s *OTPService) Verify(ctx context.Context, phone, role, code string) (bool, error) {
	key := fmt.Sprintf("otp:%s:%s", role, phone)
	stored, err := s.rdb.GetDel(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return stored == code, nil
}
