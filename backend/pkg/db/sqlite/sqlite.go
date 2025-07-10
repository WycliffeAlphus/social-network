package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

// DBFile is the location of your SQLite database file
const DBFile = "./pkg/db/data/app.db"

// MigrationsPath is the folder containing .sql migration files
const MigrationsPath = "file://pkg/db/migrations"

// ConnectAndMigrate opens the SQLite DB and runs migrations
func ConnectAndMigrate() (*sql.DB, error) {
	// Create data folder if not exists
	err := os.MkdirAll(filepath.Dir(DBFile), os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	// Open SQLite DB
	db, err := sql.Open("sqlite3", DBFile+"?_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite DB: %w", err)
	}

	// Ensure DB is reachable
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	// Setup golang-migrate with SQLite
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create SQLite driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(MigrationsPath, "sqlite3", driver)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// version is the migration version in the schema_migrations table
	// dirty is a boolean flag that indicates whether the last migration attempt failed
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return nil, fmt.Errorf("failed to get migration version: %v", err)
	}

	// Apply all up migrations
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	log.Printf("Migrations applied successfully with version %d, and %v dirty state.\n", version, dirty)
	return db, nil
}
