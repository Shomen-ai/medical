// Файл: internal/repository/review.go
// Назначение: SQL-доступ к таблице reviews — публичный список с фильтрами, создание, выборка завершённых визитов для формы и модерация.
package repository

import (
	"database/sql"
	"errors"

	"beautymed/internal/model"

	"github.com/jmoiron/sqlx"
)

// ReviewRepo — репозиторий отзывов пациентов.
type ReviewRepo struct{ db *sqlx.DB }

// NewReviewRepo создаёт ReviewRepo на базе пула sqlx.
func NewReviewRepo(db *sqlx.DB) *ReviewRepo { return &ReviewRepo{db} }

// LEFT JOIN — отзыв может быть без привязки к врачу/услуге (тогда имена пустые).
const reviewSelect = `
	SELECT r.*,
	       COALESCE(d.full_name, '') AS doctor_name,
	       COALESCE(s.name, '')      AS specialty_name,
	       COALESCE(sv.name, '')     AS service_name
	FROM reviews r
	LEFT JOIN doctors d     ON d.id  = r.doctor_id
	LEFT JOIN specialties s ON s.id  = d.specialty_id
	LEFT JOIN services sv   ON sv.id = r.service_id`

// ListPublic возвращает видимые отзывы (не скрытые), опционально фильтруя по врачу и услуге,
// новые сверху, с присоединёнными именами врача/специальности/услуги.
func (r *ReviewRepo) ListPublic(doctorID, serviceID string, limit, offset int) ([]model.Review, error) {
	var rs []model.Review
	err := r.db.Select(&rs, reviewSelect+`
		WHERE r.is_hidden = false
		  AND ($1 = '' OR r.doctor_id  = $1::uuid)
		  AND ($2 = '' OR r.service_id = $2::uuid)
		ORDER BY r.created_at DESC
		LIMIT $3 OFFSET $4`, doctorID, serviceID, limit, offset)
	return rs, err
}

// ListAll возвращает все отзывы (включая скрытые) для админ-модерации.
func (r *ReviewRepo) ListAll(limit, offset int) ([]model.Review, error) {
	var rs []model.Review
	err := r.db.Select(&rs, reviewSelect+`
		ORDER BY r.created_at DESC
		LIMIT $1 OFFSET $2`, limit, offset)
	return rs, err
}

// ReviewableAppointments возвращает завершённые визиты пациента (для выпадающего списка формы).
func (r *ReviewRepo) ReviewableAppointments(userID string) ([]model.ReviewableAppt, error) {
	var as []model.ReviewableAppt
	err := r.db.Select(&as, `
		SELECT a.id AS appointment_id, d.full_name AS doctor_name, sv.name AS service_name, a.starts_at
		FROM appointments a
		JOIN doctors d   ON d.id  = a.doctor_id
		JOIN services sv ON sv.id = a.service_id
		WHERE a.patient_id = $1 AND a.status = 'completed'
		ORDER BY a.starts_at DESC`, userID)
	return as, err
}

// AppointmentForReview проверяет, что визит принадлежит пользователю и завершён,
// и возвращает его doctor_id/service_id для записи в отзыв.
func (r *ReviewRepo) AppointmentForReview(userID, appointmentID string) (doctorID, serviceID string, ok bool, err error) {
	var row struct {
		DoctorID  string `db:"doctor_id"`
		ServiceID string `db:"service_id"`
	}
	err = r.db.Get(&row, `
		SELECT doctor_id, service_id FROM appointments
		WHERE id = $1 AND patient_id = $2 AND status = 'completed'`, appointmentID, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", "", false, nil
	}
	if err != nil {
		return "", "", false, err
	}
	return row.DoctorID, row.ServiceID, true, nil
}

// Create вставляет новый отзыв (врач/услуга — опционально, nil → NULL) и возвращает его.
func (r *ReviewRepo) Create(userID string, doctorID, serviceID *string, rating int, text string) (*model.Review, error) {
	var rv model.Review
	err := r.db.Get(&rv, `
		INSERT INTO reviews (user_id, doctor_id, service_id, rating, text)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, appointment_id, doctor_id, service_id, rating, text, is_hidden, created_at`,
		userID, doctorID, serviceID, rating, text)
	return &rv, err
}

// SetHidden скрывает или возвращает отзыв (мягкая модерация).
func (r *ReviewRepo) SetHidden(id string, hidden bool) error {
	_, err := r.db.Exec(`UPDATE reviews SET is_hidden = $1 WHERE id = $2`, hidden, id)
	return err
}
