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
	IDDocNumber   *string    `db:"id_doc_number"      json:"id_doc_number,omitempty"`   // номер паспорта
	IDDocIssuedBy *string    `db:"id_doc_issued_by"   json:"id_doc_issued_by,omitempty"` // кем выдан
	IDDocType     string     `db:"id_doc_type"        json:"id_doc_type"`                // 'domestic' | 'international'
	IDDocIssuedAt *time.Time `db:"id_doc_issued_at"   json:"id_doc_issued_at,omitempty"` // дата выдачи (загран)
	IDDocValidUntil *time.Time `db:"id_doc_valid_until" json:"id_doc_valid_until,omitempty"` // срок действия (загран)
}
