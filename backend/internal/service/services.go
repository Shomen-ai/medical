package service

import (
	"beautymed/internal/config"
	"beautymed/internal/repository"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Services struct {
	Auth    *AuthService
	OTP     *OTPService
	JWT     *JWTService
	Booking *BookingService
	Repos   *Repos
}

type Repos struct {
	Users        *repository.UserRepo
	Doctors      *repository.DoctorRepo
	Specialties  *repository.SpecialtyRepo
	Services     *repository.ServiceRepo
	Appointments *repository.AppointmentRepo
}

func New(db *sqlx.DB, rdb *redis.Client, cfg *config.Config) *Services {
	repos := &Repos{
		Users:        repository.NewUserRepo(db),
		Doctors:      repository.NewDoctorRepo(db),
		Specialties:  repository.NewSpecialtyRepo(db),
		Services:     repository.NewServiceRepo(db),
		Appointments: repository.NewAppointmentRepo(db),
	}
	otp := NewOTPService(rdb)
	j := NewJWTService(cfg.JWTSecret)
	return &Services{
		OTP:     otp,
		JWT:     j,
		Auth:    NewAuthService(repos, otp, j),
		Booking: NewBookingService(repos),
		Repos:   repos,
	}
}

// BookingService is a stub — replaced in Task 10.
type BookingService struct{}

func NewBookingService(_ *Repos) *BookingService { return &BookingService{} }
