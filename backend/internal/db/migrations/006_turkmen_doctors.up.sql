-- Adapt seeded doctor identities to Turkmenistan clinic context.
-- Names and Turkmen phone codes; gender matches the assigned photo
-- (see migration 004 for photo→ID mapping):
--   doctor-1, 2, 4, 7, 8 → male
--   doctor-3, 5, 6       → female
UPDATE doctors SET
    full_name = 'Аман Гулиев',
    phone     = '+99365100001',
    bio       = 'Врач-косметолог. Аппаратная косметология, инъекционные методики, комплексное омоложение.'
WHERE id = '22222222-0001-0000-0000-000000000001';

UPDATE doctors SET
    full_name = 'Сапарбай Ходжаев',
    phone     = '+99365100002',
    bio       = 'Сертифицированный косметолог. Контурная пластика, биоревитализация, ботулотерапия.'
WHERE id = '22222222-0001-0000-0000-000000000002';

UPDATE doctors SET
    full_name = 'Айна Атаева',
    phone     = '+99365100003',
    bio       = 'Косметолог-эстетист. Уходовые процедуры, пилинги, чистки лица.'
WHERE id = '22222222-0001-0000-0000-000000000003';

UPDATE doctors SET
    full_name = 'Меретгелди Назаров',
    phone     = '+99365100004',
    bio       = 'Врач-дерматолог. Лечение акне, розацеа, дерматитов. Удаление новообразований.'
WHERE id = '22222222-0001-0000-0000-000000000004';

UPDATE doctors SET
    full_name = 'Огулбике Гельдыева',
    phone     = '+99365100005',
    bio       = 'Дерматолог. Диагностика и лечение кожных заболеваний, консультации по уходу.'
WHERE id = '22222222-0001-0000-0000-000000000005';

UPDATE doctors SET
    full_name = 'Дженнет Овезова',
    phone     = '+99365100006',
    bio       = 'Врач-трихолог. Мезотерапия волосистой части головы, диагностика выпадения волос.'
WHERE id = '22222222-0001-0000-0000-000000000006';

UPDATE doctors SET
    full_name = 'Якуб Бердыев',
    phone     = '+99365100007',
    bio       = 'Трихолог, дерматолог. PRP-терапия, лечение заболеваний кожи головы.'
WHERE id = '22222222-0001-0000-0000-000000000007';

UPDATE doctors SET
    full_name = 'Атамурад Курбанов',
    phone     = '+99365100008',
    bio       = 'Врач эстетической медицины. Нитевой лифтинг, объёмное моделирование лица.'
WHERE id = '22222222-0001-0000-0000-000000000008';

-- Keep staff.phone in sync with doctors.phone (referenced for login).
UPDATE staff SET phone = '+99365100001' WHERE doctor_id = '22222222-0001-0000-0000-000000000001';
UPDATE staff SET phone = '+99365100002' WHERE doctor_id = '22222222-0001-0000-0000-000000000002';
UPDATE staff SET phone = '+99365100003' WHERE doctor_id = '22222222-0001-0000-0000-000000000003';
UPDATE staff SET phone = '+99365100004' WHERE doctor_id = '22222222-0001-0000-0000-000000000004';
UPDATE staff SET phone = '+99365100005' WHERE doctor_id = '22222222-0001-0000-0000-000000000005';
UPDATE staff SET phone = '+99365100006' WHERE doctor_id = '22222222-0001-0000-0000-000000000006';
UPDATE staff SET phone = '+99365100007' WHERE doctor_id = '22222222-0001-0000-0000-000000000007';
UPDATE staff SET phone = '+99365100008' WHERE doctor_id = '22222222-0001-0000-0000-000000000008';
