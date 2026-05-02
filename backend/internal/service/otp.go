package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
)

const otpTTL = 5 * time.Minute

type OTPService struct{ rdb *redis.Client }

func NewOTPService(rdb *redis.Client) *OTPService { return &OTPService{rdb} }

func (s *OTPService) Generate(ctx context.Context, phone, role string) (string, error) {
	code := fmt.Sprintf("%06d", rand.Intn(1_000_000))
	key := fmt.Sprintf("otp:%s:%s", role, phone)
	return code, s.rdb.Set(ctx, key, code, otpTTL).Err()
}

func (s *OTPService) Verify(ctx context.Context, phone, role, code string) (bool, error) {
	key := fmt.Sprintf("otp:%s:%s", role, phone)
	stored, err := s.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if stored != code {
		return false, nil
	}
	s.rdb.Del(ctx, key)
	return true, nil
}
