package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	dsn := "host=localhost port=15632 user=postgres password=postgres dbname=practice sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
