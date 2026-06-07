// Файл: internal/repository/user.go
// Назначение: SQL-доступ к таблице пациентов (users) — поиск, создание (upsert), обновление профиля.
package repository

import (
	"beautymed/internal/model"

	"github.com/jmoiron/sqlx"
)

// UserRepo — репозиторий профилей пациентов.
type UserRepo struct{ db *sqlx.DB }

// NewUserRepo создаёт UserRepo на базе пула sqlx.
func NewUserRepo(db *sqlx.DB) *UserRepo { return &UserRepo{db} }

// FindByPhone ищет пациента по номеру телефона.
func (r *UserRepo) FindByPhone(phone string) (*model.User, error) {
	var u model.User
	err := r.db.Get(&u, `SELECT * FROM users WHERE phone = $1`, phone)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// Create создаёт нового пациента по телефону (upsert по конфликту), возвращая полную запись.
func (r *UserRepo) Create(phone string) (*model.User, error) {
	var u model.User
	err := r.db.Get(&u, `
		INSERT INTO users (phone) VALUES ($1)
		ON CONFLICT (phone) DO UPDATE SET phone = EXCLUDED.phone
		RETURNING *`, phone)
	return &u, err
}

// Update обновляет личные данные пациента (ФИО, дата рождения, контакты, пол, адрес, удостоверение личности).
func (r *UserRepo) Update(u *model.User) error {
	_, err := r.db.Exec(`
		UPDATE users SET
		    full_name         = $1,
		    birth_date        = $2,
		    email             = $3,
		    gender            = $4,
		    address           = $5,
		    id_doc_number      = $6,
		    id_doc_issued_by   = $7,
		    id_doc_type        = $8,
		    id_doc_issued_at   = $9,
		    id_doc_valid_until = $10
		WHERE id = $11`,
		u.FullName, u.BirthDate, u.Email,
		u.Gender, u.Address, u.IDDocNumber, u.IDDocIssuedBy,
		u.IDDocType, u.IDDocIssuedAt, u.IDDocValidUntil,
		u.ID)
	return err
}

// FindByID возвращает пациента по UUID.
func (r *UserRepo) FindByID(id string) (*model.User, error) {
	var u model.User
	err := r.db.Get(&u, `SELECT * FROM users WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
