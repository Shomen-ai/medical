// Файл: internal/handler/doctor_portal.go
// Назначение: HTTP-обработчики рабочего места врача — приёмы на день, заполнение медкарты (черновик/финал) с переводом записи в completed, личное расписание на месяц и статистика приёмов.
package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"beautymed/internal/model"
	"beautymed/internal/service"

	"github.com/gin-gonic/gin"
)

// DoctorPortalHandler — обработчик запросов рабочего места врача.
type DoctorPortalHandler struct{ svc *service.Services }

// NewDoctorPortalHandler создаёт новый DoctorPortalHandler с подключённым сервисным слоем.
func NewDoctorPortalHandler(svc *service.Services) *DoctorPortalHandler {
	return &DoctorPortalHandler{svc: svc}
}

// TodayAppointments возвращает список приёмов врача на указанный день (по умолчанию — сегодня).
// GET /api/doctor/appointments?date=2026-05-10  (default: today)
func (h *DoctorPortalHandler) TodayAppointments(c *gin.Context) {
	doctorID := c.GetString("user_id")
	dateStr := c.Query("date")
	var date time.Time
	if dateStr == "" {
		date = time.Now().UTC().Truncate(24 * time.Hour)
	} else {
		var err error
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date, use YYYY-MM-DD"})
			return
		}
	}
	as, err := h.svc.Repos.Appointments.ListByDoctor(doctorID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, as)
}

// GetAppointment возвращает детали приёма врача вместе с уже сохранённой медкартой (если есть).
// GET /api/doctor/appointments/:id
func (h *DoctorPortalHandler) GetAppointment(c *gin.Context) {
	doctorID := c.GetString("user_id")
	a, err := h.svc.Repos.Appointments.FindByID(c.Param("id"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if a.DoctorID != doctorID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	record, err := h.svc.Repos.Appointments.GetRecord(a.ID)
	if errors.Is(err, sql.ErrNoRows) {
		record = nil
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"appointment": a, "record": record})
}

// SaveRecord сохраняет медкарту приёма (жалобы, диагноз, назначения); при is_draft=false переводит запись в статус completed.
// PATCH /api/doctor/appointments/:id/record
func (h *DoctorPortalHandler) SaveRecord(c *gin.Context) {
	doctorID := c.GetString("user_id")
	a, err := h.svc.Repos.Appointments.FindByID(c.Param("id"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if a.DoctorID != doctorID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	var req struct {
		Complaints      string  `json:"complaints"       binding:"required"`
		Diagnosis       string  `json:"diagnosis"        binding:"required"`
		Prescription    *string `json:"prescription"`
		Recommendations *string `json:"recommendations"`
		IsDraft         bool    `json:"is_draft"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rec := &model.AppointmentRecord{
		AppointmentID:   a.ID,
		Complaints:      req.Complaints,
		Diagnosis:       req.Diagnosis,
		Prescription:    req.Prescription,
		Recommendations: req.Recommendations,
		IsDraft:         req.IsDraft,
	}
	if err := h.svc.Repos.Appointments.UpsertRecord(rec); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !req.IsDraft {
		if err := h.svc.Repos.Appointments.UpdateStatus(a.ID, "completed"); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "saved"})
}

// MonthlySchedule возвращает расписание врача на указанный месяц.
// GET /api/doctor/schedule?year=2026&month=5
func (h *DoctorPortalHandler) MonthlySchedule(c *gin.Context) {
	doctorID := c.GetString("user_id")
	now := time.Now()
	year, month := now.Year(), int(now.Month())
	if y := c.Query("year"); y != "" {
		if v, err := strconv.Atoi(y); err == nil {
			year = v
		}
	}
	if m := c.Query("month"); m != "" {
		if v, err := strconv.Atoi(m); err == nil {
			month = v
		}
	}
	rows, err := h.svc.Repos.Appointments.ListScheduleByDoctor(doctorID, year, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rows)
}

// Stats возвращает агрегированную статистику приёмов текущего врача.
// GET /api/doctor/stats
func (h *DoctorPortalHandler) Stats(c *gin.Context) {
	stats, err := h.svc.Repos.Appointments.DoctorStats(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// Patients возвращает уникальных пациентов текущего врача с числом визитов (для отчёта в Excel).
// GET /api/doctor/patients
func (h *DoctorPortalHandler) Patients(c *gin.Context) {
	rows, err := h.svc.Repos.Appointments.DoctorPatients(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rows)
}
