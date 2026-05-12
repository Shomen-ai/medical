// Файл: internal/handler/auth.go
// Назначение: HTTP-обработчики аутентификации — отправка и проверка OTP-кодов для пациентов и сотрудников, логин персонала по паролю и обновление пары JWT-токенов.
package handler

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"beautymed/internal/config"
	"beautymed/internal/service"

	"github.com/gin-gonic/gin"
)

// AuthHandler — обработчик запросов аутентификации и обновления токенов.
type AuthHandler struct {
	svc     *service.Services
	devMode bool
}

// NewAuthHandler создаёт AuthHandler; в нелайв-окружениях включает режим отладки (возврат OTP в ответе).
func NewAuthHandler(svc *service.Services, cfg *config.Config) *AuthHandler {
	return &AuthHandler{svc: svc, devMode: cfg.Env != "production"}
}

// SendOTP отправляет одноразовый код пациенту по номеру телефона.
// POST /api/auth/otp  body: {"phone":"+79001234567"}
func (h *AuthHandler) SendOTP(c *gin.Context) {
	var req struct{ Phone string `json:"phone" binding:"required"` }
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phone required"})
		return
	}
	code, err := h.svc.Auth.SendPatientOTP(c.Request.Context(), req.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := gin.H{"message": "OTP sent"}
	if h.devMode && c.GetHeader("X-Dev-Mode") == "1" {
		resp["code"] = code
	}
	c.JSON(http.StatusOK, resp)
}

// VerifyOTP проверяет OTP пациента, выдаёт пару JWT и ставит refresh-токен в httpOnly cookie.
// POST /api/auth/verify  body: {"phone":"+79001234567","code":"123456"}
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required"`
		Code  string `json:"code"  binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pair, _, err := h.svc.Auth.VerifyPatientOTP(c.Request.Context(), req.Phone, req.Code)
	if err != nil {
		if errors.Is(err, service.ErrInvalidOTP) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired OTP"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.setRefreshCookie(c, pair.RefreshToken)
	c.JSON(http.StatusOK, gin.H{"access_token": pair.AccessToken})
}

// SendStaffOTP отправляет одноразовый код сотруднику (врачу или администратору) по номеру телефона.
// POST /api/staff/auth/otp
func (h *AuthHandler) SendStaffOTP(c *gin.Context) {
	var req struct{ Phone string `json:"phone" binding:"required"` }
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phone required"})
		return
	}
	code, err := h.svc.Auth.SendStaffOTP(c.Request.Context(), req.Phone)
	if err != nil {
		if strings.Contains(err.Error(), "staff not found") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "staff not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := gin.H{"message": "OTP sent"}
	if h.devMode && c.GetHeader("X-Dev-Mode") == "1" {
		resp["code"] = code
	}
	c.JSON(http.StatusOK, resp)
}

// VerifyStaffOTP проверяет OTP сотрудника и выдаёт JWT-токены с указанием его роли.
// POST /api/staff/auth/verify
func (h *AuthHandler) VerifyStaffOTP(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required"`
		Code  string `json:"code"  binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pair, _, role, err := h.svc.Auth.VerifyStaffOTP(c.Request.Context(), req.Phone, req.Code)
	if err != nil {
		if errors.Is(err, service.ErrInvalidOTP) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired OTP"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.setRefreshCookie(c, pair.RefreshToken)
	c.JSON(http.StatusOK, gin.H{"access_token": pair.AccessToken, "role": role})
}

// StaffLogin аутентифицирует сотрудника по логину и паролю и выдаёт пару JWT-токенов.
// POST /api/staff/auth/login  body: {"username":"admin","password":"..."}
func (h *AuthHandler) StaffLogin(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pair, _, role, err := h.svc.Auth.StaffLogin(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "неверный логин или пароль"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.setRefreshCookie(c, pair.RefreshToken)
	c.JSON(http.StatusOK, gin.H{"access_token": pair.AccessToken, "role": role})
}

// Refresh обновляет пару access/refresh токенов по refresh-токену из cookie.
// POST /api/auth/refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	refresh, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no refresh token"})
		return
	}
	pair, err := h.svc.Auth.RefreshTokens(c.Request.Context(), refresh)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	h.setRefreshCookie(c, pair.RefreshToken)
	c.JSON(http.StatusOK, gin.H{"access_token": pair.AccessToken})
}

func (h *AuthHandler) setRefreshCookie(c *gin.Context, token string) {
	secure := !h.devMode
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("refresh_token", token, int(30*24*time.Hour/time.Second),
		"/api/auth", "", secure, true)
}
