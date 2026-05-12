-- Откат миграции 011: восстанавливает старые поля под российский налоговый вычет
-- (пустые), удаляет туркменские поля. Данные, записанные в новые столбцы, потеряются.

ALTER TABLE users
    DROP CONSTRAINT IF EXISTS users_gender_check;

ALTER TABLE users
    DROP COLUMN IF EXISTS gender,
    DROP COLUMN IF EXISTS address,
    DROP COLUMN IF EXISTS id_doc_number,
    DROP COLUMN IF EXISTS id_doc_issued_by;

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS inn                varchar(12),
    ADD COLUMN IF NOT EXISTS passport_series    varchar(10),
    ADD COLUMN IF NOT EXISTS passport_number    varchar(20),
    ADD COLUMN IF NOT EXISTS passport_issued_at date,
    ADD COLUMN IF NOT EXISTS passport_issued_by text;
