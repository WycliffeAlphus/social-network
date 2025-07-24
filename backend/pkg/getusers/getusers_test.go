package getusers_test

import (
	"backend/internal/model"
	"backend/pkg/getusers"
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	createUsers := `
	CREATE TABLE users (
		id TEXT PRIMARY KEY,
		email TEXT,
		password TEXT,
		fname TEXT,
		lname TEXT,
		dob TEXT,
		imgurl TEXT,
		nickname TEXT,
		about TEXT,
		created_at TEXT,
		profileVisibility TEXT
	);`

	createFollowers := `
	CREATE TABLE followers (
		follower_id TEXT,
		followed_id TEXT
	);`

	_, err = db.Exec(createUsers)
	if err != nil {
		t.Fatalf("failed to create users table: %v", err)
	}
	_, err = db.Exec(createFollowers)
	if err != nil {
		t.Fatalf("failed to create followers table: %v", err)
	}

	return db
}

func insertTestUser(t *testing.T, db *sql.DB, user model.User) {
	_, err := db.Exec(`
		INSERT INTO users (id, email, password, fname, lname, dob, imgurl, nickname, about, created_at, profileVisibility)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		user.ID, user.Email, user.Password, user.FirstName, user.LastName,
		user.DOB.Format(time.RFC3339), user.ImgURL, user.Nickname, user.About,
		user.CreatedAt.Format(time.RFC3339), user.ProfileVisibility,
	)
	if err != nil {
		t.Fatalf("failed to insert test user: %v", err)
	}
}


func TestIsFollowing(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	insertTestUser(t, db, model.User{
		ID:        "1",
		Email:     "follower@example.com",
		FirstName: "Follower",
		LastName:  "User",
		DOB:       time.Now(),
		CreatedAt: time.Now(),
	})

	insertTestUser(t, db, model.User{
		ID:        "2",
		Email:     "followed@example.com",
		FirstName: "Followed",
		LastName:  "User",
		DOB:       time.Now(),
		CreatedAt: time.Now(),
	})

	_, err := db.Exec(`INSERT INTO followers (follower_id, followed_id) VALUES (?, ?)`, "1", "2")
	if err != nil {
		t.Fatalf("failed to insert follower relation: %v", err)
	}

	following, err := getusers.IsFollowing(db, "1", "2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !following {
		t.Errorf("expected user 1 to follow user 2")
	}
}
