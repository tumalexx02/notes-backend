package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

type Storage struct {
	db *sqlx.DB
}

func New() (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sqlx.Open(
		"postgres",
		"host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable",
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	err = initMigrations(db)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		db: db,
	}, nil
}

func initMigrations(db *sqlx.DB) error {
	const op = "storage.postgres.initMigrations"

	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = goose.Up(db.DB, "./migrations")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
