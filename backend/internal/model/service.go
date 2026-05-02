package model

type Service struct {
	ID          string  `db:"id"           json:"id"`
	Name        string  `db:"name"         json:"name"`
	Description string  `db:"description"  json:"description"`
	Price       float64 `db:"price"        json:"price"`
	DurationMin int     `db:"duration_min" json:"duration_min"`
	SpecialtyID string  `db:"specialty_id" json:"specialty_id"`
	IsActive    bool    `db:"is_active"    json:"is_active"`
}
