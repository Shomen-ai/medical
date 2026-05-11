-- Strip the demo content but leave seeded specialties / doctors / services intact.
DELETE FROM appointment_records WHERE appointment_id::text LIKE '44444444-0000-%';
DELETE FROM appointments         WHERE id::text           LIKE '44444444-0000-%';
DELETE FROM promo_codes WHERE code IN ('SUMMER10','WELCOME15','VIP20','NEWPATIENT5');
DELETE FROM doctor_schedules
 WHERE work_date BETWEEN DATE '2026-05-01' AND DATE '2026-06-30';
DELETE FROM users WHERE id::text LIKE '33333333-0000-%';
