-- Файл: 023_id_doc_type.up.sql
-- Назначение: тип удостоверения личности. Внутренний паспорт — как раньше (номер + кем выдан);
-- загранпаспорт добавляет дату выдачи и срок действия.
ALTER TABLE users ADD COLUMN IF NOT EXISTS id_doc_type        varchar(20) NOT NULL DEFAULT 'domestic';
ALTER TABLE users ADD COLUMN IF NOT EXISTS id_doc_issued_at   date;
ALTER TABLE users ADD COLUMN IF NOT EXISTS id_doc_valid_until date;
