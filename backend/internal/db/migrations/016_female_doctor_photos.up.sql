-- Привязка новых фото врачей-женщин к их записям. Файлы лежат в
-- frontend/public/doctors/doctor_women_{1..6}.jpg и раздаются Nuxt'ом
-- по пути /doctors/doctor_women_{1..6}.jpg.
--
-- Сопоставление — по возрастанию UUID врача (детерминированно). Если нужно
-- другое соответствие лицо↔врач, переназначайте через POST /api/admin/doctors/:id.

UPDATE doctors SET photo_url = '/doctors/doctor_women_1.jpg' WHERE id = '22222222-0001-0000-0000-000000000003'; -- Айна Атаева (косметолог)
UPDATE doctors SET photo_url = '/doctors/doctor_women_2.jpg' WHERE id = '22222222-0001-0000-0000-000000000005'; -- Огулбике Гельдыева (дерматолог)
UPDATE doctors SET photo_url = '/doctors/doctor_women_3.jpg' WHERE id = '22222222-0001-0000-0000-000000000006'; -- Дженнет Овезова (трихолог)
UPDATE doctors SET photo_url = '/doctors/doctor_women_4.jpg' WHERE id = '22222222-0001-0000-0000-000000000010'; -- Сельби Атаева (косметолог)
UPDATE doctors SET photo_url = '/doctors/doctor_women_5.jpg' WHERE id = '22222222-0001-0000-0000-000000000011'; -- Майса Бердыева (дерматолог)
UPDATE doctors SET photo_url = '/doctors/doctor_women_6.jpg' WHERE id = '22222222-0001-0000-0000-000000000014'; -- Огулджемал Гулиева (эстетист)
