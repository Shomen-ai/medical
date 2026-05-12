-- Откат: возвращает оригинальные стоковые фото для doctor 3, 5, 6 (поставленные в 004)
-- и обнуляет фото для doctor 10, 11, 14 (они были NULL до 016).
UPDATE doctors SET photo_url = '/doctors/doctor-3.jpg' WHERE id = '22222222-0001-0000-0000-000000000003';
UPDATE doctors SET photo_url = '/doctors/doctor-5.jpg' WHERE id = '22222222-0001-0000-0000-000000000005';
UPDATE doctors SET photo_url = '/doctors/doctor-6.jpg' WHERE id = '22222222-0001-0000-0000-000000000006';
UPDATE doctors SET photo_url = NULL WHERE id = '22222222-0001-0000-0000-000000000010';
UPDATE doctors SET photo_url = NULL WHERE id = '22222222-0001-0000-0000-000000000011';
UPDATE doctors SET photo_url = NULL WHERE id = '22222222-0001-0000-0000-000000000014';
