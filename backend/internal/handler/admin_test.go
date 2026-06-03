package handler

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func d(y int, m time.Month, day int) time.Time {
	return time.Date(y, m, day, 0, 0, 0, 0, time.UTC)
}

func TestResolveReportRange_Default(t *testing.T) {
	now := d(2026, 6, 15)
	from, toExcl := resolveReportRange("", "", now)
	assert.Equal(t, d(2026, 6, 1), from)   // 1-е число текущего месяца
	assert.Equal(t, d(2026, 6, 16), toExcl) // сегодня включительно (today+1)
}

func TestResolveReportRange_Valid(t *testing.T) {
	now := d(2026, 6, 15)
	from, toExcl := resolveReportRange("2026-03-10", "2026-05-20", now)
	assert.Equal(t, d(2026, 3, 10), from)
	assert.Equal(t, d(2026, 5, 21), toExcl)
}

func TestResolveReportRange_FutureToClamped(t *testing.T) {
	now := d(2026, 6, 15)
	_, toExcl := resolveReportRange("2026-01-01", "2026-12-31", now)
	assert.Equal(t, d(2026, 6, 16), toExcl) // «по» в будущем → сегодня
}

func TestResolveReportRange_FromAfterTo(t *testing.T) {
	now := d(2026, 6, 15)
	from, toExcl := resolveReportRange("2026-05-20", "2026-03-10", now)
	assert.Equal(t, d(2026, 3, 10), from) // from подтянут к «по»
	assert.Equal(t, d(2026, 3, 11), toExcl)
}
