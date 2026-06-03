// Файл: internal/service/review.go
// Назначение: бизнес-логика отзывов — валидация ввода, проверка права (завершённый визит) и создание отзыва.
package service

import (
	"errors"
	"strings"

	"beautymed/internal/model"
)

// Ошибки валидации/доступа отзывов (хендлер мапит их в 400/403).
var (
	ErrReviewText   = errors.New("review text must be 3..1000 characters")
	ErrReviewRating = errors.New("rating must be between 1 and 5")
	ErrNotEligible  = errors.New("a completed visit is required to leave a review")
)

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

// Create проверяет ввод и право пациента (завершённый визит), затем создаёт отзыв.
func (s *ReviewService) Create(userID, appointmentID string, rating int, text string) (*model.Review, error) {
	clean, err := ValidateReviewInput(rating, text)
	if err != nil {
		return nil, err
	}
	doctorID, serviceID, ok, err := s.repos.Reviews.AppointmentForReview(userID, appointmentID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrNotEligible
	}
	return s.repos.Reviews.Create(userID, appointmentID, doctorID, serviceID, rating, clean)
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
