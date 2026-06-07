-- Файл: 023_id_doc_type.down.sql
ALTER TABLE users DROP COLUMN IF EXISTS id_doc_valid_until;
ALTER TABLE users DROP COLUMN IF EXISTS id_doc_issued_at;
ALTER TABLE users DROP COLUMN IF EXISTS id_doc_type;
