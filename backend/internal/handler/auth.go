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

type AuthHandler struct {
	svc     *service.Services
	devMode bool
}

func NewAuthHandler(svc *service.Services, cfg *config.Config) *AuthHandler {
	return &AuthHandler{svc: svc, devMode: cfg.Env != "production"}
}

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
