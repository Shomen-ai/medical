-- Чистка «мусорных» пользователей, оставшихся от ручного тестирования.
-- Это записи без full_name: 4 с реальными RU-номерами (тестировались до
-- появления валидации +993), 3 с явным мусором ('вапв', '777777777',
-- '123456', '+78888888888'). У них нет медкарт; у двоих есть по одному
-- приёму — удаляем их вместе с пациентами.
--
-- Порядок DELETE учитывает FK: сперва зависимые таблицы, потом users.
-- WHERE-условие одинаковое во всех трёх запросах, чтобы атомарно
-- захватить ровно тех же 7 пользователей.

WITH junk_users AS (
    SELECT id FROM users WHERE full_name IS NULL OR full_name = ''
)
DELETE FROM appointment_records
 WHERE appointment_id IN (
     SELECT id FROM appointments WHERE patient_id IN (SELECT id FROM junk_users)
 );

WITH junk_users AS (
    SELECT id FROM users WHERE full_name IS NULL OR full_name = ''
)
DELETE FROM appointments
 WHERE patient_id IN (SELECT id FROM junk_users);

DELETE FROM users
 WHERE full_name IS NULL OR full_name = '';
