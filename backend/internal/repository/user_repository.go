package repository

import (
	"backend/internal/model"
	"database/sql"
	"fmt"
	"log"
)

// UserRepository handles database operations for users
type UserRepository struct {
	DB *sql.DB // Database connection pool
}

// CreateUser inserts a new user record into the database
func (r *UserRepository) CreateUser(user *model.User) error {
	// SQL query to insert new user
	query := `INSERT INTO users (id, email, fname, lname, dob, imgurl, nickname, about, password, profileVisibility)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// Execute the query with user data
	_, err := r.DB.Exec(query,
		user.ID,
		user.Email,
		user.FirstName,
		user.LastName,
		user.DOB,
		user.ImgURL,
		user.Nickname,
		user.About,
		user.Password,
		user.ProfileVisibility,
	)
	return err // Return any error that occurred
}

// GetUserByEmail retrieves a user by their email address
func GetUserByEmail(db *sql.DB, email string) bool {
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM users WHERE email = ?`, email).Scan(&count)
	if err != nil {
		log.Println("Error querying database", err)
	}

	if count > 0 {
		return true
	}

	return false
}

// GetUserByNickname retrieves a user by their nickname
func GetUserByNickname(db *sql.DB, nickname string) bool {
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM users WHERE nickname = ?`, nickname).Scan(&count)
	if err != nil {
		log.Println("Error querying database", err)
	}

	if count > 0 {
		return true
	}

	return false
}

// GetUserByID retrieves a user by their ID
func (r *UserRepository) GetUserByID(id int) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id, email, fname, lname, dob, imgurl, nickname, about, password, profileVisibility FROM users WHERE id = ?`
	err := r.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.DOB,
		&user.ImgURL,
		&user.Nickname,
		&user.About,
		&user.Password,
		&user.ProfileVisibility,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with ID %d not found", id)
		}
		return nil, fmt.Errorf("error getting user by ID: %w", err)
	}
	return user, nil
}
