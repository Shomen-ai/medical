package repository

import (
	"time"

	"beautymed/internal/model"

	"github.com/jmoiron/sqlx"
)

type AppointmentRepo struct{ db *sqlx.DB }

func NewAppointmentRepo(db *sqlx.DB) *AppointmentRepo { return &AppointmentRepo{db} }

func (r *AppointmentRepo) Create(a *model.Appointment) error {
	return r.db.Get(a, `
		INSERT INTO appointments
			(patient_id, doctor_id, service_id, promo_code_id, starts_at, ends_at, final_price, created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		RETURNING *`,
		a.PatientID, a.DoctorID, a.ServiceID, a.PromoCodeID,
		a.StartsAt, a.EndsAt, a.FinalPrice, a.CreatedBy)
}

func (r *AppointmentRepo) FindByID(id string) (*model.Appointment, error) {
	var a model.Appointment
	err := r.db.Get(&a, `
		SELECT a.*,
		       d.full_name AS doctor_name,
		       s.name      AS service_name,
		       u.full_name AS patient_name,
		       u.phone     AS patient_phone
		FROM appointments a
		JOIN doctors  d ON d.id = a.doctor_id
		JOIN services s ON s.id = a.service_id
		JOIN users    u ON u.id = a.patient_id
		WHERE a.id = $1`, id)
	return &a, err
}

func (r *AppointmentRepo) TakenSlots(doctorID string, date time.Time) ([]time.Time, error) {
	var slots []time.Time
	err := r.db.Select(&slots, `
		SELECT starts_at FROM appointments
		WHERE doctor_id = $1
		  AND DATE(starts_at AT TIME ZONE 'UTC') = $2
		  AND status NOT IN ('cancelled','rescheduled')`,
		doctorID, date.Format("2006-01-02"))
	return slots, err
}

func (r *AppointmentRepo) ListByPatient(patientID string) ([]model.Appointment, error) {
	var as []model.Appointment
	err := r.db.Select(&as, `
		SELECT a.*,
		       d.full_name AS doctor_name,
		       s.name      AS service_name
		FROM appointments a
		JOIN doctors  d ON d.id = a.doctor_id
		JOIN services s ON s.id = a.service_id
		WHERE a.patient_id = $1
		ORDER BY a.starts_at DESC`, patientID)
	return as, err
}

func (r *AppointmentRepo) ListByDoctor(doctorID string, date time.Time) ([]model.Appointment, error) {
	var as []model.Appointment
	err := r.db.Select(&as, `
		SELECT a.*,
		       s.name      AS service_name,
		       u.full_name AS patient_name,
		       u.phone     AS patient_phone
		FROM appointments a
		JOIN services s ON s.id = a.service_id
		JOIN users    u ON u.id = a.patient_id
		WHERE a.doctor_id = $1
		  AND DATE(a.starts_at AT TIME ZONE 'UTC') = $2
		  AND a.status NOT IN ('cancelled','rescheduled')
		ORDER BY a.starts_at`, doctorID, date.Format("2006-01-02"))
	return as, err
}

func (r *AppointmentRepo) UpdateStatus(id, status string) error {
	_, err := r.db.Exec(
		`UPDATE appointments SET status=$1 WHERE id=$2`, status, id)
	return err
}

func (r *AppointmentRepo) Reschedule(id string, startsAt, endsAt time.Time) error {
	_, err := r.db.Exec(`
		UPDATE appointments SET starts_at=$1, ends_at=$2, status='rescheduled'
		WHERE id=$3`, startsAt, endsAt, id)
	return err
}

func (r *AppointmentRepo) GetRecord(appointmentID string) (*model.AppointmentRecord, error) {
	var rec model.AppointmentRecord
	err := r.db.Get(&rec, `SELECT * FROM appointment_records WHERE appointment_id=$1`, appointmentID)
	return &rec, err
}

func (r *AppointmentRepo) UpsertRecord(rec *model.AppointmentRecord) error {
	_, err := r.db.Exec(`
		INSERT INTO appointment_records
			(appointment_id, complaints, diagnosis, prescription, recommendations, is_draft)
		VALUES ($1,$2,$3,$4,$5,$6)
		ON CONFLICT (appointment_id) DO UPDATE SET
			complaints=$2, diagnosis=$3, prescription=$4,
			recommendations=$5, is_draft=$6, updated_at=now()`,
		rec.AppointmentID, rec.Complaints, rec.Diagnosis,
		rec.Prescription, rec.Recommendations, rec.IsDraft)
	return err
}

type DoctorSchedule struct {
	StartTime string `db:"start_time"`
	EndTime   string `db:"end_time"`
	IsDayOff  bool   `db:"is_day_off"`
}

func (r *AppointmentRepo) GetSchedule(doctorID, date string) (*DoctorSchedule, error) {
	var s DoctorSchedule
	err := r.db.Get(&s, `
		SELECT start_time::text, end_time::text, is_day_off
		FROM doctor_schedules
		WHERE doctor_id=$1 AND work_date=$2`, doctorID, date)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *AppointmentRepo) FindPromoCode(code string) (*model.PromoCode, error) {
	var pc model.PromoCode
	err := r.db.Get(&pc, `
		SELECT * FROM promo_codes
		WHERE code=$1 AND is_active=true
		  AND valid_from <= CURRENT_DATE
		  AND (valid_until IS NULL OR valid_until >= CURRENT_DATE)
		  AND (max_uses IS NULL OR used_count < max_uses)`, code)
	return &pc, err
}

func (r *AppointmentRepo) IncrementPromoUsage(id string) error {
	_, err := r.db.Exec(`UPDATE promo_codes SET used_count=used_count+1 WHERE id=$1`, id)
	return err
}
