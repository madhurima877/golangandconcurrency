package repository

import (
	"database/sql"
	"errors"

	"distributed-ticketing/internal/model"

	"github.com/lib/pq"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUser(email string) (*model.User, error) {
	var user model.User

	err := r.db.QueryRow(`
		INSERT INTO users (email)
		VALUES ($1)
		RETURNING id, email, created_at
	`, email).Scan(
		&user.ID,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, errors.New("user with this email already exists")
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.QueryRow(`SELECT id, email, created_at
		FROM users
		WHERE email = $1`).Scan(&user.ID, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
