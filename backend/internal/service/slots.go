package service

import (
	"fmt"
	"time"

	"beautymed/internal/model"
)

// CalcSlots returns available time slots given work window, duration and already taken start times.
func CalcSlots(workStart, workEnd time.Time, durationMin int, taken []time.Time) []model.TimeSlot {
	takenSet := make(map[string]bool, len(taken))
	for _, t := range taken {
		takenSet[t.UTC().Format("15:04")] = true
	}

	duration := time.Duration(durationMin) * time.Minute
	var slots []model.TimeSlot
	cur := workStart
	for !cur.Add(duration).After(workEnd) {
		key := cur.UTC().Format("15:04")
		if !takenSet[key] {
			slots = append(slots, model.TimeSlot{
				StartsAt: fmt.Sprintf("%02d:%02d", cur.Hour(), cur.Minute()),
				EndsAt:   fmt.Sprintf("%02d:%02d", cur.Add(duration).Hour(), cur.Add(duration).Minute()),
			})
		}
		cur = cur.Add(duration)
	}
	return slots
}
