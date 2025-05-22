package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func Connect(driver, uri string) (*sql.DB, error) {
	database, err := sql.Open(driver, uri)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	if err := database.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	if err := migrateDatabase(database, driver); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}
	return database, nil
}

func migrateDatabase(db *sql.DB, driver string) error {
	if err := goose.SetDialect(driver); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}
	if err := goose.Up(db, "./migrations"); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}
