-- Пересчёт ротации 3/3 для май-июнь 2026 с учётом всех 15 врачей.
-- После миграции 013 в БД появилось 7 новых врачей, но doctor_schedules остался от 8.
-- Чтобы новые врачи появились в календарях и ротация перераспределилась равномерно,
-- переписываем расписание тем же алгоритмом, что и миграция 009 (SQL-зеркало Generate3x3).
--
-- Что НЕ трогается:
--   • appointments — записи на приём остаются на месте. Если новая ротация делает
--     врача «выходным» в день его существующего приёма, прямой конфликт БД отсутствует
--     (FK appointments.doctor_id → doctors, без зависимости от doctor_schedules);
--     визуально админ это увидит и при необходимости поправит вручную.
--   • appointment_records — не трогаются.
--
-- Окно: 2026-05-01 .. 2026-06-30. Будущие месяцы будут сгенерированы автоматически
-- ежемесячным cron-jobом (см. service/scheduler.go::generateNextMonthSchedule),
-- либо вручную через POST /api/admin/schedule/generate.

WITH dates AS (
  SELECT dt::date AS d,
         (EXTRACT(DOW FROM dt) = 0) OR (dt::date = DATE '2026-05-18') AS is_universal_off
  FROM generate_series(DATE '2026-05-01', DATE '2026-06-30', INTERVAL '1 day') AS dt
),
workdays AS (
  SELECT d, ROW_NUMBER() OVER (ORDER BY d) AS workday_n
  FROM dates
  WHERE NOT is_universal_off
),
doctor_idx AS (
  SELECT id, specialty_id,
         ROW_NUMBER() OVER (PARTITION BY specialty_id ORDER BY id) - 1 AS idx,
         COUNT(*)     OVER (PARTITION BY specialty_id)               AS spec_count
  FROM doctors
)
INSERT INTO doctor_schedules (doctor_id, work_date, start_time, end_time, is_day_off)
SELECT di.id,
       dates.d,
       '09:00'::time,
       '18:00'::time,
       CASE
         WHEN dates.is_universal_off       THEN TRUE
         WHEN di.spec_count <= 1           THEN FALSE
         ELSE (((w.workday_n - 1) / 3) % di.spec_count) <> di.idx
       END
FROM dates
CROSS JOIN doctor_idx di
LEFT JOIN workdays w ON w.d = dates.d
ON CONFLICT (doctor_id, work_date) DO UPDATE SET
  start_time = EXCLUDED.start_time,
  end_time   = EXCLUDED.end_time,
  is_day_off = EXCLUDED.is_day_off;
