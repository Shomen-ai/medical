// Файл: internal/service/auth.go
// Назначение: бизнес-логика входа — OTP-аутентификация пациентов и сотрудников, логин по паролю, выпуск JWT-пар.
package service

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// ErrInvalidOTP сигнализирует, что введён неверный или истёкший OTP-код.
var ErrInvalidOTP = errors.New("invalid or expired OTP")

// AuthService — сервис аутентификации пациентов и сотрудников.
type AuthService struct {
	repos *Repos
	otp   *OTPService
	jwt   *JWTService
}

// NewAuthService собирает AuthService из репозиториев, OTP- и JWT-сервисов.
func NewAuthService(repos *Repos, otp *OTPService, jwt *JWTService) *AuthService {
	return &AuthService{repos, otp, jwt}
}

// TokenPair — пара access/refresh JWT-токенов, выдаваемая клиенту.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// SendPatientOTP генерирует и сохраняет OTP пациента (в DEV-режиме возвращает код для отображения).
func (s *AuthService) SendPatientOTP(ctx context.Context, phone string) (string, error) {
	code, err := s.otp.Generate(ctx, phone, "patient")
	if err != nil {
		return "", fmt.Errorf("generate otp: %w", err)
	}
	// TODO Plan 2: send via SMSC
	return code, nil
}

// VerifyPatientOTP проверяет OTP, создаёт/находит пациента и выдаёт пару JWT-токенов.
func (s *AuthService) VerifyPatientOTP(ctx context.Context, phone, code string) (*TokenPair, string, error) {
	ok, err := s.otp.Verify(ctx, phone, "patient", code)
	if err != nil {
		return nil, "", err
	}
	if !ok {
		return nil, "", ErrInvalidOTP
	}
	user, err := s.repos.Users.Create(phone)
	if err != nil {
		return nil, "", fmt.Errorf("create user: %w", err)
	}
	pair, err := s.issuePair(user.ID, "patient")
	return pair, user.ID, err
}

// SendStaffOTP отправляет OTP сотруднику, предварительно проверив, что такой staff существует.
func (s *AuthService) SendStaffOTP(ctx context.Context, phone string) (string, error) {
	if _, err := s.repos.Doctors.FindStaffByPhone(phone); err != nil {
		return "", fmt.Errorf("staff not found: %w", err)
	}
	code, err := s.otp.Generate(ctx, phone, "staff")
	if err != nil {
		return "", fmt.Errorf("generate otp: %w", err)
	}
	return code, nil
}

// VerifyStaffOTP проверяет OTP сотрудника и выдаёт JWT-пару с подменой user_id на doctor_id для роли doctor.
func (s *AuthService) VerifyStaffOTP(ctx context.Context, phone, code string) (*TokenPair, string, string, error) {
	ok, err := s.otp.Verify(ctx, phone, "staff", code)
	if err != nil {
		return nil, "", "", err
	}
	if !ok {
		return nil, "", "", ErrInvalidOTP
	}
	st, err := s.repos.Doctors.FindStaffByPhone(phone)
	if err != nil {
		return nil, "", "", fmt.Errorf("staff not found: %w", err)
	}
	// For doctors, JWT user_id must equal doctors.id (used in appointment queries)
	userID := st.ID
	if st.Role == "doctor" && st.DoctorID != nil {
		userID = *st.DoctorID
	}
	pair, err := s.issuePair(userID, st.Role)
	return pair, userID, st.Role, err
}

// ErrInvalidCredentials сигнализирует, что введены неверные логин или пароль сотрудника.
var ErrInvalidCredentials = errors.New("invalid username or password")

// StaffLogin авторизует сотрудника по логину и паролю (bcrypt) и выдаёт пару JWT.
// StaffLogin authenticates staff by username + plain-text password.
func (s *AuthService) StaffLogin(username, password string) (*TokenPair, string, string, error) {
	st, err := s.repos.Doctors.FindStaffByUsername(username)
	if err != nil {
		return nil, "", "", ErrInvalidCredentials
	}
	if st.PasswordHash == nil {
		return nil, "", "", ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(*st.PasswordHash), []byte(password)); err != nil {
		return nil, "", "", ErrInvalidCredentials
	}
	userID := st.ID
	if st.Role == "doctor" && st.DoctorID != nil {
		userID = *st.DoctorID
	}
	pair, err := s.issuePair(userID, st.Role)
	return pair, userID, st.Role, err
}

// RefreshTokens проверяет refresh-токен и выпускает новую пару access/refresh.
func (s *AuthService) RefreshTokens(ctx context.Context, refreshToken string) (*TokenPair, error) {
	claims, err := s.jwt.ParseRefresh(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}
	return s.issuePair(claims.UserID, claims.Role)
}

func (s *AuthService) issuePair(userID, role string) (*TokenPair, error) {
	access, err := s.jwt.IssueAccess(userID, role)
	if err != nil {
		return nil, err
	}
	refresh, err := s.jwt.IssueRefresh(userID, role)
	if err != nil {
		return nil, err
	}
	return &TokenPair{AccessToken: access, RefreshToken: refresh}, nil
}
