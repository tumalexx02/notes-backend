package postgres

import (
	"fmt"
	"main/internal/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

type Storage struct {
	db *sqlx.DB
}

func New(cfg *config.Config) (*Storage, error) {
	const op = "storage.postgres.New"

	// get source name from config
	sourceName := getSourceName(cfg)

	// open connection with db
	db, err := sqlx.Open(
		"postgres",
		sourceName,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	// init migrations
	err = initMigrations(db, cfg)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		db: db,
	}, nil
}

func initMigrations(db *sqlx.DB, cfg *config.Config) error {
	const op = "storage.postgres.initMigrations"

	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = goose.Up(db.DB, cfg.MigrationsPath)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// TODO: add resetting migrations if reset flag in config is set

	return nil
}

func getSourceName(cfg *config.Config) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Name,
		cfg.Postgres.SSLMode,
	)
}
