// Файл: internal/service/slots.go
// Назначение: расчёт свободных временных слотов в рабочем окне врача с учётом длительности услуги и уже занятых интервалов.
package service

import (
	"time"

	"beautymed/internal/model"
	"beautymed/internal/repository"
)

// CalcSlots возвращает свободные слоты в рабочем окне с шагом, равным длительности услуги.
// CalcSlots returns available time slots given work window, duration and already booked ranges.
// A candidate slot is hidden if its [start, end) interval overlaps any taken range.
// All times are expected to be in UTC.
func CalcSlots(workStart, workEnd time.Time, durationMin int, taken []repository.TakenRange) []model.TimeSlot {
	duration := time.Duration(durationMin) * time.Minute
	var slots []model.TimeSlot
	cur := workStart.UTC()
	for !cur.Add(duration).After(workEnd.UTC()) {
		end := cur.Add(duration)
		if !overlapsAny(cur, end, taken) {
			slots = append(slots, model.TimeSlot{
				StartsAt: cur.Format("15:04"),
				EndsAt:   end.Format("15:04"),
			})
		}
		cur = cur.Add(duration)
	}
	return slots
}

// overlapsAny reports whether [start, end) intersects any of the given ranges.
// Intervals are treated as half-open: touching endpoints (end == other.start) do NOT overlap.
func overlapsAny(start, end time.Time, ranges []repository.TakenRange) bool {
	for _, r := range ranges {
		if start.Before(r.EndsAt) && r.StartsAt.Before(end) {
			return true
		}
	}
	return false
}
