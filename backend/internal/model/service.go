// Файл: internal/model/service.go
// Назначение: модель медицинской услуги — цена, длительность приёма, связь со специальностью.
package model

// Service — медицинская услуга клиники с ценой и длительностью в минутах.
type Service struct {
	ID          string  `db:"id"           json:"id"`
	Name        string  `db:"name"         json:"name"`
	Description string  `db:"description"  json:"description"`
	Price       float64 `db:"price"        json:"price"`
	DurationMin int     `db:"duration_min" json:"duration_min"`
	SpecialtyID string  `db:"specialty_id" json:"specialty_id"`
	IsActive    bool    `db:"is_active"    json:"is_active"`
}
