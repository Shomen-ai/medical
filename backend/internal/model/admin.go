// Файл: internal/model/admin.go
// Назначение: вспомогательные DTO/модели для админ-кабинета: статистика, ячейки расписания, напоминания.
package model

import "time"

// DoctorStats хранит KPI врача за текущий календарный месяц.
type DoctorStats struct {
	AppointmentsThisMonth int     `db:"appointments_this_month" json:"appointments_this_month"`
	UniquePatients        int     `db:"unique_patients"         json:"unique_patients"`
	FilledRecordsPct      float64 `db:"filled_records_pct"      json:"filled_records_pct"`
}

// DoctorReport — строка отчёта по врачу за период: приёмы и уникальные пациенты.
type DoctorReport struct {
	DoctorID       string `db:"doctor_id"       json:"doctor_id"`
	DoctorName     string `db:"doctor_name"     json:"doctor_name"`
	SpecialtyName  string `db:"specialty_name"  json:"specialty_name"`
	Appointments   int    `db:"appointments"    json:"appointments"`
	UniquePatients int    `db:"unique_patients" json:"unique_patients"`
}

// DoctorPatient — пациент врача в отчёте: контакты, число визитов, первый/последний приём.
type DoctorPatient struct {
	PatientName  string    `db:"patient_name"  json:"patient_name"`
	PatientPhone string    `db:"patient_phone" json:"patient_phone"`
	Visits       int       `db:"visits"        json:"visits"`
	FirstVisit   time.Time `db:"first_visit"   json:"first_visit"`
	LastVisit    time.Time `db:"last_visit"    json:"last_visit"`
}

// AdminStats хранит сводные KPI для главной страницы админ-кабинета.
type AdminStats struct {
	TotalPatients     int     `db:"total_patients"      json:"total_patients"`
	ActiveDoctors     int     `db:"active_doctors"      json:"active_doctors"`
	AppointmentsMonth int     `db:"appointments_month"  json:"appointments_month"`
	RevenueQuarter    float64 `db:"revenue_quarter"     json:"revenue_quarter"`
	TopService        string  `db:"top_service"         json:"top_service"`
}

// ScheduleRow описывает одну строку расписания (врач+дата) для массовой вставки.
// ScheduleRow is the input type for BulkUpsertSchedule.
type ScheduleRow struct {
	DoctorID  string    `db:"doctor_id"`
	WorkDate  time.Time `db:"work_date"`
	StartTime string    `db:"start_time"`
	EndTime   string    `db:"end_time"`
	IsDayOff  bool      `db:"is_day_off"`
}

// ScheduleCell — одна ячейка месячного расписания (врач+дата) с информацией о наличии записей.
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

// PeriodStats хранит ключевые метрики (записи, выручка, пациенты) за произвольный период.
// PeriodStats holds key metrics for a given time range.
type PeriodStats struct {
	Appointments   int     `db:"appointments"    json:"appointments"`
	Revenue        float64 `db:"revenue"         json:"revenue"`
	UniquePatients int     `db:"unique_patients" json:"unique_patients"`
	Cancelled      int     `db:"cancelled"       json:"cancelled"`
}

// MonthlyStatPoint — точка месячной статистики (число записей и выручка) для графиков.
// MonthlyStatPoint holds appointment count and revenue for one calendar month.
type MonthlyStatPoint struct {
	Month        string  `db:"month"        json:"month"`
	Appointments int     `db:"appointments" json:"appointments"`
	Revenue      float64 `db:"revenue"      json:"revenue"`
}

// ── Модели расширенного отчёта (Excel, листы) ───────────────────────────

// ServiceReportRow — строка отчёта «Выручка по услугам»: приёмы и выручка по услуге.
type ServiceReportRow struct {
	ServiceName   string  `db:"service_name"   json:"service_name"`
	SpecialtyName string  `db:"specialty_name" json:"specialty_name"`
	Appointments  int     `db:"appointments"   json:"appointments"`
	Revenue       float64 `db:"revenue"        json:"revenue"`
}

// TimeBucketRow — число приёмов в корзине времени (n = день недели 1..7 ИЛИ час 0..23).
type TimeBucketRow struct {
	N     int `db:"n"     json:"n"`
	Count int `db:"count" json:"count"`
}

// RatingRow — средняя оценка и число отзывов по сущности (врач или услуга).
type RatingRow struct {
	Name  string  `db:"name"  json:"name"`
	Avg   float64 `db:"avg"   json:"avg"`
	Count int     `db:"count" json:"count"`
}

// RetentionRow — новые и вернувшиеся пациенты за месяц.
type RetentionRow struct {
	Month     string `db:"month"              json:"month"`
	New       int    `db:"new_patients"       json:"new_patients"`
	Returning int    `db:"returning_patients" json:"returning_patients"`
}

// ReportSummary — сводные показатели за период (лист «Сводка»).
type ReportSummary struct {
	Total          int     `db:"total"           json:"total"`
	Completed      int     `db:"completed"       json:"completed"`
	Scheduled      int     `db:"scheduled"       json:"scheduled"`
	Cancelled      int     `db:"cancelled"       json:"cancelled"`
	Rescheduled    int     `db:"rescheduled"     json:"rescheduled"`
	Revenue        float64 `db:"revenue"         json:"revenue"`
	UniquePatients int     `db:"unique_patients" json:"unique_patients"`
}

// SpecialtyReportRow — приёмы/выручка/пациенты по специальности.
type SpecialtyReportRow struct {
	SpecialtyName  string  `db:"specialty_name"  json:"specialty_name"`
	Appointments   int     `db:"appointments"    json:"appointments"`
	UniquePatients int     `db:"unique_patients" json:"unique_patients"`
	Revenue        float64 `db:"revenue"         json:"revenue"`
}

// DoctorFullRow — строка отчёта по врачу: приёмы, пациенты, выручка, средняя оценка.
type DoctorFullRow struct {
	DoctorName     string  `db:"doctor_name"     json:"doctor_name"`
	SpecialtyName  string  `db:"specialty_name"  json:"specialty_name"`
	Appointments   int     `db:"appointments"    json:"appointments"`
	UniquePatients int     `db:"unique_patients" json:"unique_patients"`
	Revenue        float64 `db:"revenue"         json:"revenue"`
	Rating         float64 `db:"rating"          json:"rating"`
}

// DailyRow — приёмы и выручка за один день периода.
type DailyRow struct {
	Day          string  `db:"day"          json:"day"`
	Appointments int     `db:"appointments" json:"appointments"`
	Revenue      float64 `db:"revenue"      json:"revenue"`
}

// LabelCountRow — пара «подпись → количество» (демография: пол, возрастные группы).
type LabelCountRow struct {
	Label string `db:"label" json:"label"`
	Count int    `db:"count" json:"count"`
}

// AppointmentReminder — минимально достаточные данные для отправки SMS-напоминания за сутки.
// AppointmentReminder holds the data needed to send a 24h reminder SMS.
type AppointmentReminder struct {
	ID           string    `db:"id"`
	PatientPhone string    `db:"patient_phone"`
	PatientName  string    `db:"patient_name"`
	DoctorName   string    `db:"doctor_name"`
	StartsAt     time.Time `db:"starts_at"`
}
