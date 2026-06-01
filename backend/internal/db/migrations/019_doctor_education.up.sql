-- Добавляет поле «образование» врача (ВУЗ, год выпуска, специализация)
-- для отображения в публичной карточке врача на главной странице.
ALTER TABLE doctors
    ADD COLUMN IF NOT EXISTS education varchar(300) NOT NULL DEFAULT '';
