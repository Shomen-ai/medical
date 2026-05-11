-- Restore 3/3 rotation for the May-June 2026 demo schedule.
-- Migration 008 originally seeded doctor_schedules with everyone working Mon-Sat,
-- which clobbered the per-specialty 3-day rotation produced by service.Generate3x3.
-- This re-applies the rotation in-place using ON CONFLICT … DO UPDATE.
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
