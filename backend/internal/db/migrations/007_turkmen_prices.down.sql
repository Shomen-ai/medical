-- Rollback to the original Russian-rouble pricing from seed (002).
UPDATE services SET price =  1500 WHERE name = 'Консультация косметолога';
UPDATE services SET price =  3500 WHERE name = 'Чистка лица механическая';
UPDATE services SET price =  3000 WHERE name = 'Чистка лица ультразвуковая';
UPDATE services SET price =  4500 WHERE name = 'Химический пилинг';
UPDATE services SET price =  8000 WHERE name = 'Инъекции ботулотоксина (зона)';
UPDATE services SET price = 12000 WHERE name = 'Контурная пластика (1 мл)';
UPDATE services SET price =  9000 WHERE name = 'Биоревитализация';
UPDATE services SET price =  7500 WHERE name = 'Мезотерапия лица';
UPDATE services SET price = 10000 WHERE name = 'PRP-терапия лица';
UPDATE services SET price =  4000 WHERE name = 'Микротоковая терапия';

UPDATE services SET price =  1800 WHERE name = 'Консультация дерматолога';
UPDATE services SET price =  5000 WHERE name = 'Лечение акне';
UPDATE services SET price =   800 WHERE name = 'Дерматоскопия (1 элемент)';
UPDATE services SET price =  1200 WHERE name = 'Удаление новообразований (до 5 мм)';
UPDATE services SET price =  3500 WHERE name = 'Лечение розацеа';

UPDATE services SET price =  2000 WHERE name = 'Консультация трихолога';
UPDATE services SET price =  6500 WHERE name = 'Мезотерапия волосистой части головы';
UPDATE services SET price =  9500 WHERE name = 'PRP-терапия волос';
UPDATE services SET price =  3500 WHERE name = 'Озонотерапия волос';

UPDATE services SET price =  2500 WHERE name = 'Консультация эстетиста';
UPDATE services SET price = 25000 WHERE name = 'Нитевой лифтинг (нос-губы)';
UPDATE services SET price = 35000 WHERE name = 'Объёмное моделирование лица';
