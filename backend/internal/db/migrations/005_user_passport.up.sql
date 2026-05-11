-- Add fields required by Russian tax certificate form КНД 1151156.
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS inn               varchar(12),
    ADD COLUMN IF NOT EXISTS passport_series   varchar(10),
    ADD COLUMN IF NOT EXISTS passport_number   varchar(20),
    ADD COLUMN IF NOT EXISTS passport_issued_at date,
    ADD COLUMN IF NOT EXISTS passport_issued_by text;
