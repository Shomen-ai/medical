package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID    string `json:"user_id"`
	Role      string `json:"role"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

type JWTService struct{ secret []byte }

func NewJWTService(secret string) *JWTService { return &JWTService{[]byte(secret)} }

func (s *JWTService) IssueAccess(userID, role string) (string, error) {
	return s.issue(userID, role, "access", 15*time.Minute)
}

func (s *JWTService) IssueRefresh(userID, role string) (string, error) {
	return s.issue(userID, role, "refresh", 30*24*time.Hour)
}

func (s *JWTService) Parse(tokenStr string) (*Claims, error) {
	return s.parse(tokenStr)
}

func (s *JWTService) ParseRefresh(tokenStr string) (*Claims, error) {
	c, err := s.parse(tokenStr)
	if err != nil {
		return nil, err
	}
	if c.TokenType != "refresh" {
		return nil, errors.New("not a refresh token")
	}
	return c, nil
}

func (s *JWTService) issue(userID, role, tokenType string, ttl time.Duration) (string, error) {
	claims := Claims{
		UserID:    userID,
		Role:      role,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secret)
}

func (s *JWTService) parse(tokenStr string) (*Claims, error) {
	t, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, err
	}
	if c, ok := t.Claims.(*Claims); ok && t.Valid {
		return c, nil
	}
	return nil, errors.New("invalid token")
}
