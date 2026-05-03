package service

import (
	"time"

	"beautymed/internal/model"
)

// CalcSlots returns available time slots given work window, duration and already taken start times.
// All times are expected to be in UTC.
func CalcSlots(workStart, workEnd time.Time, durationMin int, taken []time.Time) []model.TimeSlot {
	takenSet := make(map[string]bool, len(taken))
	for _, t := range taken {
		takenSet[t.UTC().Format("15:04")] = true
	}

	duration := time.Duration(durationMin) * time.Minute
	var slots []model.TimeSlot
	cur := workStart.UTC()
	for !cur.Add(duration).After(workEnd.UTC()) {
		key := cur.Format("15:04")
		if !takenSet[key] {
			end := cur.Add(duration)
			slots = append(slots, model.TimeSlot{
				StartsAt: cur.Format("15:04"),
				EndsAt:   end.Format("15:04"),
			})
		}
		cur = cur.Add(duration)
	}
	return slots
}
