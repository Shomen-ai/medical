package repository

import (
	"beautymed/internal/model"

	"github.com/jmoiron/sqlx"
)

type ServiceRepo struct{ db *sqlx.DB }

func NewServiceRepo(db *sqlx.DB) *ServiceRepo { return &ServiceRepo{db} }

func (r *ServiceRepo) List(specialtyID string) ([]model.Service, error) {
	var ss []model.Service
	if specialtyID != "" {
		err := r.db.Select(&ss,
			`SELECT * FROM services WHERE specialty_id=$1 AND is_active=true ORDER BY name`,
			specialtyID)
		return ss, err
	}
	err := r.db.Select(&ss, `SELECT * FROM services WHERE is_active=true ORDER BY name`)
	return ss, err
}

func (r *ServiceRepo) FindByID(id string) (*model.Service, error) {
	var s model.Service
	err := r.db.Get(&s, `SELECT * FROM services WHERE id=$1 AND is_active=true`, id)
	return &s, err
}
