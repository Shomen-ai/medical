package handler

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidBirthDate(t *testing.T) {
	now := time.Date(2026, 6, 5, 3, 0, 0, 0, time.UTC)

	// Совершеннолетний — ок.
	_, ok := validBirthDate("1990-01-01", now)
	assert.True(t, ok, "взрослый должен проходить")

	// Ровно 18 лет (день рождения раньше сегодняшней даты-18) — ок.
	_, ok = validBirthDate("2008-06-04", now)
	assert.True(t, ok, "ровно 18 лет должно проходить")

	// Младше 18 — нельзя.
	_, ok = validBirthDate("2010-01-01", now)
	assert.False(t, ok, "младше 18 нельзя")

	// Новорождённый (родился сегодня) — нельзя.
	_, ok = validBirthDate("2026-06-05", now)
	assert.False(t, ok, "новорождённый нельзя")

	// Будущая дата — нельзя.
	_, ok = validBirthDate("2030-01-01", now)
	assert.False(t, ok, "будущая дата нельзя")

	// Некорректный формат — нельзя.
	_, ok = validBirthDate("not-a-date", now)
	assert.False(t, ok, "мусор нельзя")
}
