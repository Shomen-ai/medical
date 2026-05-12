-- Привязка новых фото врачей-мужчин (9 шт.) к их записям. Файлы — в
-- frontend/public/doctors/doctor_men_{1..9}.jpg, раздаются Nuxt'ом
-- по пути /doctors/doctor_men_{1..9}.jpg.
--
-- Соответствие — по возрастанию UUID врача (детерминированно). Поменять
-- порядок можно через POST/PATCH /api/admin/doctors/:id.

UPDATE doctors SET photo_url = '/doctors/doctor_men_1.jpg' WHERE id = '22222222-0001-0000-0000-000000000001'; -- Аман Гулиев         (косметолог)
UPDATE doctors SET photo_url = '/doctors/doctor_men_2.jpg' WHERE id = '22222222-0001-0000-0000-000000000002'; -- Сапарбай Ходжаев    (косметолог)
UPDATE doctors SET photo_url = '/doctors/doctor_men_3.jpg' WHERE id = '22222222-0001-0000-0000-000000000004'; -- Меретгелди Назаров  (дерматолог)
UPDATE doctors SET photo_url = '/doctors/doctor_men_4.jpg' WHERE id = '22222222-0001-0000-0000-000000000007'; -- Якуб Бердыев        (трихолог)
UPDATE doctors SET photo_url = '/doctors/doctor_men_5.jpg' WHERE id = '22222222-0001-0000-0000-000000000008'; -- Атамурад Курбанов   (эстетист)
UPDATE doctors SET photo_url = '/doctors/doctor_men_6.jpg' WHERE id = '22222222-0001-0000-0000-000000000009'; -- Мухамметгулы Овезов (косметолог)
UPDATE doctors SET photo_url = '/doctors/doctor_men_7.jpg' WHERE id = '22222222-0001-0000-0000-000000000012'; -- Кемал Розыев        (дерматолог)
UPDATE doctors SET photo_url = '/doctors/doctor_men_8.jpg' WHERE id = '22222222-0001-0000-0000-000000000013'; -- Мердан Аннаев       (трихолог)
UPDATE doctors SET photo_url = '/doctors/doctor_men_9.jpg' WHERE id = '22222222-0001-0000-0000-000000000015'; -- Бегенч Реджепов     (эстетист)
