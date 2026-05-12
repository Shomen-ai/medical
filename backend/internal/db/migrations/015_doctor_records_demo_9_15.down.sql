-- Откат: удаляет demo-приёмы и медкарты у новых врачей 9..15.
DELETE FROM appointment_records
 WHERE appointment_id BETWEEN '55555555-0000-0000-0000-000000000901'
                         AND  '55555555-0000-0000-0000-000000001504';

DELETE FROM appointments
 WHERE id BETWEEN '55555555-0000-0000-0000-000000000901'
             AND  '55555555-0000-0000-0000-000000001504';
