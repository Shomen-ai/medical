-- Файл: 021_reviews.up.sql
-- Назначение: таблица анонимных отзывов пациентов (привязаны к завершённому визиту) + демо-сид.
CREATE TABLE reviews (
  id             uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id        uuid NOT NULL REFERENCES users(id),
  appointment_id uuid NOT NULL REFERENCES appointments(id),
  doctor_id      uuid NOT NULL REFERENCES doctors(id),
  service_id     uuid NOT NULL REFERENCES services(id),
  rating         smallint NOT NULL CHECK (rating BETWEEN 1 AND 5),
  text           text NOT NULL CHECK (length(trim(text)) > 0),
  is_hidden      boolean NOT NULL DEFAULT false,
  created_at     timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX idx_reviews_public  ON reviews (is_hidden, created_at DESC);
CREATE INDEX idx_reviews_doctor  ON reviews (doctor_id);
CREATE INDEX idx_reviews_service ON reviews (service_id);

-- Демо-отзывы из существующих завершённых визитов (разные врачи/услуги для наглядной фильтрации).
INSERT INTO reviews (user_id, appointment_id, doctor_id, service_id, rating, text)
SELECT a.patient_id, a.id, a.doctor_id, a.service_id,
       (4 + (row_number() OVER (ORDER BY a.starts_at DESC)) % 2)::smallint,
       CASE (row_number() OVER (ORDER BY a.starts_at DESC)) % 3
         WHEN 0 THEN 'Очень довольна результатом, врач внимательный. Спасибо клинике!'
         WHEN 1 THEN 'Записалась онлайн за пару минут, приём прошёл отлично. Рекомендую.'
         ELSE 'Профессиональный подход и приятная атмосфера. Обязательно вернусь.'
       END
FROM appointments a
WHERE a.status = 'completed'
ORDER BY a.starts_at DESC
LIMIT 6;
