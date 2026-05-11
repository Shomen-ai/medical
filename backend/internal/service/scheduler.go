package service

import (
	"log"
	"time"

	"beautymed/internal/repository"

	"github.com/robfig/cron/v3"
)

type CronScheduler struct {
	c           *cron.Cron
	admin       *repository.AdminRepo
	specialties *repository.SpecialtyRepo
	doctors     *repository.DoctorRepo
	schedule    *ScheduleService
	sms         *SMSService
}

func NewCronScheduler(
	admin *repository.AdminRepo,
	specialties *repository.SpecialtyRepo,
	doctors *repository.DoctorRepo,
	schedule *ScheduleService,
	sms *SMSService,
) *CronScheduler {
	return &CronScheduler{
		c:           cron.New(cron.WithLocation(time.UTC)),
		admin:       admin,
		specialties: specialties,
		doctors:     doctors,
		schedule:    schedule,
		sms:         sms,
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

// generateNextMonthSchedule generates the 3/3 rotation schedule for the next calendar month
// across all specialties. Already-existing schedule rows for the same (doctor, date) are
// preserved with the new start/end/is_day_off via the ON CONFLICT upsert, but appointments
// are NOT touched — booked slots remain valid since they reference the appointments table.
func (s *CronScheduler) generateNextMonthSchedule() {
	now := time.Now().UTC()
	nextMonth := now.AddDate(0, 1, 0)
	year, month := nextMonth.Year(), int(nextMonth.Month())

	specialties, err := s.specialties.List()
	if err != nil {
		log.Printf("scheduler: list specialties: %v", err)
		return
	}

	var groups []SpecialtyGroup
	for _, sp := range specialties {
		doctors, err := s.doctors.ListBySpecialty(sp.ID)
		if err != nil {
			log.Printf("scheduler: list doctors for specialty %s: %v", sp.ID, err)
			continue
		}
		if len(doctors) == 0 {
			continue
		}
		ids := make([]string, len(doctors))
		for i, d := range doctors {
			ids[i] = d.ID
		}
		groups = append(groups, SpecialtyGroup{
			SpecialtyID: sp.ID,
			DoctorIDs:   ids,
			StartTime:   "09:00",
			EndTime:     "18:00",
		})
	}

	if len(groups) == 0 {
		log.Printf("scheduler: no specialties with doctors, nothing to generate for %s %d",
			nextMonth.Month(), year)
		return
	}

	rows := s.schedule.Generate(year, month, groups)
	if err := s.admin.BulkUpsertSchedule(rows); err != nil {
		log.Printf("scheduler: upsert schedule for %s %d: %v", nextMonth.Month(), year, err)
		return
	}
	log.Printf("scheduler: generated schedule for %s %d (%d rows, %d specialties)",
		nextMonth.Month(), year, len(rows), len(groups))
}
