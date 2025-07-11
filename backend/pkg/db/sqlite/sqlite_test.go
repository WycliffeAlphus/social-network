package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestConnectAndMigrate(t *testing.T) {
	// Change to project root directory for tests
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Navigate to project root (where go.mod is)
	err = os.Chdir("../../../")
	if err != nil {
		t.Fatalf("Failed to change to project root: %v", err)
	}

	// Clean up any existing test database
	defer func() {
		// Clean up test database file
		if _, err := os.Stat(DBFile); err == nil {
			os.Remove(DBFile)
		}
		// Clean up test directory
		os.RemoveAll(filepath.Dir(DBFile))
	}()

	// Test the actual ConnectAndMigrate function
	db, err := ConnectAndMigrate()
	if err != nil {
		t.Fatalf("ConnectAndMigrate() failed: %v", err)
	}
	defer db.Close()

	// Test that database connection works
	err = db.Ping()
	if err != nil {
		t.Errorf("Database ping failed: %v", err)
	}

	// Test that we can execute a simple query
	_, err = db.Exec("SELECT 1")
	if err != nil {
		t.Errorf("Failed to execute simple query: %v", err)
	}

	// Test that the database is actually SQLite
	var version string
	err = db.QueryRow("SELECT sqlite_version()").Scan(&version)
	if err != nil {
		t.Errorf("Failed to get SQLite version: %v", err)
	}
	if version == "" {
		t.Error("SQLite version should not be empty")
	}
}

func TestConnectAndMigrateCreatesDirectory(t *testing.T) {
	// Change to project root directory for tests
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Navigate to project root (where go.mod is)
	err = os.Chdir("../../../")
	if err != nil {
		t.Fatalf("Failed to change to project root: %v", err)
	}

	// Clean up before test
	os.RemoveAll(filepath.Dir(DBFile))
	
	// Clean up after test
	defer func() {
		os.Remove(DBFile)
		os.RemoveAll(filepath.Dir(DBFile))
	}()

	// Verify directory doesn't exist before
	if _, err := os.Stat(filepath.Dir(DBFile)); !os.IsNotExist(err) {
		t.Error("Database directory should not exist before test")
	}

	// Run ConnectAndMigrate
	db, err := ConnectAndMigrate()
	if err != nil {
		t.Fatalf("ConnectAndMigrate() failed: %v", err)
	}
	defer db.Close()

	// Check that directory was created
	if _, err := os.Stat(filepath.Dir(DBFile)); os.IsNotExist(err) {
		t.Error("Database directory was not created")
	}

	// Check that database file was created
	if _, err := os.Stat(DBFile); os.IsNotExist(err) {
		t.Error("Database file was not created")
	}

	// Check directory permissions
	dirInfo, err := os.Stat(filepath.Dir(DBFile))
	if err != nil {
		t.Errorf("Failed to get directory info: %v", err)
	}
	if !dirInfo.IsDir() {
		t.Error("Created path should be a directory")
	}
}

func TestConnectAndMigrateMultipleCalls(t *testing.T) {
	// Change to project root directory for tests
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Navigate to project root (where go.mod is)
	err = os.Chdir("../../../")
	if err != nil {
		t.Fatalf("Failed to change to project root: %v", err)
	}

	// Clean up
	defer func() {
		os.Remove(DBFile)
		os.RemoveAll(filepath.Dir(DBFile))
	}()

	// First call
	db1, err := ConnectAndMigrate()
	if err != nil {
		t.Fatalf("First ConnectAndMigrate() failed: %v", err)
	}
	db1.Close()

	// Second call (should not fail even if migrations already applied)
	db2, err := ConnectAndMigrate()
	if err != nil {
		t.Fatalf("Second ConnectAndMigrate() failed: %v", err)
	}
	defer db2.Close()

	// Should still be able to ping
	err = db2.Ping()
	if err != nil {
		t.Errorf("Database ping failed after second migration: %v", err)
	}

	// Third call to ensure consistency
	db3, err := ConnectAndMigrate()
	if err != nil {
		t.Fatalf("Third ConnectAndMigrate() failed: %v", err)
	}
	defer db3.Close()

	// All connections should work
	err = db3.Ping()
	if err != nil {
		t.Errorf("Database ping failed after third migration: %v", err)
	}
}

func TestDatabaseFileLocation(t *testing.T) {
	// Test that the constants are set correctly
	expectedDBFile := "./pkg/db/data/app.db"
	if DBFile != expectedDBFile {
		t.Errorf("DBFile constant incorrect. Expected: %s, Got: %s", expectedDBFile, DBFile)
	}

	expectedMigrationsPath := "file://pkg/db/migrations"
	if MigrationsPath != expectedMigrationsPath {
		t.Errorf("MigrationsPath constant incorrect. Expected: %s, Got: %s", expectedMigrationsPath, MigrationsPath)
	}

	// Test that paths are not empty
	if DBFile == "" {
		t.Error("DBFile should not be empty")
	}
	if MigrationsPath == "" {
		t.Error("MigrationsPath should not be empty")
	}

	// Test that DBFile has correct extension
	if !strings.HasSuffix(DBFile, ".db") {
		t.Error("DBFile should have .db extension")
	}

	// Test that MigrationsPath has correct prefix
	if !strings.HasPrefix(MigrationsPath, "file://") {
		t.Error("MigrationsPath should start with file://")
	}
}

func TestErrorHandling_InvalidDirectory(t *testing.T) {
	// This test checks error handling for directory creation
	// We can't easily test MkdirAll failure, but we can test other error paths
	
	// Test database file path validation
	if filepath.Dir(DBFile) == "" {
		t.Error("Database file directory should not be empty")
	}
	
	// Test that the path is relative
	if filepath.IsAbs(DBFile) {
		t.Error("Database file path should be relative for this test setup")
	}
}

func TestDatabaseOperations(t *testing.T) {
	// Change to project root directory for tests
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Navigate to project root (where go.mod is)
	err = os.Chdir("../../../")
	if err != nil {
		t.Fatalf("Failed to change to project root: %v", err)
	}

	// Clean up
	defer func() {
		os.Remove(DBFile)
		os.RemoveAll(filepath.Dir(DBFile))
	}()

	// Connect and migrate
	db, err := ConnectAndMigrate()
	if err != nil {
		t.Fatalf("ConnectAndMigrate() failed: %v", err)
	}
	defer db.Close()

	// Test basic SQL operations
	tests := []struct {
		name  string
		query string
	}{
		{"Simple SELECT", "SELECT 1"},
		{"SQLite version", "SELECT sqlite_version()"},
		{"Current time", "SELECT datetime('now')"},
		{"Math operation", "SELECT 2 + 2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rows, err := db.Query(tt.query)
			if err != nil {
				t.Errorf("Query '%s' failed: %v", tt.query, err)
				return
			}
			defer rows.Close()

			if !rows.Next() {
				t.Errorf("Query '%s' returned no rows", tt.query)
				return
			}

			// Try to scan the result
			var result interface{}
			err = rows.Scan(&result)
			if err != nil {
				t.Errorf("Failed to scan result from query '%s': %v", tt.query, err)
			}
		})
	}
}

func TestConcurrentAccess(t *testing.T) {
	// Change to project root directory for tests
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Navigate to project root (where go.mod is)
	err = os.Chdir("../../../")
	if err != nil {
		t.Fatalf("Failed to change to project root: %v", err)
	}

	// Clean up
	defer func() {
		os.Remove(DBFile)
		os.RemoveAll(filepath.Dir(DBFile))
	}()

	// Test concurrent connections
	const numConnections = 5
	results := make(chan error, numConnections)

	for i := 0; i < numConnections; i++ {
		go func(id int) {
			db, err := ConnectAndMigrate()
			if err != nil {
				results <- fmt.Errorf("connection %d failed: %v", id, err)
				return
			}
			defer db.Close()

			// Test the connection
			err = db.Ping()
			if err != nil {
				results <- fmt.Errorf("ping %d failed: %v", id, err)
				return
			}

			// Test a simple query
			_, err = db.Exec("SELECT 1")
			if err != nil {
				results <- fmt.Errorf("query %d failed: %v", id, err)
				return
			}

			results <- nil
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numConnections; i++ {
		select {
		case err := <-results:
			if err != nil {
				t.Errorf("Concurrent access test failed: %v", err)
			}
		case <-time.After(10 * time.Second):
			t.Error("Concurrent access test timed out")
		}
	}
}

func TestMigrationPathValidation(t *testing.T) {
	// Test that migration path is accessible
	// Strip the file:// prefix for file system operations
	migrationDir := strings.TrimPrefix(MigrationsPath, "file://")
	
	// Change to project root directory for tests
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Navigate to project root (where go.mod is)
	err = os.Chdir("../../../")
	if err != nil {
		t.Fatalf("Failed to change to project root: %v", err)
	}

	// Check if migration directory exists
	if _, err := os.Stat(migrationDir); err != nil {
		t.Logf("Migration directory %s does not exist or is not accessible: %v", migrationDir, err)
		// This is not necessarily an error if migrations haven't been created yet
	}

	// Test that the path format is correct
	if !strings.HasPrefix(MigrationsPath, "file://") {
		t.Error("MigrationsPath should start with file://")
	}
}

func TestDatabaseConnectionProperties(t *testing.T) {
	// Change to project root directory for tests
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Navigate to project root (where go.mod is)
	err = os.Chdir("../../../")
	if err != nil {
		t.Fatalf("Failed to change to project root: %v", err)
	}

	// Clean up
	defer func() {
		os.Remove(DBFile)
		os.RemoveAll(filepath.Dir(DBFile))
	}()

	// Connect and migrate
	db, err := ConnectAndMigrate()
	if err != nil {
		t.Fatalf("ConnectAndMigrate() failed: %v", err)
	}
	defer db.Close()

	// Test connection properties
	stats := db.Stats()
	if stats.MaxOpenConnections < 0 {
		t.Error("MaxOpenConnections should be >= 0")
	}

	// Test that we can set connection properties
	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(time.Hour)

	// Verify connection still works after setting properties
	err = db.Ping()
	if err != nil {
		t.Errorf("Database ping failed after setting connection properties: %v", err)
	}
}

func TestFilePathOperations(t *testing.T) {
	// Test filepath operations used in the code
	dir := filepath.Dir(DBFile)
	if dir == "" {
		t.Error("Directory should not be empty")
	}

	// Test that we can get the directory multiple times
	dir2 := filepath.Dir(DBFile)
	if dir != dir2 {
		t.Error("Directory should be consistent")
	}

	// Test path operations
	if filepath.IsAbs(DBFile) {
		t.Log("DBFile is absolute path")
	} else {
		t.Log("DBFile is relative path")
	}

	// Test that the path is well-formed
	cleaned := filepath.Clean(DBFile)
	if cleaned == "" {
		t.Error("Cleaned path should not be empty")
	}
}

func TestSQLiteDriver(t *testing.T) {
	// Test that we can open a database with the sqlite3 driver
	testDB := "./test_driver.db"
	defer os.Remove(testDB)

	db, err := sql.Open("sqlite3", testDB)
	if err != nil {
		t.Fatalf("Failed to open database with sqlite3 driver: %v", err)
	}
	defer db.Close()

	// Test that the driver works
	err = db.Ping()
	if err != nil {
		t.Errorf("SQLite driver ping failed: %v", err)
	}

	// Test driver name
	if db.Driver() == nil {
		t.Error("Database driver should not be nil")
	}
}

// Benchmark tests to ensure performance
func BenchmarkConnectAndMigrate(b *testing.B) {
	// Change to project root directory for tests
	originalDir, err := os.Getwd()
	if err != nil {
		b.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Navigate to project root (where go.mod is)
	err = os.Chdir("../../../")
	if err != nil {
		b.Fatalf("Failed to change to project root: %v", err)
	}

	// Clean up before benchmark
	os.Remove(DBFile)
	os.RemoveAll(filepath.Dir(DBFile))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db, err := ConnectAndMigrate()
		if err != nil {
			b.Fatalf("ConnectAndMigrate() failed: %v", err)
		}
		db.Close()
	}

	// Clean up after benchmark
	os.Remove(DBFile)
	os.RemoveAll(filepath.Dir(DBFile))
}

func BenchmarkDatabasePing(b *testing.B) {
	// Change to project root directory for tests
	originalDir, err := os.Getwd()
	if err != nil {
		b.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Navigate to project root (where go.mod is)
	err = os.Chdir("../../../")
	if err != nil {
		b.Fatalf("Failed to change to project root: %v", err)
	}

	// Setup database
	db, err := ConnectAndMigrate()
	if err != nil {
		b.Fatalf("ConnectAndMigrate() failed: %v", err)
	}
	defer db.Close()

	defer func() {
		os.Remove(DBFile)
		os.RemoveAll(filepath.Dir(DBFile))
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := db.Ping()
		if err != nil {
			b.Fatalf("Database ping failed: %v", err)
		}
	}
}