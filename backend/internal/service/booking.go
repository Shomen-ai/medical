// Файл: internal/service/booking.go
// Назначение: бизнес-логика бронирования — расчёт свободных слотов, проверка промокода, атомарное создание записи.
package service

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"beautymed/internal/model"
	"beautymed/internal/repository"
)

// BookingService — сервис бронирования приёмов.
type BookingService struct{ repos *Repos }

// NewBookingService создаёт сервис бронирования поверх набора репозиториев.
func NewBookingService(repos *Repos) *BookingService { return &BookingService{repos} }

// ErrSlotTaken — реэкспорт ошибки занятого слота из repository для использования в хендлерах.
// ErrSlotTaken is re-exported from the repository so handlers can keep using service.ErrSlotTaken.
var ErrSlotTaken = repository.ErrSlotTaken

// ErrDayOff сигнализирует, что у врача в выбранный день нет рабочей смены.
var ErrDayOff = errors.New("doctor does not work on this day")

// GetSlots возвращает свободные временные слоты у врача на дату с учётом длительности услуги.
func (s *BookingService) GetSlots(doctorID, serviceID, dateStr string) ([]model.TimeSlot, error) {
	svc, err := s.repos.Services.FindByID(serviceID)
	if err != nil {
		return nil, fmt.Errorf("service not found: %w", err)
	}
	sched, err := s.repos.Appointments.GetSchedule(doctorID, dateStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrDayOff
		}
		return nil, fmt.Errorf("get schedule: %w", err)
	}
	if sched.IsDayOff {
		return nil, ErrDayOff
	}
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, errors.New("invalid date format, use YYYY-MM-DD")
	}
	parseTime := func(s string) time.Time {
		t, _ := time.Parse("15:04:05", s)
		return time.Date(date.Year(), date.Month(), date.Day(),
			t.Hour(), t.Minute(), 0, 0, time.UTC)
	}
	workStart := parseTime(sched.StartTime)
	workEnd := parseTime(sched.EndTime)

	taken, err := s.repos.Appointments.TakenSlots(doctorID, date)
	if err != nil {
		return nil, err
	}
	return CalcSlots(workStart, workEnd, svc.DurationMin, taken), nil
}

// GetAvailableDates возвращает рабочие даты врача в указанном месяце (формат YYYY-MM).
func (s *BookingService) GetAvailableDates(doctorID, monthStr string) ([]string, error) {
	t, err := time.Parse("2006-01", monthStr)
	if err != nil {
		return nil, errors.New("invalid month, use YYYY-MM")
	}
	return s.repos.Appointments.ListWorkDates(doctorID, t.Year(), int(t.Month()))
}

// PromoCheckResult — результат проверки промокода с пересчитанной ценой услуги.
type PromoCheckResult struct {
	Valid         bool    `json:"valid"`
	DiscountPct   int     `json:"discount_pct"`
	OriginalPrice float64 `json:"original_price"`
	FinalPrice    float64 `json:"final_price"`
}

// CheckPromo проверяет промокод применительно к услуге и возвращает итоговую цену со скидкой.
// CheckPromo validates a promo code against a service and returns the recalculated price.
// Returns a result with Valid=false if the code is missing, unknown, expired, or exhausted.
func (s *BookingService) CheckPromo(code, serviceID string) (*PromoCheckResult, error) {
	svc, err := s.repos.Services.FindByID(serviceID)
	if err != nil {
		return nil, fmt.Errorf("service not found: %w", err)
	}
	result := &PromoCheckResult{
		OriginalPrice: svc.Price,
		FinalPrice:    svc.Price,
	}
	if code == "" {
		return result, nil
	}
	pc, err := s.repos.Appointments.FindPromoCode(code)
	if err != nil {
		return result, nil
	}
	result.Valid = true
	result.DiscountPct = pc.DiscountPct
	result.FinalPrice = svc.Price * float64(100-pc.DiscountPct) / 100
	return result, nil
}

// BookRequest — входные параметры для создания записи на приём.
type BookRequest struct {
	PatientID string
	DoctorID  string
	ServiceID string
	PromoCode string
	StartsAt  time.Time
	CreatedBy string
}

// Book создаёт запись на приём, применяя скидку и инкрементируя счётчик промокода в одной транзакции.
func (s *BookingService) Book(req BookRequest) (*model.Appointment, error) {
	svc, err := s.repos.Services.FindByID(req.ServiceID)
	if err != nil {
		return nil, fmt.Errorf("service not found: %w", err)
	}
	price := svc.Price
	var promoID *string

	if req.PromoCode != "" {
		pc, err := s.repos.Appointments.FindPromoCode(req.PromoCode)
		if err == nil {
			price = price * float64(100-pc.DiscountPct) / 100
			promoID = &pc.ID
		}
	}

	endsAt := req.StartsAt.Add(time.Duration(svc.DurationMin) * time.Minute)

	a := &model.Appointment{
		PatientID:   req.PatientID,
		DoctorID:    req.DoctorID,
		ServiceID:   req.ServiceID,
		PromoCodeID: promoID,
		StartsAt:    req.StartsAt,
		EndsAt:      endsAt,
		FinalPrice:  price,
		CreatedBy:   req.CreatedBy,
	}
	if err := s.repos.Appointments.BookAtomic(a, promoID); err != nil {
		if errors.Is(err, repository.ErrSlotTaken) {
			return nil, ErrSlotTaken
		}
		return nil, fmt.Errorf("create appointment: %w", err)
	}
	return a, nil
}
