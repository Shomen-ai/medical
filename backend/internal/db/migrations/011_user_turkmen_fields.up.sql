-- Замена российских налоговых полей на минимальный набор данных,
-- который туркменская клиника фиксирует при первом приёме пациента:
--   • gender            — пол (m / f) для медконтекста;
--   • address           — адрес проживания (welaýat, etrap, köçe, jaý);
--   • id_doc_number     — номер удостоверения личности (паспорт Туркменистана);
--   • id_doc_issued_by  — кем выдан документ.
-- Старые поля (inn, passport_*) удаляются — они вводились под российский налоговый вычет
-- и в туркменских клиниках не используются.

ALTER TABLE users
    DROP COLUMN IF EXISTS inn,
    DROP COLUMN IF EXISTS passport_series,
    DROP COLUMN IF EXISTS passport_number,
    DROP COLUMN IF EXISTS passport_issued_at,
    DROP COLUMN IF EXISTS passport_issued_by;

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS gender            char(1),
    ADD COLUMN IF NOT EXISTS address           text,
    ADD COLUMN IF NOT EXISTS id_doc_number     varchar(30),
    ADD COLUMN IF NOT EXISTS id_doc_issued_by  text;

-- Допустимые значения для пола: 'm' (мужской), 'f' (женский). NULL = не указан.
ALTER TABLE users
    ADD CONSTRAINT users_gender_check CHECK (gender IS NULL OR gender IN ('m', 'f'));
