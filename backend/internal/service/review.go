// Файл: internal/service/review.go
// Назначение: бизнес-логика отзывов — валидация ввода, проверка права (завершённый визит) и создание отзыва.
package service

import (
	"errors"
	"strings"

	"beautymed/internal/model"
)

// Ошибки валидации отзывов (хендлер мапит их в 400).
var (
	ErrReviewText   = errors.New("review text must be 3..1000 characters")
	ErrReviewRating = errors.New("rating must be between 1 and 5")
)

// emptyToNil превращает пустую строку (врач/услуга не выбраны) в nil → NULL в БД.
func emptyToNil(s string) *string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return &s
}

// ReviewService — сервис отзывов поверх набора репозиториев.
type ReviewService struct{ repos *Repos }

// NewReviewService создаёт сервис отзывов.
func NewReviewService(repos *Repos) *ReviewService { return &ReviewService{repos} }

// ValidateReviewInput проверяет рейтинг (1..5) и длину текста (3..1000 рун), возвращает очищенный текст.
// Чистая функция — покрыта unit-тестами.
func ValidateReviewInput(rating int, text string) (clean string, err error) {
	if rating < 1 || rating > 5 {
		return "", ErrReviewRating
	}
	clean = strings.TrimSpace(text)
	n := len([]rune(clean))
	if n < 3 || n > 1000 {
		return "", ErrReviewText
	}
	return clean, nil
}

// Create валидирует ввод и создаёт отзыв от любого авторизованного пациента.
// Врач/услуга указываются по желанию (пустая строка → NULL).
func (s *ReviewService) Create(userID, doctorID, serviceID string, rating int, text string) (*model.Review, error) {
	clean, err := ValidateReviewInput(rating, text)
	if err != nil {
		return nil, err
	}
	return s.repos.Reviews.Create(userID, emptyToNil(doctorID), emptyToNil(serviceID), rating, clean)
}

// ListPublic — видимые отзывы с фильтрами.
func (s *ReviewService) ListPublic(doctorID, serviceID string, limit, offset int) ([]model.Review, error) {
	return s.repos.Reviews.ListPublic(doctorID, serviceID, limit, offset)
}

// Reviewable — завершённые визиты пациента для формы.
func (s *ReviewService) Reviewable(userID string) ([]model.ReviewableAppt, error) {
	return s.repos.Reviews.ReviewableAppointments(userID)
}

// ListAll — все отзывы (админ-модерация).
func (s *ReviewService) ListAll(limit, offset int) ([]model.Review, error) {
	return s.repos.Reviews.ListAll(limit, offset)
}

// SetHidden — скрыть/вернуть отзыв (админ).
func (s *ReviewService) SetHidden(id string, hidden bool) error {
	return s.repos.Reviews.SetHidden(id, hidden)
}
