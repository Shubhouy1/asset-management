package database

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	// source/file import is required for migration files to read
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"

	// load pq as database driver
	_ "github.com/lib/pq"
)

var Asset *sqlx.DB

type SSLMode string

const (
	SSLModeDisabled SSLMode = "disable"
	SSLModeEnabled  SSLMode = "enable"
)

func CreateAndMigrate(host, port, user, password, dbname string, sslmode SSLMode) error {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)
	DataBase, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return err
	}
	err = DataBase.Ping()
	if err != nil {
		return err
	}
	Asset = DataBase
	return migrateUp(DataBase)
}
func migrateUp(db *sqlx.DB) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://database/migrations",
		"postgres", driver)

	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}
func Tx(fn func(tx *sqlx.Tx) error) (err error) {
	tx, err := Asset.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}

		if err != nil {
			_ = tx.Rollback()
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			err = fmt.Errorf("failed to commit tx: %w", commitErr)
		}
	}()
	err = fn(tx)
	return
}
