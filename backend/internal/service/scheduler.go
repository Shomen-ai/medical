package service

import (
	"log"
	"time"

	"beautymed/internal/repository"

	"github.com/robfig/cron/v3"
)

type CronScheduler struct {
	c     *cron.Cron
	admin *repository.AdminRepo
	sms   *SMSService
}

func NewCronScheduler(admin *repository.AdminRepo, sms *SMSService) *CronScheduler {
	return &CronScheduler{
		c:     cron.New(cron.WithLocation(time.UTC)),
		admin: admin,
		sms:   sms,
	}
}

func (s *CronScheduler) Start() {
	// Send 24h appointment reminders daily at 10:00 UTC.
	s.c.AddFunc("0 10 * * *", s.sendReminders)

	// Log reminder to generate next month's schedule on the 25th at 08:00 UTC.
	s.c.AddFunc("0 8 25 * *", s.generateNextMonthSchedule)

	s.c.Start()
}

func (s *CronScheduler) Stop() {
	s.c.Stop()
}

func (s *CronScheduler) sendReminders() {
	reminders, err := s.admin.ListTomorrowAppointments()
	if err != nil {
		log.Printf("scheduler: list tomorrow appointments: %v", err)
		return
	}
	for _, r := range reminders {
		if err := s.sms.SendReminder(r.PatientPhone, r.DoctorName, r.StartsAt); err != nil {
			log.Printf("scheduler: send reminder to %s: %v", r.PatientPhone, err)
		}
	}
	log.Printf("scheduler: sent %d reminders", len(reminders))
}

// generateNextMonthSchedule logs a reminder — actual generation is done via admin panel.
func (s *CronScheduler) generateNextMonthSchedule() {
	now := time.Now().UTC()
	nextMonth := now.AddDate(0, 1, 0)
	log.Printf("scheduler: reminder — generate schedule for %s %d via admin panel",
		nextMonth.Month().String(), nextMonth.Year())
}
