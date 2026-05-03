package repository

import (
	"beautymed/internal/model"

	"github.com/jmoiron/sqlx"
)

type DoctorRepo struct{ db *sqlx.DB }

func NewDoctorRepo(db *sqlx.DB) *DoctorRepo { return &DoctorRepo{db} }

func (r *DoctorRepo) List() ([]model.Doctor, error) {
	var ds []model.Doctor
	err := r.db.Select(&ds, `
		SELECT d.*, s.name AS specialty_name
		FROM doctors d
		JOIN specialties s ON s.id = d.specialty_id
		WHERE d.is_active = true
		ORDER BY d.full_name`)
	return ds, err
}

func (r *DoctorRepo) ListBySpecialty(specialtyID string) ([]model.Doctor, error) {
	var ds []model.Doctor
	err := r.db.Select(&ds, `
		SELECT d.*, s.name AS specialty_name
		FROM doctors d
		JOIN specialties s ON s.id = d.specialty_id
		WHERE d.specialty_id = $1 AND d.is_active = true
		ORDER BY d.full_name`, specialtyID)
	return ds, err
}

func (r *DoctorRepo) FindByID(id string) (*model.Doctor, error) {
	var d model.Doctor
	err := r.db.Get(&d, `
		SELECT d.*, s.name AS specialty_name
		FROM doctors d
		JOIN specialties s ON s.id = d.specialty_id
		WHERE d.id = $1`, id)
	return &d, err
}

func (r *DoctorRepo) FindStaffByPhone(phone string) (*model.Staff, error) {
	var st model.Staff
	err := r.db.Get(&st, `SELECT * FROM staff WHERE phone = $1 AND is_active = true`, phone)
	return &st, err
}

// CreateDoctor inserts a new doctor and returns the generated UUID.
func (r *DoctorRepo) CreateDoctor(d *model.Doctor) error {
	return r.db.QueryRow(`
		INSERT INTO doctors (full_name, specialty_id, phone, bio, photo_url, experience_years, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at`,
		d.FullName, d.SpecialtyID, d.Phone, d.Bio, d.PhotoURL, d.ExperienceYears, d.IsActive,
	).Scan(&d.ID, &d.CreatedAt)
}

// UpdateDoctor updates mutable doctor fields.
func (r *DoctorRepo) UpdateDoctor(d *model.Doctor) error {
	_, err := r.db.Exec(`
		UPDATE doctors
		SET full_name=COALESCE(NULLIF($1,''), full_name),
		    bio=COALESCE(NULLIF($2,''), bio),
		    photo_url=COALESCE(NULLIF($3,''), photo_url),
		    experience_years=CASE WHEN $4=0 THEN experience_years ELSE $4 END,
		    is_active=$5
		WHERE id=$6`,
		d.FullName, d.Bio, d.PhotoURL, d.ExperienceYears, d.IsActive, d.ID)
	return err
}

// CreateStaff creates a staff login record for a doctor.
func (r *DoctorRepo) CreateStaff(doctorID, phone, role string) error {
	_, err := r.db.Exec(`
		INSERT INTO staff (doctor_id, phone, role)
		VALUES ($1, $2, $3)
		ON CONFLICT (phone) DO NOTHING`,
		doctorID, phone, role)
	return err
}
