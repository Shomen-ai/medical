package service

import (
	"beautymed/internal/config"
	"beautymed/internal/model"
	"beautymed/internal/repository"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Services struct {
	Auth      *AuthService
	OTP       *OTPService
	JWT       *JWTService
	Booking   *BookingService
	Schedule  *ScheduleService
	Admin     *AdminService
	Repos     *Repos
}

type Repos struct {
	Users        *repository.UserRepo
	Doctors      *repository.DoctorRepo
	Specialties  *repository.SpecialtyRepo
	Services     *repository.ServiceRepo
	Appointments *repository.AppointmentRepo
}

// AdminService wraps AdminRepo — replaced by direct *repository.AdminRepo in Task 9.
type AdminService struct{ *repository.AdminRepo }

// ScheduleService stub — replaced by real implementation in Task 5.
type ScheduleService struct{}

// SpecialtyGroup stub — moved to schedule.go in Task 5.
type SpecialtyGroup struct {
	SpecialtyID string
	DoctorIDs   []string
	StartTime   string
	EndTime     string
}

func (s *ScheduleService) Generate(year, month int, groups []SpecialtyGroup) []model.ScheduleRow {
	return nil
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
		OTP:      otp,
		JWT:      j,
		Auth:     NewAuthService(repos, otp, j),
		Booking:  NewBookingService(repos),
		Schedule: &ScheduleService{},
		Admin:    &AdminService{repository.NewAdminRepo(db)},
		Repos:    repos,
	}
}
