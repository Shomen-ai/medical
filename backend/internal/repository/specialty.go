package repository

import (
	"beautymed/internal/model"

	"github.com/jmoiron/sqlx"
)

type SpecialtyRepo struct{ db *sqlx.DB }

func NewSpecialtyRepo(db *sqlx.DB) *SpecialtyRepo { return &SpecialtyRepo{db} }

func (r *SpecialtyRepo) List() ([]model.Specialty, error) {
	var ss []model.Specialty
	err := r.db.Select(&ss, `SELECT * FROM specialties ORDER BY name`)
	return ss, err
}

func (r *SpecialtyRepo) FindByID(id string) (*model.Specialty, error) {
	var s model.Specialty
	err := r.db.Get(&s, `SELECT * FROM specialties WHERE id = $1`, id)
	return &s, err
}
