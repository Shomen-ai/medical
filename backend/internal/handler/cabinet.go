// Файл: internal/handler/cabinet.go
// Назначение: HTTP-обработчики личного кабинета пациента — просмотр и управление своими записями (отмена, перенос) и редактирование профиля.
package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"beautymed/internal/service"

	"github.com/gin-gonic/gin"
)

// CabinetHandler — обработчик запросов личного кабинета пациента.
type CabinetHandler struct{ svc *service.Services }

// NewCabinetHandler создаёт новый CabinetHandler с подключённым сервисным слоем.
func NewCabinetHandler(svc *service.Services) *CabinetHandler { return &CabinetHandler{svc: svc} }

// ListAppointments возвращает все записи текущего пациента.
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

// GetAppointment возвращает детали записи пациента вместе с медкартой (если она уже заполнена).
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
	record, err := h.svc.Repos.Appointments.GetRecord(a.ID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		record = nil
	}
	c.JSON(http.StatusOK, gin.H{"appointment": a, "record": record})
}

// Cancel отменяет запись пациента (только если она ещё в статусе scheduled).
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

// Reschedule переносит запись на другое время с проверкой свободного слота и правила «не позже чем за 2 часа».
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
	if a.Status != "scheduled" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only scheduled appointments can be rescheduled"})
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

// GetProfile возвращает анкету текущего пациента.
// GET /api/cabinet/profile
func (h *CabinetHandler) GetProfile(c *gin.Context) {
	u, err := h.svc.Repos.Users.FindByID(c.GetString("user_id"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, u)
}

// UpdateProfile обновляет данные медкарты пациента: ФИО, дату рождения, пол, адрес, контакты, номер удостоверения личности.
// PATCH /api/cabinet/profile
func (h *CabinetHandler) UpdateProfile(c *gin.Context) {
	u, err := h.svc.Repos.Users.FindByID(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	var req struct {
		FullName      string  `json:"full_name"`
		Email         *string `json:"email"`
		BirthDate     *string `json:"birth_date"`
		Gender        *string `json:"gender"`
		Address       *string `json:"address"`
		IDDocNumber   *string `json:"id_doc_number"`
		IDDocIssuedBy *string `json:"id_doc_issued_by"`
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
	if req.BirthDate != nil && *req.BirthDate != "" {
		d, ok := validBirthDate(*req.BirthDate, time.Now().UTC())
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "birth_date must be in the past and the patient must be at least 18 years old"})
			return
		}
		u.BirthDate = &d
	} else if req.BirthDate != nil {
		// Пустая строка — сброс даты рождения.
		u.BirthDate = nil
	}
	if req.Gender != nil {
		// Принимаем только 'm', 'f' либо пустую строку (=сброс).
		g := *req.Gender
		if g == "m" || g == "f" {
			u.Gender = &g
		} else if g == "" {
			u.Gender = nil
		}
	}
	if req.Address != nil {
		u.Address = req.Address
	}
	if req.IDDocNumber != nil {
		u.IDDocNumber = req.IDDocNumber
	}
	if req.IDDocIssuedBy != nil {
		u.IDDocIssuedBy = req.IDDocIssuedBy
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

// validBirthDate парсит дату рождения (YYYY-MM-DD) и проверяет, что она не в будущем и
// что на момент now пациенту исполнилось не меньше 18 лет. Возвращает (дата, ok).
// Чистая функция — покрыта unit-тестами.
func validBirthDate(s string, now time.Time) (time.Time, bool) {
	d, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, false
	}
	// Граница: самая поздняя допустимая дата рождения — ровно «now минус 18 лет».
	// Будущие даты и возраст < 18 (включая новорождённых) оказываются позже границы.
	eighteenAgo := now.AddDate(-18, 0, 0)
	if d.After(eighteenAgo) {
		return time.Time{}, false
	}
	return d, true
}
