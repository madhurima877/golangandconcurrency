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
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	err = tx.QueryRow(`
		INSERT INTO users (email)
		VALUES ($1)
		RETURNING id, email, created_at
	`, email).Scan(
		&user.ID,
		&user.Email,
		&user.CreatedAt)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, errors.New("user with this email already exists")
		}
		return nil, err
	}
	_, err = tx.Exec(`INSERT INTO outbox_events (event_type,aggregate_id,payload) VALUES ($1,$2,jsonb_build_object('id',$2::bigint,'email',$3::text))`, "USER_CREATED", user.ID, user.Email)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
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
