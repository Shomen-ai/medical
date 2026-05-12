// Файл: internal/model/user.go
// Назначение: модель пациента (таблица users) с контактами и данными, фиксируемыми в медкарте туркменской клиники.
package model

import "time"

// User — пациент клиники: контакты, базовые медданные и удостоверение личности.
type User struct {
	ID            string     `db:"id"                 json:"id"`
	Phone         string     `db:"phone"              json:"phone"`
	FullName      string     `db:"full_name"          json:"full_name"`
	BirthDate     *time.Time `db:"birth_date"         json:"birth_date,omitempty"`
	Email         *string    `db:"email"              json:"email,omitempty"`
	CreatedAt     time.Time  `db:"created_at"         json:"created_at"`
	Gender        *string    `db:"gender"             json:"gender,omitempty"`           // 'm' | 'f' | nil
	Address       *string    `db:"address"            json:"address,omitempty"`          // адрес проживания
	IDDocNumber   *string    `db:"id_doc_number"      json:"id_doc_number,omitempty"`   // номер паспорта ТМ
	IDDocIssuedBy *string    `db:"id_doc_issued_by"   json:"id_doc_issued_by,omitempty"` // кем выдан
}
