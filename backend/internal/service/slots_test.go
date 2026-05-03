package service_test

import (
	"testing"
	"time"

	"beautymed/internal/service"

	"github.com/stretchr/testify/assert"
)

func TestSlots_NoTaken(t *testing.T) {
	start := time.Date(2026, 5, 10, 9, 0, 0, 0, time.UTC)
	end := time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC)
	taken := []time.Time{}

	slots := service.CalcSlots(start, end, 60, taken)
	assert.Len(t, slots, 3)
	assert.Equal(t, "09:00", slots[0].StartsAt)
	assert.Equal(t, "10:00", slots[1].StartsAt)
	assert.Equal(t, "11:00", slots[2].StartsAt)
}

func TestSlots_WithTaken(t *testing.T) {
	start := time.Date(2026, 5, 10, 9, 0, 0, 0, time.UTC)
	end := time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC)
	taken := []time.Time{
		time.Date(2026, 5, 10, 10, 0, 0, 0, time.UTC),
	}

	slots := service.CalcSlots(start, end, 60, taken)
	assert.Len(t, slots, 2)
	assert.Equal(t, "09:00", slots[0].StartsAt)
	assert.Equal(t, "11:00", slots[1].StartsAt)
}
