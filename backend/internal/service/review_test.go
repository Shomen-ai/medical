package service_test

import (
	"strings"
	"testing"

	"beautymed/internal/service"

	"github.com/stretchr/testify/assert"
)

func TestValidateReviewInput_OK(t *testing.T) {
	clean, err := service.ValidateReviewInput(5, "  Отличная клиника  ")
	assert.NoError(t, err)
	assert.Equal(t, "Отличная клиника", clean)
}

func TestValidateReviewInput_BadRating(t *testing.T) {
	_, err := service.ValidateReviewInput(0, "нормальный текст")
	assert.ErrorIs(t, err, service.ErrReviewRating)
	_, err = service.ValidateReviewInput(6, "нормальный текст")
	assert.ErrorIs(t, err, service.ErrReviewRating)
}

func TestValidateReviewInput_TooShort(t *testing.T) {
	_, err := service.ValidateReviewInput(5, "  a  ")
	assert.ErrorIs(t, err, service.ErrReviewText)
}

func TestValidateReviewInput_TooLong(t *testing.T) {
	_, err := service.ValidateReviewInput(5, strings.Repeat("я", 1001))
	assert.ErrorIs(t, err, service.ErrReviewText)
}
