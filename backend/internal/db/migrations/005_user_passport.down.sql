ALTER TABLE users
    DROP COLUMN IF EXISTS passport_issued_by,
    DROP COLUMN IF EXISTS passport_issued_at,
    DROP COLUMN IF EXISTS passport_number,
    DROP COLUMN IF EXISTS passport_series,
    DROP COLUMN IF EXISTS inn;
