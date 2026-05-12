// Файл: internal/service/jwt.go
// Назначение: выпуск и валидация JWT-токенов (access на 15 минут и refresh на 30 дней) для API.
package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims — полезная нагрузка JWT с идентификатором пользователя, ролью и типом токена.
type Claims struct {
	UserID    string `json:"user_id"`
	Role      string `json:"role"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

// JWTService подписывает и проверяет JWT-токены симметричным ключом HS256.
type JWTService struct{ secret []byte }

// NewJWTService создаёт сервис JWT с заданным секретным ключом.
func NewJWTService(secret string) *JWTService { return &JWTService{[]byte(secret)} }

// IssueAccess выпускает access-токен (срок жизни 15 минут).
func (s *JWTService) IssueAccess(userID, role string) (string, error) {
	return s.issue(userID, role, "access", 15*time.Minute)
}

// IssueRefresh выпускает refresh-токен (срок жизни 30 дней).
func (s *JWTService) IssueRefresh(userID, role string) (string, error) {
	return s.issue(userID, role, "refresh", 30*24*time.Hour)
}

// Parse валидирует любой токен (access или refresh) и возвращает его claims.
func (s *JWTService) Parse(tokenStr string) (*Claims, error) {
	return s.parse(tokenStr)
}

// ParseRefresh валидирует токен и убеждается, что это именно refresh-токен.
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
