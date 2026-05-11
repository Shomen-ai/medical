package service

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidOTP = errors.New("invalid or expired OTP")

type AuthService struct {
	repos *Repos
	otp   *OTPService
	jwt   *JWTService
}

func NewAuthService(repos *Repos, otp *OTPService, jwt *JWTService) *AuthService {
	return &AuthService{repos, otp, jwt}
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *AuthService) SendPatientOTP(ctx context.Context, phone string) (string, error) {
	code, err := s.otp.Generate(ctx, phone, "patient")
	if err != nil {
		return "", fmt.Errorf("generate otp: %w", err)
	}
	// TODO Plan 2: send via SMSC
	return code, nil
}

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

var ErrInvalidCredentials = errors.New("invalid username or password")

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
