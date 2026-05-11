package service_test

import (
	"testing"
	"time"

	"beautymed/internal/repository"
	"beautymed/internal/service"

	"github.com/stretchr/testify/assert"
)

func TestSlots_NoTaken(t *testing.T) {
	start := time.Date(2026, 5, 10, 9, 0, 0, 0, time.UTC)
	end := time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC)

	slots := service.CalcSlots(start, end, 60, nil)
	assert.Len(t, slots, 3)
	assert.Equal(t, "09:00", slots[0].StartsAt)
	assert.Equal(t, "10:00", slots[1].StartsAt)
	assert.Equal(t, "11:00", slots[2].StartsAt)
}

func TestSlots_WithTaken(t *testing.T) {
	start := time.Date(2026, 5, 10, 9, 0, 0, 0, time.UTC)
	end := time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC)
	taken := []repository.TakenRange{
		{
			StartsAt: time.Date(2026, 5, 10, 10, 0, 0, 0, time.UTC),
			EndsAt:   time.Date(2026, 5, 10, 11, 0, 0, 0, time.UTC),
		},
	}

	slots := service.CalcSlots(start, end, 60, taken)
	assert.Len(t, slots, 2)
	assert.Equal(t, "09:00", slots[0].StartsAt)
	assert.Equal(t, "11:00", slots[1].StartsAt)
}

// Regression: a 60-min appointment at 10:00 should also block a 30-min booking at 10:30.
func TestSlots_OverlapBlocksShorterSlot(t *testing.T) {
	start := time.Date(2026, 5, 10, 9, 0, 0, 0, time.UTC)
	end := time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC)
	taken := []repository.TakenRange{
		{
			StartsAt: time.Date(2026, 5, 10, 10, 0, 0, 0, time.UTC),
			EndsAt:   time.Date(2026, 5, 10, 11, 0, 0, 0, time.UTC),
		},
	}

	slots := service.CalcSlots(start, end, 30, taken)
	starts := make([]string, len(slots))
	for i, s := range slots {
		starts[i] = s.StartsAt
	}
	assert.NotContains(t, starts, "10:00")
	assert.NotContains(t, starts, "10:30")
	assert.Contains(t, starts, "09:00")
	assert.Contains(t, starts, "09:30")
	assert.Contains(t, starts, "11:00")
	assert.Contains(t, starts, "11:30")
}

// Adjacent appointments (10:00–10:30 and 10:30–11:00) should both be allowed.
func TestSlots_AdjacentNotOverlap(t *testing.T) {
	start := time.Date(2026, 5, 10, 10, 0, 0, 0, time.UTC)
	end := time.Date(2026, 5, 10, 11, 0, 0, 0, time.UTC)
	taken := []repository.TakenRange{
		{
			StartsAt: time.Date(2026, 5, 10, 10, 0, 0, 0, time.UTC),
			EndsAt:   time.Date(2026, 5, 10, 10, 30, 0, 0, time.UTC),
		},
	}

	slots := service.CalcSlots(start, end, 30, taken)
	assert.Len(t, slots, 1)
	assert.Equal(t, "10:30", slots[0].StartsAt)
}
