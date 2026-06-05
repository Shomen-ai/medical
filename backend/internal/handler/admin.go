// Файл: internal/handler/admin.go
// Назначение: HTTP-обработчики административной панели — KPI на главном экране, управление записями, врачами, расписанием, промокодами, выручкой и сводной статистикой.
package handler

import (
	"net/http"
	"strconv"
	"time"

	"beautymed/internal/model"
	"beautymed/internal/service"

	"github.com/gin-gonic/gin"
)

// AdminHandler — обработчик запросов административной части API.
type AdminHandler struct{ svc *service.Services }

// NewAdminHandler создаёт новый AdminHandler с подключённым сервисным слоем.
func NewAdminHandler(svc *service.Services) *AdminHandler { return &AdminHandler{svc: svc} }

// Dashboard возвращает KPI текущего дня и список сегодняшних записей для главного экрана админки.
// GET /api/admin/dashboard
func (h *AdminHandler) Dashboard(c *gin.Context) {
	count, revenue, free, err := h.svc.Admin.TodayKPI()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	today := time.Now().UTC().Truncate(24 * time.Hour)
	appts, err := h.svc.Admin.ListAllAppointments("", "", &today)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"kpi": gin.H{
			"appointments_today": count,
			"revenue_today":      revenue,
			"free_slots":         free,
		},
		"appointments": appts,
	})
}

// ListAppointments возвращает список записей с фильтрацией по статусу, врачу и дате.
// GET /api/admin/appointments?status=&doctor_id=&date=
func (h *AdminHandler) ListAppointments(c *gin.Context) {
	var date *time.Time
	if d := c.Query("date"); d != "" {
		t, err := time.Parse("2006-01-02", d)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date"})
			return
		}
		date = &t
	}
	as, err := h.svc.Admin.ListAllAppointments(c.Query("status"), c.Query("doctor_id"), date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, as)
}

// CreateAppointment создаёт запись пациента вручную через админ-панель (при необходимости создаёт пользователя по телефону).
// POST /api/admin/appointments — manual booking
// body: {"patient_phone":"+7...","doctor_id":"uuid","service_id":"uuid","starts_at":"RFC3339","promo_code":""}
func (h *AdminHandler) CreateAppointment(c *gin.Context) {
	var req struct {
		PatientPhone string `json:"patient_phone" binding:"required"`
		DoctorID     string `json:"doctor_id"     binding:"required"`
		ServiceID    string `json:"service_id"    binding:"required"`
		StartsAt     string `json:"starts_at"     binding:"required"`
		PromoCode    string `json:"promo_code"`
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
	patient, err := h.svc.Repos.Users.Create(req.PatientPhone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	a, err := h.svc.Booking.Book(service.BookRequest{
		PatientID: patient.ID,
		DoctorID:  req.DoctorID,
		ServiceID: req.ServiceID,
		PromoCode: req.PromoCode,
		StartsAt:  startsAt,
		CreatedBy: "admin",
	})
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, a)
}

// CreateDoctor добавляет нового врача и создаёт для него сотрудника с ролью doctor.
// POST /api/admin/doctors
// body: {"full_name":"...","specialty_id":"uuid","phone":"+7...","bio":"...","photo_url":"...","experience_years":5}
func (h *AdminHandler) CreateDoctor(c *gin.Context) {
	var req struct {
		FullName        string `json:"full_name"        binding:"required"`
		SpecialtyID     string `json:"specialty_id"     binding:"required"`
		Phone           string `json:"phone"            binding:"required"`
		Bio             string `json:"bio"`
		PhotoURL        string `json:"photo_url"`
		ExperienceYears int    `json:"experience_years"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	d := &model.Doctor{
		FullName:        req.FullName,
		SpecialtyID:     req.SpecialtyID,
		Phone:           req.Phone,
		Bio:             req.Bio,
		PhotoURL:        req.PhotoURL,
		ExperienceYears: req.ExperienceYears,
		IsActive:        true,
	}
	if err := h.svc.Repos.Doctors.CreateDoctor(d); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.Repos.Doctors.CreateStaff(d.ID, d.Phone, "doctor"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, d)
}

// UpdateDoctor обновляет данные карточки врача по идентификатору.
// PATCH /api/admin/doctors/:id
func (h *AdminHandler) UpdateDoctor(c *gin.Context) {
	d, err := h.svc.Repos.Doctors.FindByID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "doctor not found"})
		return
	}
	var req struct {
		FullName        string `json:"full_name"`
		Bio             string `json:"bio"`
		PhotoURL        string `json:"photo_url"`
		ExperienceYears int    `json:"experience_years"`
		IsActive        *bool  `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.FullName != "" {
		d.FullName = req.FullName
	}
	if req.Bio != "" {
		d.Bio = req.Bio
	}
	if req.PhotoURL != "" {
		d.PhotoURL = req.PhotoURL
	}
	if req.ExperienceYears != 0 {
		d.ExperienceYears = req.ExperienceYears
	}
	if req.IsActive != nil {
		d.IsActive = *req.IsActive
	}
	if err := h.svc.Repos.Doctors.UpdateDoctor(d); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, d)
}

// GetSchedule возвращает сетку расписания врачей на указанный месяц и (опционально) специальность.
// GET /api/admin/schedule?year=2026&month=6&specialty_id=uuid
func (h *AdminHandler) GetSchedule(c *gin.Context) {
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
	cells, err := h.svc.Admin.ListScheduleForMonth(year, month, c.Query("specialty_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"year": year, "month": month, "cells": cells})
}

// UpdateScheduleCell обновляет одну ячейку расписания (рабочие часы или выходной) для пары врач+дата.
// PATCH /api/admin/schedule/:doctor_id/:date
// body: {"start_time":"09:00","end_time":"18:00","is_day_off":false}
func (h *AdminHandler) UpdateScheduleCell(c *gin.Context) {
	var req struct {
		StartTime string `json:"start_time" binding:"required"`
		EndTime   string `json:"end_time"   binding:"required"`
		IsDayOff  bool   `json:"is_day_off"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.Admin.UpsertScheduleCell(
		c.Param("doctor_id"), c.Param("date"),
		req.StartTime, req.EndTime, req.IsDayOff,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

// GenerateSchedule автоматически формирует расписание на месяц для всех активных врачей указанной специальности.
// POST /api/admin/schedule/generate
// body: {"year":2026,"month":6,"specialty_id":"uuid"}
func (h *AdminHandler) GenerateSchedule(c *gin.Context) {
	var req struct {
		Year        int    `json:"year"         binding:"required"`
		Month       int    `json:"month"        binding:"required"`
		SpecialtyID string `json:"specialty_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	doctors, err := h.svc.Repos.Doctors.ListBySpecialty(req.SpecialtyID)
	if err != nil || len(doctors) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no active doctors in specialty"})
		return
	}
	ids := make([]string, len(doctors))
	for i, d := range doctors {
		ids[i] = d.ID
	}
	rows := h.svc.Schedule.Generate(req.Year, req.Month, []service.SpecialtyGroup{
		{SpecialtyID: req.SpecialtyID, DoctorIDs: ids, StartTime: "09:00", EndTime: "18:00"},
	})
	if err := h.svc.Admin.BulkUpsertSchedule(rows); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"generated": len(rows)})
}

// ListPromos возвращает все промокоды клиники.
// GET /api/admin/promo
func (h *AdminHandler) ListPromos(c *gin.Context) {
	pcs, err := h.svc.Admin.ListPromos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pcs)
}

// CreatePromo создаёт новый промокод с заданной скидкой и периодом действия.
// POST /api/admin/promo
// body: {"code":"BEAUTY20","discount_pct":20,"max_uses":100,"valid_from":"2026-06-01","valid_until":"2026-08-31"}
func (h *AdminHandler) CreatePromo(c *gin.Context) {
	var req struct {
		Code        string  `json:"code"         binding:"required"`
		DiscountPct int     `json:"discount_pct" binding:"required,min=1,max=100"`
		MaxUses     *int    `json:"max_uses"`
		ValidFrom   string  `json:"valid_from"   binding:"required"`
		ValidUntil  *string `json:"valid_until"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pc := &model.PromoCode{
		Code:        req.Code,
		DiscountPct: req.DiscountPct,
		MaxUses:     req.MaxUses,
		ValidFrom:   req.ValidFrom,
		ValidUntil:  req.ValidUntil,
		IsActive:    true,
	}
	if err := h.svc.Admin.CreatePromoCode(pc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, pc)
}

// periodRange returns [from, to) for the named period. Returns false if period is unknown.
func periodRange(period string) (from, to time.Time, ok bool) {
	now := time.Now().UTC()
	ok = true
	switch period {
	case "day":
		from = now.Truncate(24 * time.Hour)
		to = from.Add(24 * time.Hour)
	case "week":
		offset := int(now.Weekday()) - 1
		if offset < 0 {
			offset = 6
		}
		from = now.Truncate(24 * time.Hour).AddDate(0, 0, -offset)
		to = from.AddDate(0, 0, 7)
	case "month":
		from = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		to = from.AddDate(0, 1, 0)
	case "quarter":
		q := (int(now.Month()) - 1) / 3
		from = time.Date(now.Year(), time.Month(q*3+1), 1, 0, 0, 0, 0, time.UTC)
		to = from.AddDate(0, 3, 0)
	case "year":
		from = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
		to = from.AddDate(1, 0, 0)
	default:
		ok = false
	}
	return
}

// resolveReportRange парсит даты «с»/«по» (YYYY-MM-DD) в полузакрытый интервал [from, toExcl).
// Дефолт — текущий месяц. «По» клампится к сегодняшнему дню (запрет будущих дат). Гарантирует from<=to.
func resolveReportRange(fromStr, toStr string, now time.Time) (from, toExcl time.Time) {
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	// from: по умолчанию — 1-е число текущего месяца.
	from = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	if d, err := time.Parse("2006-01-02", fromStr); err == nil {
		from = time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
	}
	// to: по умолчанию — сегодня; будущее клампим к сегодня.
	toDay := today
	if d, err := time.Parse("2006-01-02", toStr); err == nil {
		toDay = time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
	}
	if toDay.After(today) {
		toDay = today
	}
	if from.After(toDay) {
		from = toDay
	}
	toExcl = toDay.AddDate(0, 0, 1) // «по» включительно
	return
}

// Revenue возвращает суммарную выручку клиники за выбранный период.
// GET /api/admin/revenue?period=week  (day|week|month|quarter|year)
func (h *AdminHandler) Revenue(c *gin.Context) {
	period := c.DefaultQuery("period", "day")
	from, to, ok := periodRange(period)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid period"})
		return
	}
	total, err := h.svc.Admin.Revenue(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"period": period, "from": from, "to": to, "total": total})
}

// PeriodStats возвращает статистику записей и выручки за интервал дат [from, to].
// GET /api/admin/stats/period?from=2026-01-01&to=2026-06-03  (даты YYYY-MM-DD; «to» не больше сегодня)
func (h *AdminHandler) PeriodStats(c *gin.Context) {
	from, to := resolveReportRange(c.Query("from"), c.Query("to"), time.Now().UTC())
	s, err := h.svc.Admin.PeriodStats(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, s)
}

// StatsByDoctor возвращает разбивку приёмов и уникальных пациентов по каждому врачу за интервал дат.
// GET /api/admin/stats/by-doctor?from=2026-01-01&to=2026-06-03
func (h *AdminHandler) StatsByDoctor(c *gin.Context) {
	from, to := resolveReportRange(c.Query("from"), c.Query("to"), time.Now().UTC())
	rows, err := h.svc.Admin.ByDoctorStats(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"from": from, "to": to, "doctors": rows})
}

// Stats возвращает общую сводную статистику по клинике за всё время.
// GET /api/admin/stats
func (h *AdminHandler) Stats(c *gin.Context) {
	s, err := h.svc.Admin.OverallStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, s)
}

// FullReport возвращает данные для расширенного Excel-отчёта (4 листа) за интервал дат:
// выручка по услугам, загрузка по времени, рейтинги, новые/вернувшиеся пациенты.
// GET /api/admin/report/full?from=2026-01-01&to=2026-06-03
func (h *AdminHandler) FullReport(c *gin.Context) {
	from, to := resolveReportRange(c.Query("from"), c.Query("to"), time.Now().UTC())
	a := h.svc.Admin
	summary, err := a.ReportSummary(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	bySpecialty, err := a.ReportBySpecialty(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	byDoctor, err := a.ReportByDoctorFull(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	daily, err := a.DailyStats(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	genderRows, err := a.DemographicsGender()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ageRows, err := a.DemographicsAge()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	byService, err := a.ReportByService(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	byWeekday, err := a.LoadByWeekday(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	byHour, err := a.LoadByHour(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ratingsDoctor, err := a.RatingsByDoctor(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ratingsService, err := a.RatingsByService(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	retention, err := a.Retention(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"from": from, "to": to,
		"summary":            summary,
		"by_service":         byService,
		"by_specialty":       bySpecialty,
		"by_doctor":          byDoctor,
		"daily":              daily,
		"by_weekday":         byWeekday,
		"by_hour":            byHour,
		"ratings_by_doctor":  ratingsDoctor,
		"ratings_by_service": ratingsService,
		"retention":          retention,
		"gender":             genderRows,
		"age":                ageRows,
	})
}

// MonthlyStats возвращает помесячную разбивку за интервал дат — каждый месяц периода
// (включая месяцы без записей) для построения графика.
// GET /api/admin/stats/monthly?from=2026-01-01&to=2026-06-03
func (h *AdminHandler) MonthlyStats(c *gin.Context) {
	from, to := resolveReportRange(c.Query("from"), c.Query("to"), time.Now().UTC())
	pts, err := h.svc.Admin.MonthlyStatsRange(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pts)
}
