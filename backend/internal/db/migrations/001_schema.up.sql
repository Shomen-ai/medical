CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE specialties (
    id              uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name            varchar(100) NOT NULL,
    slot_duration_min int NOT NULL DEFAULT 30
);

CREATE TABLE doctors (
    id               uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    full_name        varchar(200) NOT NULL,
    specialty_id     uuid NOT NULL REFERENCES specialties(id),
    phone            varchar(20) UNIQUE NOT NULL,
    bio              text NOT NULL DEFAULT '',
    photo_url        varchar(500) NOT NULL DEFAULT '',
    experience_years int NOT NULL DEFAULT 0,
    is_active        boolean NOT NULL DEFAULT true,
    created_at       timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE staff (
    id        uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    doctor_id uuid REFERENCES doctors(id),
    phone     varchar(20) UNIQUE NOT NULL,
    role      varchar(20) NOT NULL CHECK (role IN ('doctor','admin')),
    is_active boolean NOT NULL DEFAULT true
);

CREATE TABLE users (
    id         uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    phone      varchar(20) UNIQUE NOT NULL,
    full_name  varchar(200) NOT NULL DEFAULT '',
    birth_date date,
    email      varchar(200),
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE services (
    id           uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name         varchar(200) NOT NULL,
    description  text NOT NULL DEFAULT '',
    price        numeric(10,2) NOT NULL,
    duration_min int NOT NULL DEFAULT 30,
    specialty_id uuid NOT NULL REFERENCES specialties(id),
    is_active    boolean NOT NULL DEFAULT true
);

CREATE TABLE promo_codes (
    id          uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    code        varchar(50) UNIQUE NOT NULL,
    discount_pct int NOT NULL CHECK (discount_pct BETWEEN 1 AND 100),
    max_uses    int,
    used_count  int NOT NULL DEFAULT 0,
    valid_from  date NOT NULL DEFAULT CURRENT_DATE,
    valid_until date,
    is_active   boolean NOT NULL DEFAULT true,
    created_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TYPE appointment_status AS ENUM ('scheduled','completed','cancelled','rescheduled');
CREATE TYPE appointment_creator AS ENUM ('patient','admin');

CREATE TABLE appointments (
    id            uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    patient_id    uuid NOT NULL REFERENCES users(id),
    doctor_id     uuid NOT NULL REFERENCES doctors(id),
    service_id    uuid NOT NULL REFERENCES services(id),
    promo_code_id uuid REFERENCES promo_codes(id),
    starts_at     timestamptz NOT NULL,
    ends_at       timestamptz NOT NULL,
    status        appointment_status NOT NULL DEFAULT 'scheduled',
    final_price   numeric(10,2) NOT NULL,
    created_by    appointment_creator NOT NULL DEFAULT 'patient',
    created_at    timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_appointments_patient ON appointments(patient_id);
CREATE INDEX idx_appointments_doctor  ON appointments(doctor_id);
CREATE INDEX idx_appointments_starts  ON appointments(starts_at);

CREATE TABLE appointment_records (
    id             uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    appointment_id uuid UNIQUE NOT NULL REFERENCES appointments(id),
    complaints     text NOT NULL DEFAULT '',
    diagnosis      text NOT NULL DEFAULT '',
    prescription   text,
    recommendations text,
    is_draft       boolean NOT NULL DEFAULT true,
    created_at     timestamptz NOT NULL DEFAULT now(),
    updated_at     timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE doctor_schedules (
    id         uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    doctor_id  uuid NOT NULL REFERENCES doctors(id),
    work_date  date NOT NULL,
    start_time time NOT NULL DEFAULT '09:00',
    end_time   time NOT NULL DEFAULT '18:00',
    is_day_off boolean NOT NULL DEFAULT false,
    UNIQUE(doctor_id, work_date)
);

CREATE INDEX idx_schedules_doctor_date ON doctor_schedules(doctor_id, work_date);

CREATE TABLE tax_receipts (
    id           uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    patient_id   uuid NOT NULL REFERENCES users(id),
    year         int NOT NULL,
    generated_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE(patient_id, year)
);
