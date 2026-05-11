package model

import "time"

type User struct {
	ID               string     `db:"id"                  json:"id"`
	Phone            string     `db:"phone"               json:"phone"`
	FullName         string     `db:"full_name"           json:"full_name"`
	BirthDate        *time.Time `db:"birth_date"          json:"birth_date,omitempty"`
	Email            *string    `db:"email"               json:"email,omitempty"`
	CreatedAt        time.Time  `db:"created_at"          json:"created_at"`
	INN              *string    `db:"inn"                 json:"inn,omitempty"`
	PassportSeries   *string    `db:"passport_series"     json:"passport_series,omitempty"`
	PassportNumber   *string    `db:"passport_number"     json:"passport_number,omitempty"`
	PassportIssuedAt *time.Time `db:"passport_issued_at"  json:"passport_issued_at,omitempty"`
	PassportIssuedBy *string    `db:"passport_issued_by"  json:"passport_issued_by,omitempty"`
}
