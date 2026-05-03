package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"beautymed/internal/service"

	"github.com/gin-gonic/gin"
)

type CabinetHandler struct{ svc *service.Services }

func NewCabinetHandler(svc *service.Services) *CabinetHandler { return &CabinetHandler{svc: svc} }

// GET /api/cabinet/appointments
func (h *CabinetHandler) ListAppointments(c *gin.Context) {
	patientID := c.GetString("user_id")
	as, err := h.svc.Repos.Appointments.ListByPatient(patientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, as)
}

// GET /api/cabinet/appointments/:id
func (h *CabinetHandler) GetAppointment(c *gin.Context) {
	a, err := h.svc.Repos.Appointments.FindByID(c.Param("id"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "appointment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if a.PatientID != c.GetString("user_id") {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	record, _ := h.svc.Repos.Appointments.GetRecord(a.ID)
	c.JSON(http.StatusOK, gin.H{"appointment": a, "record": record})
}

// PATCH /api/cabinet/appointments/:id/cancel
func (h *CabinetHandler) Cancel(c *gin.Context) {
	a, err := h.svc.Repos.Appointments.FindByID(c.Param("id"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if a.PatientID != c.GetString("user_id") {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	if a.Status != "scheduled" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only scheduled appointments can be cancelled"})
		return
	}
	if err := h.svc.Repos.Appointments.UpdateStatus(a.ID, "cancelled"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "cancelled"})
}

// PATCH /api/cabinet/appointments/:id/reschedule
// body: {"starts_at":"2026-05-15T11:00:00Z"}
func (h *CabinetHandler) Reschedule(c *gin.Context) {
	a, err := h.svc.Repos.Appointments.FindByID(c.Param("id"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if a.PatientID != c.GetString("user_id") {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	if time.Until(a.StartsAt) < 2*time.Hour {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot reschedule less than 2 hours before appointment"})
		return
	}
	var req struct {
		StartsAt string `json:"starts_at" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newStart, err := time.Parse(time.RFC3339, req.StartsAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid starts_at"})
		return
	}
	svc, err := h.svc.Repos.Services.FindByID(a.ServiceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := h.validateSlot(a.DoctorID, a.ServiceID, newStart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newEnd := newStart.Add(time.Duration(svc.DurationMin) * time.Minute)
	if err := h.svc.Repos.Appointments.Reschedule(a.ID, newStart, newEnd); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "rescheduled", "starts_at": newStart})
}

// GET /api/cabinet/profile
func (h *CabinetHandler) GetProfile(c *gin.Context) {
	u, err := h.svc.Repos.Users.FindByID(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, u)
}

// PATCH /api/cabinet/profile
func (h *CabinetHandler) UpdateProfile(c *gin.Context) {
	u, err := h.svc.Repos.Users.FindByID(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	var req struct {
		FullName string  `json:"full_name"`
		Email    *string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.FullName != "" {
		u.FullName = req.FullName
	}
	if req.Email != nil {
		u.Email = req.Email
	}
	if err := h.svc.Repos.Users.Update(u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, u)
}

func (h *CabinetHandler) validateSlot(doctorID, serviceID string, newStart time.Time) error {
	dateStr := newStart.UTC().Format("2006-01-02")
	slots, err := h.svc.Booking.GetSlots(doctorID, serviceID, dateStr)
	if err != nil {
		return err
	}
	targetTime := newStart.UTC().Format("15:04")
	for _, s := range slots {
		if s.StartsAt == targetTime {
			return nil
		}
	}
	return errors.New("selected slot is not available")
}
