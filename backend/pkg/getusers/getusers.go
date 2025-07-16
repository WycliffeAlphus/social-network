package getusers

import (
	"database/sql"
	"backend/pkg/models"
	"time"
)

func GetUserByEmail(db *sql.DB, email string) (models.User, error) {
	query := `SELECT id, email, password, fname, lname, dob, 
		imgurl, nickname, about, created_at 
		FROM users WHERE email = ?`

	row := db.QueryRow(query, email)

	var user models.User
	var dob string
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&dob,
		&user.AvatarImage,
		&user.Nickname,
		&user.AboutMe,
		&user.CreatedAt,
	)
	if err != nil {
		return user, err
	}

	user.DateOfBirth, err = time.Parse(time.RFC3339, dob)
	if err != nil {
		return user, err
	}

	return user, nil
}
