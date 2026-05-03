package router

import (
	"beautymed/internal/config"
	"beautymed/internal/handler"
	"beautymed/internal/middleware"
	"beautymed/internal/service"

	"github.com/gin-gonic/gin"
)

func New(svc *service.Services, cfg *config.Config) *gin.Engine {
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	patientMw := middleware.Auth(svc.JWT, "patient")
	doctorMw := middleware.Auth(svc.JWT, "doctor")
	adminMw := middleware.Auth(svc.JWT, "admin")

	authH := handler.NewAuthHandler(svc, cfg)
	specH := handler.NewSpecialtyHandler(svc)
	docH := handler.NewDoctorHandler(svc)
	svcH := handler.NewServiceHandler(svc)
	bookH := handler.NewBookingHandler(svc)
	cabH := handler.NewCabinetHandler(svc)
	drH := handler.NewDoctorPortalHandler(svc)
	adminH := handler.NewAdminHandler(svc)

	api := r.Group("/api")
	{
		// Auth — patients
		api.POST("/auth/otp", authH.SendOTP)
		api.POST("/auth/verify", authH.VerifyOTP)
		api.POST("/auth/refresh", authH.Refresh)

		// Auth — staff
		api.POST("/staff/auth/otp", authH.SendStaffOTP)
		api.POST("/staff/auth/verify", authH.VerifyStaffOTP)

		// Public read
		api.GET("/specialties", specH.List)
		api.GET("/doctors", docH.List)
		api.GET("/doctors/:id", docH.Get)
		api.GET("/services", svcH.List)
		api.GET("/services/:id", svcH.Get)
		api.GET("/doctors/:id/slots", bookH.GetSlots)
		api.GET("/doctors/:id/available-dates", bookH.GetAvailableDates)

		// Booking (patient)
		api.POST("/appointments", patientMw, bookH.Create)

		// Patient cabinet
		cabinet := api.Group("/cabinet", patientMw)
		{
			cabinet.GET("/appointments", cabH.ListAppointments)
			cabinet.GET("/appointments/:id", cabH.GetAppointment)
			cabinet.PATCH("/appointments/:id/cancel", cabH.Cancel)
			cabinet.PATCH("/appointments/:id/reschedule", cabH.Reschedule)
			cabinet.GET("/profile", cabH.GetProfile)
			cabinet.PATCH("/profile", cabH.UpdateProfile)
			cabinet.GET("/receipts", cabH.Receipts)
		}

		// Doctor portal
		doctor := api.Group("/doctor", doctorMw)
		{
			doctor.GET("/appointments", drH.TodayAppointments)
			doctor.GET("/appointments/:id", drH.GetAppointment)
			doctor.PATCH("/appointments/:id/record", drH.SaveRecord)
			doctor.GET("/schedule", drH.MonthlySchedule)
			doctor.GET("/stats", drH.Stats)
		}

		// Admin panel
		admin := api.Group("/admin", adminMw)
		{
			admin.GET("/dashboard", adminH.Dashboard)
			admin.GET("/appointments", adminH.ListAppointments)
			admin.POST("/appointments", adminH.CreateAppointment)
			admin.GET("/doctors", docH.List)
			admin.POST("/doctors", adminH.CreateDoctor)
			admin.PATCH("/doctors/:id", adminH.UpdateDoctor)
			admin.GET("/schedule", adminH.GetSchedule)
			admin.POST("/schedule/generate", adminH.GenerateSchedule)
			admin.PATCH("/schedule/:doctor_id/:date", adminH.UpdateScheduleCell)
			admin.GET("/promo", adminH.ListPromos)
			admin.POST("/promo", adminH.CreatePromo)
			admin.GET("/revenue", adminH.Revenue)
			admin.GET("/stats", adminH.Stats)
		}
	}

	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	return r
}
