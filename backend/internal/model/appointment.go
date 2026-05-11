package model

import "time"

type Appointment struct {
	ID          string    `db:"id"           json:"id"`
	PatientID   string    `db:"patient_id"   json:"patient_id"`
	DoctorID    string    `db:"doctor_id"    json:"doctor_id"`
	ServiceID   string    `db:"service_id"   json:"service_id"`
	PromoCodeID *string   `db:"promo_code_id" json:"promo_code_id,omitempty"`
	StartsAt    time.Time `db:"starts_at"    json:"starts_at"`
	EndsAt      time.Time `db:"ends_at"      json:"ends_at"`
	Status      string    `db:"status"       json:"status"`
	FinalPrice  float64   `db:"final_price"  json:"final_price"`
	CreatedBy   string    `db:"created_by"   json:"created_by"`
	CreatedAt   time.Time `db:"created_at"   json:"created_at"`

	// Joined fields for display
	DoctorName   string  `db:"doctor_name"   json:"doctor_name,omitempty"`
	ServiceName  string  `db:"service_name"  json:"service_name,omitempty"`
	PatientName  string  `db:"patient_name"  json:"patient_name,omitempty"`
	PatientPhone string  `db:"patient_phone" json:"patient_phone,omitempty"`
}

type AppointmentRecord struct {
	ID              string    `db:"id"               json:"id"`
	AppointmentID   string    `db:"appointment_id"   json:"appointment_id"`
	Complaints      string    `db:"complaints"       json:"complaints"`
	Diagnosis       string    `db:"diagnosis"        json:"diagnosis"`
	Prescription    *string   `db:"prescription"     json:"prescription,omitempty"`
	Recommendations *string   `db:"recommendations"  json:"recommendations,omitempty"`
	IsDraft         bool      `db:"is_draft"         json:"is_draft"`
	CreatedAt       time.Time `db:"created_at"       json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"       json:"updated_at"`
}

type TimeSlot struct {
	StartsAt string `json:"starts_at"`
	EndsAt   string `json:"ends_at"`
}

type PromoCode struct {
	ID          string    `db:"id"           json:"id"`
	Code        string    `db:"code"         json:"code"`
	DiscountPct int       `db:"discount_pct" json:"discount_pct"`
	MaxUses     *int      `db:"max_uses"     json:"max_uses,omitempty"`
	UsedCount   int       `db:"used_count"   json:"used_count"`
	ValidFrom   string    `db:"valid_from"   json:"valid_from"`
	ValidUntil  *string   `db:"valid_until"  json:"valid_until,omitempty"`
	IsActive    bool      `db:"is_active"    json:"is_active"`
	CreatedAt   time.Time `db:"created_at"   json:"created_at"`
}
