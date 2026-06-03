// Файл: internal/model/review.go
// Назначение: модель анонимного отзыва пациента (привязан к завершённому визиту) с присоединёнными именами врача/услуги для отображения.
package model

import "time"

// Review — отзыв пациента. user_id/appointment_id хранятся для аудита, но наружу не отдаются.
type Review struct {
	ID            string    `db:"id"             json:"id"`
	UserID        string    `db:"user_id"        json:"-"`
	AppointmentID string    `db:"appointment_id" json:"-"`
	DoctorID      string    `db:"doctor_id"      json:"doctor_id"`
	ServiceID     string    `db:"service_id"     json:"service_id"`
	Rating        int       `db:"rating"         json:"rating"`
	Text          string    `db:"text"           json:"text"`
	IsHidden      bool      `db:"is_hidden"      json:"is_hidden"`
	CreatedAt     time.Time `db:"created_at"     json:"created_at"`
	// Присоединяются при выборке для отображения:
	DoctorName    string `db:"doctor_name"    json:"doctor_name,omitempty"`
	SpecialtyName string `db:"specialty_name" json:"specialty_name,omitempty"`
	ServiceName   string `db:"service_name"   json:"service_name,omitempty"`
}

// ReviewableAppt — завершённый визит пациента, по которому можно оставить отзыв (для формы).
type ReviewableAppt struct {
	AppointmentID string    `db:"appointment_id" json:"appointment_id"`
	DoctorName    string    `db:"doctor_name"    json:"doctor_name"`
	ServiceName   string    `db:"service_name"   json:"service_name"`
	StartsAt      time.Time `db:"starts_at"      json:"starts_at"`
}
