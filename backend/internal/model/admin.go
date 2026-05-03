package model

import "time"

type DoctorStats struct {
	AppointmentsThisMonth int     `db:"appointments_this_month" json:"appointments_this_month"`
	UniquePatients        int     `db:"unique_patients"         json:"unique_patients"`
	FilledRecordsPct      float64 `db:"filled_records_pct"      json:"filled_records_pct"`
}

type AdminStats struct {
	TotalPatients     int     `db:"total_patients"      json:"total_patients"`
	ActiveDoctors     int     `db:"active_doctors"      json:"active_doctors"`
	AppointmentsMonth int     `db:"appointments_month"  json:"appointments_month"`
	RevenueQuarter    float64 `db:"revenue_quarter"     json:"revenue_quarter"`
	TopService        string  `db:"top_service"         json:"top_service"`
}

// ScheduleRow is the input type for BulkUpsertSchedule.
type ScheduleRow struct {
	DoctorID  string    `db:"doctor_id"`
	WorkDate  time.Time `db:"work_date"`
	StartTime string    `db:"start_time"`
	EndTime   string    `db:"end_time"`
	IsDayOff  bool      `db:"is_day_off"`
}

// ScheduleCell is returned by ListScheduleForMonth — one cell per (doctor, date).
type ScheduleCell struct {
	DoctorID        string `db:"doctor_id"        json:"doctor_id"`
	DoctorName      string `db:"doctor_name"      json:"doctor_name"`
	SpecialtyID     string `db:"specialty_id"     json:"specialty_id"`
	WorkDate        string `db:"work_date"        json:"work_date"`
	StartTime       string `db:"start_time"       json:"start_time"`
	EndTime         string `db:"end_time"         json:"end_time"`
	IsDayOff        bool   `db:"is_day_off"       json:"is_day_off"`
	HasAppointments bool   `db:"has_appointments" json:"has_appointments"`
}

// AppointmentReminder holds the data needed to send a 24h reminder SMS.
type AppointmentReminder struct {
	ID           string    `db:"id"`
	PatientPhone string    `db:"patient_phone"`
	PatientName  string    `db:"patient_name"`
	DoctorName   string    `db:"doctor_name"`
	StartsAt     time.Time `db:"starts_at"`
}
