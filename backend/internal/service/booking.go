package service

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"beautymed/internal/model"
)

type BookingService struct{ repos *Repos }

func NewBookingService(repos *Repos) *BookingService { return &BookingService{repos} }

var ErrSlotTaken = errors.New("slot is already taken")
var ErrDayOff = errors.New("doctor does not work on this day")

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

func (s *BookingService) GetAvailableDates(doctorID, monthStr string) ([]string, error) {
	t, err := time.Parse("2006-01", monthStr)
	if err != nil {
		return nil, errors.New("invalid month, use YYYY-MM")
	}
	return s.repos.Appointments.ListWorkDates(doctorID, t.Year(), int(t.Month()))
}

type BookRequest struct {
	PatientID string
	DoctorID  string
	ServiceID string
	PromoCode string
	StartsAt  time.Time
	CreatedBy string
}

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
			_ = s.repos.Appointments.IncrementPromoUsage(pc.ID)
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
	if err := s.repos.Appointments.Create(a); err != nil {
		return nil, fmt.Errorf("create appointment: %w", err)
	}
	return a, nil
}
