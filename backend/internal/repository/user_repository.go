package repository

import (
	"backend/internal/model"
	"database/sql"
	"errors"
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
func (r *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	// SQL query to select user by email
	query := `SELECT id, email, fname, lname, dob, imgurl, nickname, about, password, profileVisibility, created_at 
		FROM users WHERE email = ?`
	
	// Execute query for single row
	row := r.DB.QueryRow(query, email)
	
	// Prepare user object and nullable fields
	user := &model.User{}
	var nickname, imgurl, about sql.NullString // Fields that might be NULL in DB
	
	// Scan row data into variables
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.DOB,
		&imgurl,
		&nickname,
		&about,
		&user.Password,
		&user.ProfileVisibility,
		&user.CreatedAt,
	)
	
	// Handle errors
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found") // Special case for no results
		}
		return nil, err // Return other errors
	}
	
	// Handle nullable fields - only set if they contain valid data
	if nickname.Valid {
		user.Nickname = nickname.String
	}
	if imgurl.Valid {
		user.ImgURL = imgurl.String
	}
	if about.Valid {
		user.About = about.String
	}
	
	return user, nil
}

// GetUserByNickname retrieves a user by their nickname
func (r *UserRepository) GetUserByNickname(nickname string) (*model.User, error) {
	// SQL query to select user by nickname
	query := `SELECT id, email, fname, lname, dob, imgurl, nickname, about, password, profileVisibility, created_at 
		FROM users WHERE nickname = ?`
	
	// Execute query for single row
	row := r.DB.QueryRow(query, nickname)
	
	// Prepare user object and nullable fields
	user := &model.User{}
	var nicknameField, imgurl, about sql.NullString // Fields that might be NULL in DB
	
	// Scan row data into variables
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.DOB,
		&imgurl,
		&nicknameField,
		&about,
		&user.Password,
		&user.ProfileVisibility,
		&user.CreatedAt,
	)
	
	// Handle errors
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found") // Special case for no results
		}
		return nil, err // Return other errors
	}
	
	// Handle nullable fields - only set if they contain valid data
	if nicknameField.Valid {
		user.Nickname = nicknameField.String
	}
	if imgurl.Valid {
		user.ImgURL = imgurl.String
	}
	if about.Valid {
		user.About = about.String
	}
	
	return user, nil
}
