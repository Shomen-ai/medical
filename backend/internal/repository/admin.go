// Файл: internal/repository/admin.go
// Назначение: SQL-запросы для админ-кабинета — статистика, расписание, промокоды, KPI, напоминания.
package repository

import (
	"fmt"
	"strings"
	"time"

	"beautymed/internal/model"

	"github.com/jmoiron/sqlx"
)

// AdminRepo — репозиторий админских операций над всеми записями и расписанием.
type AdminRepo struct{ db *sqlx.DB }

// NewAdminRepo создаёт AdminRepo, обёрнутый вокруг переданного пула sqlx.
func NewAdminRepo(db *sqlx.DB) *AdminRepo { return &AdminRepo{db} }

// ListAllAppointments возвращает все записи на приём с опциональными фильтрами по статусу/врачу/дате.
// ListAllAppointments returns all appointments with optional filters.
// status: "" | "scheduled" | "completed" | "cancelled" | "rescheduled"
// doctorID: "" means all doctors
// date: nil means all dates
func (r *AdminRepo) ListAllAppointments(status, doctorID string, date *time.Time) ([]model.Appointment, error) {
	q := `
		SELECT a.*,
		       d.full_name AS doctor_name,
		       s.name      AS service_name,
		       u.full_name AS patient_name,
		       u.phone     AS patient_phone
		FROM appointments a
		JOIN doctors  d ON d.id = a.doctor_id
		JOIN services s ON s.id = a.service_id
		JOIN users    u ON u.id = a.patient_id
		WHERE 1=1`
	args := []interface{}{}
	i := 1
	if status != "" {
		q += fmt.Sprintf(" AND a.status = $%d", i)
		args = append(args, status)
		i++
	}
	if doctorID != "" {
		q += fmt.Sprintf(" AND a.doctor_id = $%d", i)
		args = append(args, doctorID)
		i++
	}
	if date != nil {
		q += fmt.Sprintf(" AND a.starts_at::date = $%d", i)
		args = append(args, date.Format("2006-01-02"))
		i++
	}
	_ = i
	q += " ORDER BY a.starts_at DESC"
	var as []model.Appointment
	return as, r.db.Select(&as, q, args...)
}

// PeriodStats возвращает агрегированную статистику записей и выручки за интервал [from, to).
// PeriodStats returns appointment counts and revenue for the given time range.
func (r *AdminRepo) PeriodStats(from, to time.Time) (*model.PeriodStats, error) {
	var s model.PeriodStats
	err := r.db.Get(&s, `
		SELECT
		  COUNT(*)                                                   AS appointments,
		  COALESCE(SUM(final_price) FILTER (WHERE status='completed'), 0) AS revenue,
		  COUNT(DISTINCT patient_id)                                 AS unique_patients,
		  COUNT(*) FILTER (WHERE status='cancelled')                 AS cancelled
		FROM appointments
		WHERE starts_at >= $1 AND starts_at < $2`, from, to)
	return &s, err
}

// ByDoctorStats возвращает по каждому активному врачу число приёмов и уникальных
// пациентов за интервал [from, to). Отменённые/перенесённые приёмы не считаются;
// врачи без приёмов попадают в отчёт с нулями (LEFT JOIN).
func (r *AdminRepo) ByDoctorStats(from, to time.Time) ([]model.DoctorReport, error) {
	var rows []model.DoctorReport
	err := r.db.Select(&rows, `
		SELECT d.id        AS doctor_id,
		       d.full_name AS doctor_name,
		       s.name      AS specialty_name,
		       COUNT(a.id) FILTER (WHERE a.status NOT IN ('cancelled','rescheduled'))                  AS appointments,
		       COUNT(DISTINCT a.patient_id) FILTER (WHERE a.status NOT IN ('cancelled','rescheduled')) AS unique_patients
		FROM doctors d
		JOIN specialties s ON s.id = d.specialty_id
		LEFT JOIN appointments a
		       ON a.doctor_id = d.id
		      AND a.starts_at >= $1 AND a.starts_at < $2
		WHERE d.is_active = true
		GROUP BY d.id, d.full_name, s.name
		ORDER BY appointments DESC, d.full_name`, from, to)
	return rows, err
}

// Revenue возвращает суммарную выручку (final_price) по завершённым приёмам за интервал [from, to).
// Revenue returns the total final_price for completed appointments in [from, to).
func (r *AdminRepo) Revenue(from, to time.Time) (float64, error) {
	var total float64
	err := r.db.Get(&total, `
		SELECT COALESCE(SUM(final_price), 0)
		FROM appointments
		WHERE status = 'completed'
		  AND starts_at >= $1 AND starts_at < $2`, from, to)
	return total, err
}

// OverallStats возвращает сводные KPI для главной страницы админ-кабинета.
// OverallStats returns aggregate KPIs for the admin stats page.
func (r *AdminRepo) OverallStats() (*model.AdminStats, error) {
	var s model.AdminStats
	err := r.db.Get(&s, `
		SELECT
		  (SELECT COUNT(DISTINCT patient_id) FROM appointments)           AS total_patients,
		  (SELECT COUNT(*) FROM doctors WHERE is_active = true)           AS active_doctors,
		  (SELECT COUNT(*) FROM appointments
		     WHERE starts_at >= date_trunc('month', NOW())
		       AND starts_at <  date_trunc('month', NOW()) + INTERVAL '1 month'
		  )                                                               AS appointments_month,
		  (SELECT COALESCE(SUM(final_price), 0) FROM appointments
		     WHERE status = 'completed'
		       AND starts_at >= date_trunc('quarter', NOW())
		       AND starts_at <  date_trunc('quarter', NOW()) + INTERVAL '3 months'
		  )                                                               AS revenue_quarter,
		  (SELECT COALESCE(s.name, 'Нет данных')
		     FROM appointments a
		     JOIN services s ON s.id = a.service_id
		     WHERE a.starts_at >= date_trunc('month', NOW())
		       AND a.starts_at <  date_trunc('month', NOW()) + INTERVAL '1 month'
		     GROUP BY s.name
		     ORDER BY COUNT(*) DESC
		     LIMIT 1
		  )                                                               AS top_service`)
	return &s, err
}

// TodayKPI возвращает на сегодня: число записей, выручку и количество свободных слотов.
// TodayKPI returns appointment count, revenue and free slot count for today.
func (r *AdminRepo) TodayKPI() (int, float64, int, error) {
	var count int
	var revenue float64
	if err := r.db.QueryRow(`
		SELECT COUNT(*), COALESCE(SUM(final_price) FILTER (WHERE status='completed'), 0)
		FROM appointments
		WHERE starts_at::date = CURRENT_DATE`).Scan(&count, &revenue); err != nil {
		return 0, 0, 0, err
	}
	var totalSlots, booked int
	r.db.QueryRow(`
		SELECT
		  (SELECT COUNT(*) FROM doctor_schedules WHERE work_date = CURRENT_DATE AND is_day_off = false),
		  (SELECT COUNT(*) FROM appointments WHERE starts_at::date = CURRENT_DATE AND status = 'scheduled')
	`).Scan(&totalSlots, &booked)
	free := totalSlots - booked
	if free < 0 {
		free = 0
	}
	return count, revenue, free, nil
}

// ListTomorrowAppointments возвращает все назначенные на завтра приёмы с контактами для SMS-напоминаний.
// ListTomorrowAppointments returns all scheduled appointments for tomorrow with contact info.
func (r *AdminRepo) ListTomorrowAppointments() ([]model.AppointmentReminder, error) {
	var reminders []model.AppointmentReminder
	err := r.db.Select(&reminders, `
		SELECT a.id, u.phone AS patient_phone, u.full_name AS patient_name,
		       d.full_name AS doctor_name, a.starts_at
		FROM appointments a
		JOIN users   u ON u.id = a.patient_id
		JOIN doctors d ON d.id = a.doctor_id
		WHERE a.starts_at::date = CURRENT_DATE + INTERVAL '1 day'
		  AND a.status = 'scheduled'
		ORDER BY a.starts_at`)
	return reminders, err
}

// ListScheduleForMonth возвращает ячейки расписания за указанный месяц с инфо о врачах и наличии записей.
// ListScheduleForMonth returns all schedule cells for the given month, joined with doctor info.
func (r *AdminRepo) ListScheduleForMonth(year, month int, specialtyID string) ([]model.ScheduleCell, error) {
	q := `
		SELECT ds.doctor_id, d.full_name AS doctor_name, d.specialty_id,
		       ds.work_date::text, ds.start_time::text, ds.end_time::text, ds.is_day_off,
		       EXISTS(
		         SELECT 1 FROM appointments a
		         WHERE a.doctor_id = ds.doctor_id
		           AND a.starts_at::date = ds.work_date
		           AND a.status = 'scheduled'
		       ) AS has_appointments
		FROM doctor_schedules ds
		JOIN doctors d ON d.id = ds.doctor_id
		WHERE EXTRACT(YEAR  FROM ds.work_date) = $1
		  AND EXTRACT(MONTH FROM ds.work_date) = $2`
	args := []interface{}{year, month}
	if specialtyID != "" {
		q += " AND d.specialty_id = $3"
		args = append(args, specialtyID)
	}
	q += " ORDER BY ds.work_date, d.full_name"
	var cells []model.ScheduleCell
	return cells, r.db.Select(&cells, q, args...)
}

// BulkUpsertSchedule массово вставляет или обновляет строки расписания (по конфликту doctor_id+work_date).
// BulkUpsertSchedule inserts or updates schedule rows in bulk.
func (r *AdminRepo) BulkUpsertSchedule(rows []model.ScheduleRow) error {
	if len(rows) == 0 {
		return nil
	}
	vals := make([]string, 0, len(rows))
	args := make([]interface{}, 0, len(rows)*5)
	i := 1
	for _, row := range rows {
		vals = append(vals, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d)", i, i+1, i+2, i+3, i+4))
		args = append(args, row.DoctorID, row.WorkDate.Format("2006-01-02"),
			row.StartTime, row.EndTime, row.IsDayOff)
		i += 5
	}
	q := `INSERT INTO doctor_schedules (doctor_id, work_date, start_time, end_time, is_day_off)
	      VALUES ` + strings.Join(vals, ",") + `
	      ON CONFLICT (doctor_id, work_date) DO UPDATE SET
	        start_time = EXCLUDED.start_time,
	        end_time   = EXCLUDED.end_time,
	        is_day_off = EXCLUDED.is_day_off`
	_, err := r.db.Exec(q, args...)
	return err
}

// UpsertScheduleCell вставляет или обновляет одну ячейку расписания для пары (врач, дата).
// UpsertScheduleCell updates a single doctor/date schedule cell.
func (r *AdminRepo) UpsertScheduleCell(doctorID, date, startTime, endTime string, isDayOff bool) error {
	_, err := r.db.Exec(`
		INSERT INTO doctor_schedules (doctor_id, work_date, start_time, end_time, is_day_off)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (doctor_id, work_date) DO UPDATE SET
		  start_time = EXCLUDED.start_time,
		  end_time   = EXCLUDED.end_time,
		  is_day_off = EXCLUDED.is_day_off`,
		doctorID, date, startTime, endTime, isDayOff)
	return err
}

// CreatePromoCode создаёт новый промокод и записывает в pc.ID сгенерированный UUID.
// CreatePromoCode inserts a new promo code.
func (r *AdminRepo) CreatePromoCode(pc *model.PromoCode) error {
	return r.db.QueryRow(`
		INSERT INTO promo_codes (code, discount_pct, max_uses, valid_from, valid_until, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`,
		pc.Code, pc.DiscountPct, pc.MaxUses, pc.ValidFrom, pc.ValidUntil, pc.IsActive,
	).Scan(&pc.ID)
}

// ListPromos возвращает все промокоды, упорядоченные по дате создания (новые сверху).
// ListPromos returns all promo codes ordered by creation date desc.
func (r *AdminRepo) ListPromos() ([]model.PromoCode, error) {
	var pcs []model.PromoCode
	return pcs, r.db.Select(&pcs, `
		SELECT * FROM promo_codes ORDER BY created_at DESC`)
}

// MonthlyStats возвращает помесячные метрики (число записей, выручка) за последние 6 календарных месяцев.
// MonthlyStats returns appointment counts and revenue for the last 6 calendar months.
func (r *AdminRepo) MonthlyStats() ([]model.MonthlyStatPoint, error) {
	var pts []model.MonthlyStatPoint
	err := r.db.Select(&pts, `
		SELECT
		  TO_CHAR(DATE_TRUNC('month', starts_at), 'YYYY-MM') AS month,
		  COUNT(*)                                            AS appointments,
		  COALESCE(SUM(final_price) FILTER (WHERE status = 'completed'), 0) AS revenue
		FROM appointments
		WHERE starts_at >= DATE_TRUNC('month', NOW()) - INTERVAL '5 months'
		GROUP BY DATE_TRUNC('month', starts_at)
		ORDER BY DATE_TRUNC('month', starts_at)`)
	return pts, err
}

// MonthlyStatsRange возвращает помесячные метрики за интервал [from, to) — строку на КАЖДЫЙ
// месяц периода, включая месяцы без записей (через generate_series), для графика по всему периоду.
func (r *AdminRepo) MonthlyStatsRange(from, to time.Time) ([]model.MonthlyStatPoint, error) {
	var pts []model.MonthlyStatPoint
	err := r.db.Select(&pts, `
		SELECT TO_CHAR(m, 'YYYY-MM') AS month,
		       COALESCE(c.appointments, 0) AS appointments,
		       COALESCE(c.revenue, 0)      AS revenue
		FROM generate_series(
		         DATE_TRUNC('month', $1::timestamptz),
		         DATE_TRUNC('month', ($2::timestamptz - INTERVAL '1 day')),
		         INTERVAL '1 month') AS m
		LEFT JOIN (
		    SELECT DATE_TRUNC('month', starts_at) AS mo,
		           COUNT(*) AS appointments,
		           COALESCE(SUM(final_price) FILTER (WHERE status = 'completed'), 0) AS revenue
		    FROM appointments
		    WHERE starts_at >= $1 AND starts_at < $2
		    GROUP BY 1
		) c ON c.mo = m
		ORDER BY m`, from, to)
	return pts, err
}

// ── Расширенный отчёт (4 листа Excel) ────────────────────────────────────

// ReportByService — приёмы и выручка по каждой услуге за интервал [from, to).
func (r *AdminRepo) ReportByService(from, to time.Time) ([]model.ServiceReportRow, error) {
	var rows []model.ServiceReportRow
	err := r.db.Select(&rows, `
		SELECT sv.name AS service_name, sp.name AS specialty_name,
		       COUNT(a.id) FILTER (WHERE a.status NOT IN ('cancelled','rescheduled'))      AS appointments,
		       COALESCE(SUM(a.final_price) FILTER (WHERE a.status = 'completed'), 0)        AS revenue
		FROM services sv
		JOIN specialties sp ON sp.id = sv.specialty_id
		LEFT JOIN appointments a ON a.service_id = sv.id AND a.starts_at >= $1 AND a.starts_at < $2
		GROUP BY sv.id, sv.name, sp.name
		ORDER BY revenue DESC, appointments DESC`, from, to)
	return rows, err
}

// LoadByWeekday — число приёмов по дням недели (n: 1=Пн..7=Вс) за интервал.
func (r *AdminRepo) LoadByWeekday(from, to time.Time) ([]model.TimeBucketRow, error) {
	var rows []model.TimeBucketRow
	err := r.db.Select(&rows, `
		SELECT EXTRACT(ISODOW FROM starts_at)::int AS n, COUNT(*) AS count
		FROM appointments
		WHERE status <> 'cancelled' AND starts_at >= $1 AND starts_at < $2
		GROUP BY 1 ORDER BY 1`, from, to)
	return rows, err
}

// LoadByHour — число приёмов по часам начала (n: 0..23) за интервал.
func (r *AdminRepo) LoadByHour(from, to time.Time) ([]model.TimeBucketRow, error) {
	var rows []model.TimeBucketRow
	err := r.db.Select(&rows, `
		SELECT EXTRACT(HOUR FROM starts_at)::int AS n, COUNT(*) AS count
		FROM appointments
		WHERE status <> 'cancelled' AND starts_at >= $1 AND starts_at < $2
		GROUP BY 1 ORDER BY 1`, from, to)
	return rows, err
}

// RatingsByDoctor — средняя оценка и число отзывов по врачу за интервал (по дате отзыва).
func (r *AdminRepo) RatingsByDoctor(from, to time.Time) ([]model.RatingRow, error) {
	var rows []model.RatingRow
	err := r.db.Select(&rows, `
		SELECT d.full_name AS name, ROUND(AVG(r.rating)::numeric, 2)::float8 AS avg, COUNT(*) AS count
		FROM reviews r
		JOIN doctors d ON d.id = r.doctor_id
		WHERE r.is_hidden = false AND r.created_at >= $1 AND r.created_at < $2
		GROUP BY d.id, d.full_name
		ORDER BY avg DESC, count DESC`, from, to)
	return rows, err
}

// RatingsByService — средняя оценка и число отзывов по услуге за интервал.
func (r *AdminRepo) RatingsByService(from, to time.Time) ([]model.RatingRow, error) {
	var rows []model.RatingRow
	err := r.db.Select(&rows, `
		SELECT sv.name AS name, ROUND(AVG(r.rating)::numeric, 2)::float8 AS avg, COUNT(*) AS count
		FROM reviews r
		JOIN services sv ON sv.id = r.service_id
		WHERE r.is_hidden = false AND r.created_at >= $1 AND r.created_at < $2
		GROUP BY sv.id, sv.name
		ORDER BY avg DESC, count DESC`, from, to)
	return rows, err
}

// Retention — новые и вернувшиеся пациенты помесячно за интервал. Новый = первый в истории
// (не отменённый) приём пациента пришёлся на этот месяц; вернувшийся — первый приём был раньше.
func (r *AdminRepo) Retention(from, to time.Time) ([]model.RetentionRow, error) {
	var rows []model.RetentionRow
	err := r.db.Select(&rows, `
		WITH firsts AS (
		    SELECT patient_id, MIN(starts_at) AS first_at
		    FROM appointments WHERE status <> 'cancelled'
		    GROUP BY patient_id
		),
		month_patients AS (
		    SELECT DISTINCT DATE_TRUNC('month', a.starts_at) AS mo, a.patient_id
		    FROM appointments a
		    WHERE a.status <> 'cancelled' AND a.starts_at >= $1 AND a.starts_at < $2
		)
		SELECT TO_CHAR(mp.mo, 'YYYY-MM') AS month,
		       COUNT(*) FILTER (WHERE DATE_TRUNC('month', f.first_at) = mp.mo) AS new_patients,
		       COUNT(*) FILTER (WHERE DATE_TRUNC('month', f.first_at) < mp.mo) AS returning_patients
		FROM month_patients mp
		JOIN firsts f ON f.patient_id = mp.patient_id
		GROUP BY mp.mo
		ORDER BY mp.mo`, from, to)
	return rows, err
}

// ReportSummary — сводные KPI за интервал.
func (r *AdminRepo) ReportSummary(from, to time.Time) (*model.ReportSummary, error) {
	var s model.ReportSummary
	err := r.db.Get(&s, `
		SELECT COUNT(*) AS total,
		       COUNT(*) FILTER (WHERE status='completed')   AS completed,
		       COUNT(*) FILTER (WHERE status='scheduled')   AS scheduled,
		       COUNT(*) FILTER (WHERE status='cancelled')   AS cancelled,
		       COUNT(*) FILTER (WHERE status='rescheduled') AS rescheduled,
		       COALESCE(SUM(final_price) FILTER (WHERE status='completed'), 0) AS revenue,
		       COUNT(DISTINCT patient_id) AS unique_patients
		FROM appointments
		WHERE starts_at >= $1 AND starts_at < $2`, from, to)
	return &s, err
}

// ReportBySpecialty — приёмы/пациенты/выручка по специальности за интервал.
func (r *AdminRepo) ReportBySpecialty(from, to time.Time) ([]model.SpecialtyReportRow, error) {
	var rows []model.SpecialtyReportRow
	err := r.db.Select(&rows, `
		SELECT sp.name AS specialty_name,
		       COUNT(a.id) FILTER (WHERE a.status NOT IN ('cancelled','rescheduled'))                  AS appointments,
		       COUNT(DISTINCT a.patient_id) FILTER (WHERE a.status NOT IN ('cancelled','rescheduled')) AS unique_patients,
		       COALESCE(SUM(a.final_price) FILTER (WHERE a.status='completed'), 0)                      AS revenue
		FROM specialties sp
		LEFT JOIN services sv     ON sv.specialty_id = sp.id
		LEFT JOIN appointments a  ON a.service_id = sv.id AND a.starts_at >= $1 AND a.starts_at < $2
		GROUP BY sp.id, sp.name
		ORDER BY revenue DESC, appointments DESC`, from, to)
	return rows, err
}

// ReportByDoctorFull — по врачу: приёмы, уникальные пациенты, выручка и средняя оценка (за всё время).
func (r *AdminRepo) ReportByDoctorFull(from, to time.Time) ([]model.DoctorFullRow, error) {
	var rows []model.DoctorFullRow
	err := r.db.Select(&rows, `
		SELECT d.full_name AS doctor_name, sp.name AS specialty_name,
		       COUNT(a.id) FILTER (WHERE a.status NOT IN ('cancelled','rescheduled'))                  AS appointments,
		       COUNT(DISTINCT a.patient_id) FILTER (WHERE a.status NOT IN ('cancelled','rescheduled')) AS unique_patients,
		       COALESCE(SUM(a.final_price) FILTER (WHERE a.status='completed'), 0)                      AS revenue,
		       COALESCE((SELECT ROUND(AVG(rv.rating)::numeric, 2)::float8 FROM reviews rv
		                 WHERE rv.doctor_id = d.id AND rv.is_hidden = false), 0)                        AS rating
		FROM doctors d
		JOIN specialties sp ON sp.id = d.specialty_id
		LEFT JOIN appointments a ON a.doctor_id = d.id AND a.starts_at >= $1 AND a.starts_at < $2
		WHERE d.is_active = true
		GROUP BY d.id, d.full_name, sp.name
		ORDER BY revenue DESC, appointments DESC`, from, to)
	return rows, err
}

// DailyStats — приёмы и выручка по каждому дню периода (включая дни без записей).
func (r *AdminRepo) DailyStats(from, to time.Time) ([]model.DailyRow, error) {
	var rows []model.DailyRow
	err := r.db.Select(&rows, `
		SELECT TO_CHAR(g::date, 'YYYY-MM-DD') AS day,
		       COALESCE(c.appointments, 0) AS appointments,
		       COALESCE(c.revenue, 0)      AS revenue
		FROM generate_series($1::date, ($2::date - 1), INTERVAL '1 day') AS g
		LEFT JOIN (
		    SELECT starts_at::date AS dt, COUNT(*) AS appointments,
		           COALESCE(SUM(final_price) FILTER (WHERE status='completed'), 0) AS revenue
		    FROM appointments WHERE starts_at >= $1 AND starts_at < $2
		    GROUP BY 1
		) c ON c.dt = g::date
		ORDER BY g`, from, to)
	return rows, err
}

// DemographicsGender — распределение пациентов по полу. Возвращает код пола (m/f/unknown);
// человекочитаемая подпись формируется на фронте через i18n.
func (r *AdminRepo) DemographicsGender() ([]model.LabelCountRow, error) {
	var rows []model.LabelCountRow
	err := r.db.Select(&rows, `
		SELECT CASE LOWER(COALESCE(NULLIF(gender, ''), '?'))
		         WHEN 'm' THEN 'm' WHEN 'f' THEN 'f' ELSE 'unknown' END AS label,
		       COUNT(*) AS count
		FROM users
		GROUP BY 1 ORDER BY count DESC`)
	return rows, err
}

// DemographicsAge — распределение пациентов по возрастным группам. Возвращает код группы
// (lt18/18_30/31_45/46_60/gt60); подпись формируется на фронте через i18n.
func (r *AdminRepo) DemographicsAge() ([]model.LabelCountRow, error) {
	var rows []model.LabelCountRow
	err := r.db.Select(&rows, `
		WITH ages AS (
		    SELECT DATE_PART('year', AGE(birth_date))::int AS age
		    FROM users WHERE birth_date IS NOT NULL
		)
		SELECT g.label, COALESCE(COUNT(a.age), 0) AS count
		FROM (VALUES (1,'lt18'),(2,'18_30'),(3,'31_45'),(4,'46_60'),(5,'gt60')) AS g(ord, label)
		LEFT JOIN ages a ON
		     (g.ord=1 AND a.age < 18) OR
		     (g.ord=2 AND a.age BETWEEN 18 AND 30) OR
		     (g.ord=3 AND a.age BETWEEN 31 AND 45) OR
		     (g.ord=4 AND a.age BETWEEN 46 AND 60) OR
		     (g.ord=5 AND a.age > 60)
		GROUP BY g.ord, g.label ORDER BY g.ord`)
	return rows, err
}
