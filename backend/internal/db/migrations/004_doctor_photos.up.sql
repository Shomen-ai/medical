-- Bind each seeded doctor to a specific photo slot in /public/doctors/.
-- Photo files must match gender of each doctor:
--   doctor-1.jpg .. doctor-5.jpg — female (5 photos)
--   doctor-6.jpg               — male
--   doctor-7.jpg               — female
--   doctor-8.jpg               — male
UPDATE doctors SET photo_url = '/doctors/doctor-1.jpg' WHERE id = '22222222-0001-0000-0000-000000000001';
UPDATE doctors SET photo_url = '/doctors/doctor-2.jpg' WHERE id = '22222222-0001-0000-0000-000000000002';
UPDATE doctors SET photo_url = '/doctors/doctor-3.jpg' WHERE id = '22222222-0001-0000-0000-000000000003';
UPDATE doctors SET photo_url = '/doctors/doctor-4.jpg' WHERE id = '22222222-0001-0000-0000-000000000004';
UPDATE doctors SET photo_url = '/doctors/doctor-5.jpg' WHERE id = '22222222-0001-0000-0000-000000000005';
UPDATE doctors SET photo_url = '/doctors/doctor-6.jpg' WHERE id = '22222222-0001-0000-0000-000000000006';
UPDATE doctors SET photo_url = '/doctors/doctor-7.jpg' WHERE id = '22222222-0001-0000-0000-000000000007';
UPDATE doctors SET photo_url = '/doctors/doctor-8.jpg' WHERE id = '22222222-0001-0000-0000-000000000008';
