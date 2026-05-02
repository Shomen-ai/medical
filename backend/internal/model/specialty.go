package model

type Specialty struct {
	ID              string `db:"id"                json:"id"`
	Name            string `db:"name"              json:"name"`
	SlotDurationMin int    `db:"slot_duration_min" json:"slot_duration_min"`
}
