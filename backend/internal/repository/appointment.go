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

// DoctorStats returns aggregate stats for the current calendar month.
func (r *AppointmentRepo) DoctorStats(doctorID string) (*model.DoctorStats, error) {
	var s model.DoctorStats
	err := r.db.Get(&s, `
		WITH monthly AS (
			SELECT a.patient_id, a.id, ar.id AS record_id
			FROM appointments a
			LEFT JOIN appointment_records ar
				ON ar.appointment_id = a.id AND ar.is_draft = false
			WHERE a.doctor_id = $1
			  AND a.starts_at >= date_trunc('month', NOW())
			  AND a.starts_at <  date_trunc('month', NOW()) + INTERVAL '1 month'
		)
		SELECT
			COUNT(*)                                          AS appointments_this_month,
			COUNT(DISTINCT patient_id)                        AS unique_patients,
			CASE WHEN COUNT(*) = 0 THEN 0.0
			     ELSE ROUND(COUNT(record_id)::numeric / COUNT(*) * 100, 1)
			END                                               AS filled_records_pct
		FROM monthly`, doctorID)
	return &s, err
}

// ListScheduleByDoctor returns the month's schedule rows for a single doctor.
func (r *AppointmentRepo) ListScheduleByDoctor(doctorID string, year, month int) ([]map[string]interface{}, error) {
	rows, err := r.db.Queryx(`
		SELECT work_date::text, start_time::text, end_time::text, is_day_off
		FROM doctor_schedules
		WHERE doctor_id = $1
		  AND EXTRACT(YEAR  FROM work_date) = $2
		  AND EXTRACT(MONTH FROM work_date) = $3
		ORDER BY work_date`, doctorID, year, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []map[string]interface{}
	for rows.Next() {
		m := make(map[string]interface{})
		if err := rows.MapScan(m); err != nil {
			return nil, err
		}
		for k, v := range m {
			if b, ok := v.([]byte); ok {
				m[k] = string(b)
			}
		}
		result = append(result, m)
	}
	return result, rows.Err()
}

// ListCompletedByPatientYear returns completed appointments for a patient in a given year.
func (r *AppointmentRepo) ListCompletedByPatientYear(patientID string, year int) ([]model.Appointment, error) {
	var as []model.Appointment
	err := r.db.Select(&as, `
		SELECT a.*,
		       d.full_name AS doctor_name,
		       s.name      AS service_name
		FROM appointments a
		JOIN doctors  d ON d.id = a.doctor_id
		JOIN services s ON s.id = a.service_id
		WHERE a.patient_id = $1
		  AND a.status     = 'completed'
		  AND EXTRACT(YEAR FROM a.starts_at) = $2
		ORDER BY a.starts_at`, patientID, year)
	return as, err
}
