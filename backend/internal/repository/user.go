package repository

import (
	"beautymed/internal/model"

	"github.com/jmoiron/sqlx"
)

type UserRepo struct{ db *sqlx.DB }

func NewUserRepo(db *sqlx.DB) *UserRepo { return &UserRepo{db} }

func (r *UserRepo) FindByPhone(phone string) (*model.User, error) {
	var u model.User
	err := r.db.Get(&u, `SELECT * FROM users WHERE phone = $1`, phone)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) Create(phone string) (*model.User, error) {
	var u model.User
	err := r.db.Get(&u, `
		INSERT INTO users (phone) VALUES ($1)
		ON CONFLICT (phone) DO UPDATE SET phone = EXCLUDED.phone
		RETURNING *`, phone)
	return &u, err
}

func (r *UserRepo) Update(u *model.User) error {
	_, err := r.db.Exec(`
		UPDATE users SET
		    full_name           = $1,
		    birth_date          = $2,
		    email               = $3,
		    inn                 = $4,
		    passport_series     = $5,
		    passport_number     = $6,
		    passport_issued_at  = $7,
		    passport_issued_by  = $8
		WHERE id = $9`,
		u.FullName, u.BirthDate, u.Email,
		u.INN, u.PassportSeries, u.PassportNumber,
		u.PassportIssuedAt, u.PassportIssuedBy,
		u.ID)
	return err
}

func (r *UserRepo) FindByID(id string) (*model.User, error) {
	var u model.User
	err := r.db.Get(&u, `SELECT * FROM users WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
