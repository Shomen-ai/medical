// Файл: internal/repository/service.go
// Назначение: SQL-доступ к таблице услуг (services) — листинг и поиск по ID.
package repository

import (
	"beautymed/internal/model"

	"github.com/jmoiron/sqlx"
)

// ServiceRepo — репозиторий медицинских услуг клиники.
type ServiceRepo struct{ db *sqlx.DB }

// NewServiceRepo создаёт ServiceRepo на базе пула sqlx.
func NewServiceRepo(db *sqlx.DB) *ServiceRepo { return &ServiceRepo{db} }

// List возвращает все активные услуги, опционально отфильтрованные по специальности.
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

// FindByID возвращает активную услугу по ID.
func (r *ServiceRepo) FindByID(id string) (*model.Service, error) {
	var s model.Service
	err := r.db.Get(&s, `SELECT * FROM services WHERE id=$1 AND is_active=true`, id)
	return &s, err
}
