package model

type Staff struct {
	ID       string  `db:"id"        json:"id"`
	DoctorID *string `db:"doctor_id" json:"doctor_id,omitempty"`
	Phone    string  `db:"phone"     json:"-"`
	Role     string  `db:"role"      json:"role"`
	IsActive bool    `db:"is_active" json:"is_active"`
}
