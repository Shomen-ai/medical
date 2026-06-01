-- Откат сида: очищаем образование у засеянных врачей.
UPDATE doctors SET education = '' WHERE id LIKE '22222222-0001-0000-0000-%';
