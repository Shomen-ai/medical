package service_test

import (
	"testing"
	"time"

	"beautymed/internal/service"

	"github.com/stretchr/testify/assert"
)

func TestGenerate3x3_TwoDoctors(t *testing.T) {
	// June 2026: 30 days, Sundays on 7,14,21,28 = 4 Sundays → 26 working days
	groups := []service.SpecialtyGroup{
		{
			SpecialtyID: "spec-1",
			DoctorIDs:   []string{"doc-A", "doc-B"},
			StartTime:   "09:00",
			EndTime:     "18:00",
		},
	}
	rows := service.Generate3x3(2026, 6, groups, nil)

	// 30 days × 2 doctors = 60 rows total
	assert.Len(t, rows, 60)

	// June 1 (Mon) — working day 0 (block 0) → doc-A on, doc-B off
	june1 := rowsFor(rows, "2026-06-01")
	assert.Len(t, june1, 2)
	assert.False(t, findDoctor(june1, "doc-A").IsDayOff, "doc-A should work on day 1")
	assert.True(t, findDoctor(june1, "doc-B").IsDayOff, "doc-B should be off on day 1")

	// June 3 (Wed) — working day 2 (block 0) → doc-A on
	june3 := rowsFor(rows, "2026-06-03")
	assert.False(t, findDoctor(june3, "doc-A").IsDayOff, "doc-A works day 3")

	// June 4 (Thu) — working day 3 (block 1) → doc-B on
	june4 := rowsFor(rows, "2026-06-04")
	assert.True(t, findDoctor(june4, "doc-A").IsDayOff, "doc-A off day 4")
	assert.False(t, findDoctor(june4, "doc-B").IsDayOff, "doc-B works day 4")

	// June 7 (Sun) — all is_day_off = true
	june7 := rowsFor(rows, "2026-06-07")
	assert.Len(t, june7, 2)
	for _, r := range june7 {
		assert.True(t, r.IsDayOff, "Sunday must be day off")
	}
}

func TestGenerate3x3_Holiday(t *testing.T) {
	groups := []service.SpecialtyGroup{
		{SpecialtyID: "s", DoctorIDs: []string{"doc-A", "doc-B"}, StartTime: "09:00", EndTime: "18:00"},
	}
	holiday := time.Date(2026, 6, 12, 0, 0, 0, 0, time.UTC) // arbitrary mid-month holiday
	rows := service.Generate3x3(2026, 6, groups, []time.Time{holiday})

	june12 := rowsFor(rows, "2026-06-12")
	for _, r := range june12 {
		assert.True(t, r.IsDayOff, "holiday must be day off")
	}
}

func TestGenerate3x3_SingleDoctor(t *testing.T) {
	groups := []service.SpecialtyGroup{
		{SpecialtyID: "s", DoctorIDs: []string{"doc-A"}, StartTime: "09:00", EndTime: "18:00"},
	}
	rows := service.Generate3x3(2026, 6, groups, nil)

	// Single doctor works all non-Sunday days
	june1 := rowsFor(rows, "2026-06-01")
	assert.Len(t, june1, 1)
	assert.False(t, june1[0].IsDayOff)

	june7 := rowsFor(rows, "2026-06-07")
	assert.Len(t, june7, 1)
	assert.True(t, june7[0].IsDayOff)
}

// helpers

func rowsFor(rows []service.ScheduleGenRow, date string) []service.ScheduleGenRow {
	var result []service.ScheduleGenRow
	for _, r := range rows {
		if r.WorkDate.Format("2006-01-02") == date {
			result = append(result, r)
		}
	}
	return result
}

func findDoctor(rows []service.ScheduleGenRow, id string) service.ScheduleGenRow {
	for _, r := range rows {
		if r.DoctorID == id {
			return r
		}
	}
	return service.ScheduleGenRow{}
}
