package sqlite

import (
	"backend/internal/model"
	"database/sql"
)

func UpdateUserVisibility(db *sql.DB, userID string, user model.User) error {
	query := `
		UPDATE users SET
			profileVisibility = COALESCE(?, profileVisibility)
		WHERE id = ?
	`
	_, err := db.Exec(query, user.ProfileVisibility, userID)
	return err
}
