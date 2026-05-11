-- Re-price seeded services to Turkmen manat (TMT) market rates.
-- Calibrated against the 2026 average monthly salary of ~1,410 TMT
-- and Russian/CIS clinic price benchmarks scaled by purchasing-power parity.
-- Aesthetic-medicine premium services keep a higher relative margin.

-- ── Косметология ─────────────────────────────────────────────────────
UPDATE services SET price =   50 WHERE name = 'Консультация косметолога';
UPDATE services SET price =  120 WHERE name = 'Чистка лица механическая';
UPDATE services SET price =  100 WHERE name = 'Чистка лица ультразвуковая';
UPDATE services SET price =  180 WHERE name = 'Химический пилинг';
UPDATE services SET price =  350 WHERE name = 'Инъекции ботулотоксина (зона)';
UPDATE services SET price =  600 WHERE name = 'Контурная пластика (1 мл)';
UPDATE services SET price =  450 WHERE name = 'Биоревитализация';
UPDATE services SET price =  350 WHERE name = 'Мезотерапия лица';
UPDATE services SET price =  500 WHERE name = 'PRP-терапия лица';
UPDATE services SET price =  180 WHERE name = 'Микротоковая терапия';

-- ── Дерматология ─────────────────────────────────────────────────────
UPDATE services SET price =   60 WHERE name = 'Консультация дерматолога';
UPDATE services SET price =  200 WHERE name = 'Лечение акне';
UPDATE services SET price =   40 WHERE name = 'Дерматоскопия (1 элемент)';
UPDATE services SET price =   60 WHERE name = 'Удаление новообразований (до 5 мм)';
UPDATE services SET price =  180 WHERE name = 'Лечение розацеа';

-- ── Трихология ───────────────────────────────────────────────────────
UPDATE services SET price =   80 WHERE name = 'Консультация трихолога';
UPDATE services SET price =  300 WHERE name = 'Мезотерапия волосистой части головы';
UPDATE services SET price =  450 WHERE name = 'PRP-терапия волос';
UPDATE services SET price =  180 WHERE name = 'Озонотерапия волос';

-- ── Эстетическая медицина ────────────────────────────────────────────
UPDATE services SET price =  100 WHERE name = 'Консультация эстетиста';
UPDATE services SET price = 1200 WHERE name = 'Нитевой лифтинг (нос-губы)';
UPDATE services SET price = 1800 WHERE name = 'Объёмное моделирование лица';
