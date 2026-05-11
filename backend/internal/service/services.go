package service

import (
	"beautymed/internal/config"
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
	SMS       *SMSService
	PDF       *PDFService
	Admin     *repository.AdminRepo
	Scheduler *CronScheduler
	Repos     *Repos
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
	adminRepo := repository.NewAdminRepo(db)
	otp := NewOTPService(rdb)
	j := NewJWTService(cfg.JWTSecret)
	sms := NewSMSService(cfg.SMSCLogin, cfg.SMSCPassword)
	schedule := &ScheduleService{}
	return &Services{
		OTP:       otp,
		JWT:       j,
		Auth:      NewAuthService(repos, otp, j),
		Booking:   NewBookingService(repos),
		Schedule:  schedule,
		SMS:       sms,
		PDF:       NewPDFService(cfg.Clinic),
		Admin:     adminRepo,
		Scheduler: NewCronScheduler(adminRepo, repos.Specialties, repos.Doctors, schedule, sms),
		Repos:     repos,
	}
}
