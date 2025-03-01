package testinfra

import (
	"database/sql"
	"testing"

	"github.com/pressly/goose/v3"
)

func SetupTestDB(t *testing.T, migrationsPath string) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open sqlite database: %v", err)
	}

	if err := goose.SetDialect("sqlite"); err != nil {
		t.Fatalf("failed to set dialect: %v", err)
	}
	if err := goose.Up(db, migrationsPath); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	return db
}
