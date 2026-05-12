// Файл: internal/repository/specialty.go
// Назначение: SQL-доступ к таблице специальностей (specialties).
package repository

import (
	"beautymed/internal/model"

	"github.com/jmoiron/sqlx"
)

// SpecialtyRepo — репозиторий медицинских специальностей.
type SpecialtyRepo struct{ db *sqlx.DB }

// NewSpecialtyRepo создаёт SpecialtyRepo на базе пула sqlx.
func NewSpecialtyRepo(db *sqlx.DB) *SpecialtyRepo { return &SpecialtyRepo{db} }

// List возвращает все специальности, отсортированные по названию.
func (r *SpecialtyRepo) List() ([]model.Specialty, error) {
	var ss []model.Specialty
	err := r.db.Select(&ss, `SELECT * FROM specialties ORDER BY name`)
	return ss, err
}

// FindByID возвращает специальность по ID.
func (r *SpecialtyRepo) FindByID(id string) (*model.Specialty, error) {
	var s model.Specialty
	err := r.db.Get(&s, `SELECT * FROM specialties WHERE id = $1`, id)
	return &s, err
}
