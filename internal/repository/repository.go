package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Repo struct {
	db *sql.DB
}

func New() (Repo, error) {
	db, err := sql.Open("pgx", getPostgresUrl())
	if err != nil {
		return Repo{}, fmt.Errorf("open failed: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return Repo{}, fmt.Errorf("ping failed: %w", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return Repo{}, fmt.Errorf("with instance failed: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(os.Getenv("MIGRATE_PATH"), "postgres", driver)
	if err != nil {
		return Repo{}, fmt.Errorf("new with database instance failed: %w", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return Repo{}, fmt.Errorf("migrate up failed: %w", err)
	}

	return Repo{
		db: db,
	}, nil
}

func (r Repo) Close() error {
	return r.db.Close()
}

func getPostgresUrl() string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s", os.Getenv("POSTGRES_DB_USERNAME"), os.Getenv("POSTGRES_DB_PASSWORD"), os.Getenv("POSTGRES_DB_HOST"), os.Getenv("POSTGRES_DB_NAME"))
}
