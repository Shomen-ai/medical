// Файл: internal/model/specialty.go
// Назначение: модель медицинской специальности с длительностью базового слота приёма.
package model

// Specialty — медицинская специальность (кардиология, дерматология и т. п.) с шагом расписания.
type Specialty struct {
	ID              string `db:"id"                json:"id"`
	Name            string `db:"name"              json:"name"`
	SlotDurationMin int    `db:"slot_duration_min" json:"slot_duration_min"`
}
