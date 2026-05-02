package model

import "time"

type User struct {
	ID        string     `db:"id"         json:"id"`
	Phone     string     `db:"phone"      json:"phone"`
	FullName  string     `db:"full_name"  json:"full_name"`
	BirthDate *time.Time `db:"birth_date" json:"birth_date,omitempty"`
	Email     *string    `db:"email"      json:"email,omitempty"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
}
