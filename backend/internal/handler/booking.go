package handler

import (
	"errors"
	"net/http"
	"time"

	"beautymed/internal/service"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct{ svc *service.Services }

func NewBookingHandler(svc *service.Services) *BookingHandler { return &BookingHandler{svc: svc} }

// GET /api/doctors/:id/slots?service_id=uuid&date=2026-05-10
func (h *BookingHandler) GetSlots(c *gin.Context) {
	slots, err := h.svc.Booking.GetSlots(
		c.Param("id"),
		c.Query("service_id"),
		c.Query("date"),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, slots)
}

// GET /api/doctors/:id/available-dates?month=2026-05
func (h *BookingHandler) GetAvailableDates(c *gin.Context) {
	dates, err := h.svc.Booking.GetAvailableDates(c.Param("id"), c.Query("month"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dates)
}

// POST /api/promo/check
// body: {"code":"SUMMER10","service_id":"uuid"}
func (h *BookingHandler) CheckPromo(c *gin.Context) {
	var req struct {
		Code      string `json:"code"`
		ServiceID string `json:"service_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.svc.Booking.CheckPromo(req.Code, req.ServiceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// POST /api/appointments
// body: {"doctor_id":"uuid","service_id":"uuid","starts_at":"2026-05-10T10:00:00Z","promo_code":""}
func (h *BookingHandler) Create(c *gin.Context) {
	patientID := c.GetString("user_id")
	var req struct {
		DoctorID  string `json:"doctor_id"  binding:"required"`
		ServiceID string `json:"service_id" binding:"required"`
		StartsAt  string `json:"starts_at"  binding:"required"`
		PromoCode string `json:"promo_code"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	startsAt, err := time.Parse(time.RFC3339, req.StartsAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid starts_at, use RFC3339"})
		return
	}
	a, err := h.svc.Booking.Book(service.BookRequest{
		PatientID: patientID,
		DoctorID:  req.DoctorID,
		ServiceID: req.ServiceID,
		PromoCode: req.PromoCode,
		StartsAt:  startsAt,
		CreatedBy: "patient",
	})
	if err != nil {
		if errors.Is(err, service.ErrSlotTaken) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, a)
}
