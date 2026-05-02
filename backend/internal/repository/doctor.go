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
