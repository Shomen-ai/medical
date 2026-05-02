package model

import "time"

type Doctor struct {
	ID              string    `db:"id"               json:"id"`
	FullName        string    `db:"full_name"        json:"full_name"`
	SpecialtyID     string    `db:"specialty_id"     json:"specialty_id"`
	SpecialtyName   string    `db:"specialty_name"   json:"specialty_name,omitempty"`
	Phone           string    `db:"phone"            json:"-"`
	Bio             string    `db:"bio"              json:"bio"`
	PhotoURL        string    `db:"photo_url"        json:"photo_url"`
	ExperienceYears int       `db:"experience_years" json:"experience_years"`
	IsActive        bool      `db:"is_active"        json:"is_active"`
	CreatedAt       time.Time `db:"created_at"       json:"created_at"`
}
