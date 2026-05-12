// Файл: internal/service/services.go
// Назначение: фабрика-сборщик всех сервисов и репозиториев BeautyMed в одну структуру для DI.
package service

import (
	"beautymed/internal/config"
	"beautymed/internal/repository"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

// Services — контейнер со всеми инициализированными сервисами и репозиториями приложения.
type Services struct {
	Auth      *AuthService
	OTP       *OTPService
	JWT       *JWTService
	Booking   *BookingService
	Schedule  *ScheduleService
	SMS       *SMSService
	Admin     *repository.AdminRepo
	Scheduler *CronScheduler
	Repos     *Repos
}

// Repos группирует базовые data-репозитории для удобного проброса в сервисы.
type Repos struct {
	Users        *repository.UserRepo
	Doctors      *repository.DoctorRepo
	Specialties  *repository.SpecialtyRepo
	Services     *repository.ServiceRepo
	Appointments *repository.AppointmentRepo
}

// New собирает все сервисы и репозитории BeautyMed на основе подключений к Postgres, Redis и конфига.
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
		Admin:     adminRepo,
		Scheduler: NewCronScheduler(adminRepo, repos.Specialties, repos.Doctors, schedule, sms),
		Repos:     repos,
	}
}
