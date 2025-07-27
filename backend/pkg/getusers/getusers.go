package getusers

import (
	"backend/internal/model"
	"database/sql"
	"time"
)

func GetUserByEmail(db *sql.DB, email string) (model.User, error) {
	query := `SELECT id, email, password, fname, lname, dob, 
		imgurl, nickname, about, created_at 
		FROM users WHERE email = ?`

	row := db.QueryRow(query, email)

	var user model.User
	var dob string
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&dob,
		&user.ImgURL,
		&user.Nickname,
		&user.About,
		&user.CreatedAt,
	)
	if err != nil {
		return user, err
	}

	user.DOB, err = time.Parse(time.RFC3339, dob)
	if err != nil {
		return user, err
	}

	return user, nil
}

func GetUserByID(db *sql.DB, userID string) (model.User, error) {
	query := `SELECT id, email, password, fname, lname, dob,
		imgurl, nickname, about, created_at, profileVisibility
		FROM users WHERE id = ?`

	row := db.QueryRow(query, userID)

	var user model.User
	var dob string
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&dob,
		&user.ImgURL,
		&user.Nickname,
		&user.About,
		&user.CreatedAt,
		&user.ProfileVisibility,
	)
	if err != nil {
		return user, err
	}

	user.DOB, err = time.Parse(time.RFC3339, dob)
	if err != nil {
		return user, err
	}

	return user, nil
}
