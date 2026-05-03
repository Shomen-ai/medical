package service

import (
	"fmt"
	"time"

	"beautymed/internal/model"
)

// SpecialtyGroup describes one specialty's doctors for schedule generation.
type SpecialtyGroup struct {
	SpecialtyID string
	DoctorIDs   []string
	StartTime   string // "09:00"
	EndTime     string // "18:00"
}

// ScheduleService wraps the Generate3x3 function for use in the Services struct.
type ScheduleService struct{}

// ScheduleGenRow is the output of Generate3x3 — one row per (doctor, date).
type ScheduleGenRow struct {
	DoctorID  string
	WorkDate  time.Time
	StartTime string
	EndTime   string
	IsDayOff  bool
}

// Generate3x3 generates a 3/3 rotation schedule for the given month.
// For specialties with 2+ doctors, doctors alternate in blocks of 3 working days.
// Sundays and dates in holidays are always day-off for every doctor.
func Generate3x3(year, month int, groups []SpecialtyGroup, holidays []time.Time) []ScheduleGenRow {
	holidaySet := make(map[string]bool, len(holidays))
	for _, h := range holidays {
		holidaySet[h.UTC().Format("2006-01-02")] = true
	}

	firstDay := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDay.AddDate(0, 1, -1)

	// Per-specialty counter of working days elapsed (for 3/3 block logic).
	workdayCount := make(map[string]int, len(groups))

	var rows []ScheduleGenRow
	for d := firstDay; !d.After(lastDay); d = d.AddDate(0, 0, 1) {
		isDayOff := d.Weekday() == time.Sunday || holidaySet[d.Format("2006-01-02")]

		for _, g := range groups {
			if isDayOff {
				for _, id := range g.DoctorIDs {
					rows = append(rows, ScheduleGenRow{
						DoctorID:  id,
						WorkDate:  d,
						StartTime: g.StartTime,
						EndTime:   g.EndTime,
						IsDayOff:  true,
					})
				}
				continue
			}

			if len(g.DoctorIDs) <= 1 {
				if len(g.DoctorIDs) == 1 {
					rows = append(rows, ScheduleGenRow{
						DoctorID:  g.DoctorIDs[0],
						WorkDate:  d,
						StartTime: g.StartTime,
						EndTime:   g.EndTime,
						IsDayOff:  false,
					})
				}
				continue
			}

			// 3/3 rotation: each block of 3 working days belongs to one doctor.
			n := workdayCount[g.SpecialtyID]
			onDuty := (n / 3) % len(g.DoctorIDs)
			workdayCount[g.SpecialtyID]++

			for i, id := range g.DoctorIDs {
				rows = append(rows, ScheduleGenRow{
					DoctorID:  id,
					WorkDate:  d,
					StartTime: g.StartTime,
					EndTime:   g.EndTime,
					IsDayOff:  i != onDuty,
				})
			}
		}
	}
	return rows
}

// Generate converts ScheduleGenRow results to model.ScheduleRow for DB insertion.
func (s *ScheduleService) Generate(year, month int, groups []SpecialtyGroup) []model.ScheduleRow {
	genRows := Generate3x3(year, month, groups, russianHolidays(year))
	result := make([]model.ScheduleRow, len(genRows))
	for i, r := range genRows {
		result[i] = model.ScheduleRow{
			DoctorID:  r.DoctorID,
			WorkDate:  r.WorkDate,
			StartTime: r.StartTime,
			EndTime:   r.EndTime,
			IsDayOff:  r.IsDayOff,
		}
	}
	return result
}

// russianHolidays returns the fixed Russian public holidays for the given year.
func russianHolidays(year int) []time.Time {
	dates := []string{
		// New Year block
		"%d-01-01", "%d-01-02", "%d-01-03", "%d-01-04",
		"%d-01-05", "%d-01-06", "%d-01-07", "%d-01-08",
		// Defender of the Fatherland Day
		"%d-02-23",
		// International Women's Day
		"%d-03-08",
		// Spring and Labor Day
		"%d-05-01",
		// Victory Day
		"%d-05-09",
		// Russia Day
		"%d-06-12",
		// National Unity Day
		"%d-11-04",
	}
	holidays := make([]time.Time, 0, len(dates))
	for _, f := range dates {
		t, err := time.Parse("2006-01-02", fmt.Sprintf(f, year))
		if err == nil {
			holidays = append(holidays, t)
		}
	}
	return holidays
}
