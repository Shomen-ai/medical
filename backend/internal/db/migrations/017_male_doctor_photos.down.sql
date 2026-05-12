-- Откат: возвращает стоковые фото для doctor 1, 2, 4, 7, 8 (как было в 004)
-- и обнуляет фото для новых doctor 9, 12, 13, 15.
UPDATE doctors SET photo_url = '/doctors/doctor-1.jpg' WHERE id = '22222222-0001-0000-0000-000000000001';
UPDATE doctors SET photo_url = '/doctors/doctor-2.jpg' WHERE id = '22222222-0001-0000-0000-000000000002';
UPDATE doctors SET photo_url = '/doctors/doctor-4.jpg' WHERE id = '22222222-0001-0000-0000-000000000004';
UPDATE doctors SET photo_url = '/doctors/doctor-7.jpg' WHERE id = '22222222-0001-0000-0000-000000000007';
UPDATE doctors SET photo_url = '/doctors/doctor-8.jpg' WHERE id = '22222222-0001-0000-0000-000000000008';
UPDATE doctors SET photo_url = NULL WHERE id = '22222222-0001-0000-0000-000000000009';
UPDATE doctors SET photo_url = NULL WHERE id = '22222222-0001-0000-0000-000000000012';
UPDATE doctors SET photo_url = NULL WHERE id = '22222222-0001-0000-0000-000000000013';
UPDATE doctors SET photo_url = NULL WHERE id = '22222222-0001-0000-0000-000000000015';
