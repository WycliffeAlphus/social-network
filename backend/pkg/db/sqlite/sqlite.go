package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

// DBFile is the location of your SQLite database file
const DBFile = "./pkg/db/data/app.db"

// MigrationsPath is the folder containing .sql migration files
const MigrationsPath = "file://pkg/db/migrations"

// Session represents a user session.
type Session struct {
	ID        string // UUID v4
	UserID    string
	CreatedAt time.Time
	ExpiresAt time.Time
}

func CreateMigrationFile() {
	if len(os.Args) == 4 && os.Args[1] == "migrate" && os.Args[2] == "create" {
		fileName := os.Args[3]
		migrationsDir := "pkg/db/migrations"

		// Validate file name format (lowercase with underscores)
		validName := regexp.MustCompile(`^[a-z]+(_[a-z]+)*$`).MatchString
		if !validName(fileName) {
			fmt.Println("Error: Migration name must use lowercase letters and underscores only (e.g. 'add_followers_table')")
			os.Exit(1)
		}

		cmd := exec.Command("migrate", "create", "-ext", "sql", "-dir", migrationsDir, "-seq", fileName)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		fmt.Printf("Creating migration: %s\n", fileName)
		if err := cmd.Run(); err != nil {
			fmt.Printf("Error creating migration: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Migration files created successfully")
	} else if len(os.Args) > 1 && os.Args[1] != "migrate" && os.Args[2] != "create" {
		fmt.Println("Unrecognized CLI commands")
		fmt.Println()
		fmt.Println("Migration usage:")
		fmt.Println("go run main.go migrate create <migration-file-name> - Creates new migration files")
		os.Exit(1)
	}
}

// InsertSession inserts a new session into the sessions table.
func InsertSession(db *sql.DB, userID string, expiresAt time.Time) (string, error) {
	sessionID := uuid.NewString()
	_, err := db.Exec(
		"INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)",
		sessionID, userID, expiresAt,
	)
	if err != nil {
		return "", err
	}
	return sessionID, nil
}

// GetSession retrieves a session by its ID.
func GetSession(db *sql.DB, sessionID string) (*Session, error) {
	var s Session
	row := db.QueryRow(
		"SELECT id, user_id, created_at, expires_at FROM sessions WHERE id = ?",
		sessionID,
	)
	err := row.Scan(&s.ID, &s.UserID, &s.CreatedAt, &s.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// DeleteSession deletes a session by its ID.
func DeleteSession(db *sql.DB, sessionID string) error {
	_, err := db.Exec("DELETE FROM sessions WHERE id = ?", sessionID)
	return err
}

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
	if err != nil {
		if err == migrate.ErrNoChange {
			log.Println("No new migrations to apply")
		} else {
			return nil, fmt.Errorf("migration failed: %v", err)
		}
	}

	log.Printf("Migrations applied successfully with version %d, and %v dirty state.\n", version, dirty)
	return db, nil
}
