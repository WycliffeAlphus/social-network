package sqlite

import (
	"backend/internal/model"
	"database/sql"
)

func UpdateUserProfile(db *sql.DB, userID string, user model.User) error {
	query := `
		UPDATE users SET
			fname = COALESCE(?, fname),
			lname = COALESCE(?, lname),
			dob = COALESCE(?, dob),
			imgurl = COALESCE(?, imgurl),
			nickname = COALESCE(?, nickname),
			about = COALESCE(?, about),
			profileVisibility = COALESCE(?, profileVisibility)
		WHERE id = ?
	`
	_, err := db.Exec(query, user.FirstName, user.LastName, user.DOB, user.ImgURL, user.Nickname, user.About, user.ProfileVisibility, userID)
	return err
}
