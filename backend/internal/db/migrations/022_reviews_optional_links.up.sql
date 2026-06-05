-- Файл: 022_reviews_optional_links.up.sql
-- Назначение: отзыв может оставить любой авторизованный пациент — без привязки к завершённому
-- визиту. Поэтому appointment_id/doctor_id/service_id становятся необязательными (пациент
-- указывает врача/услугу по желанию).
ALTER TABLE reviews ALTER COLUMN appointment_id DROP NOT NULL;
ALTER TABLE reviews ALTER COLUMN doctor_id      DROP NOT NULL;
ALTER TABLE reviews ALTER COLUMN service_id     DROP NOT NULL;
