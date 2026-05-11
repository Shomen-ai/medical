DELETE FROM appointment_records WHERE appointment_id::text LIKE '55555555-0000-%';
DELETE FROM appointments         WHERE id::text           LIKE '55555555-0000-%';
