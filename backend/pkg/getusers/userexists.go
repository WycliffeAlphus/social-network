package getusers

import (
	"database/sql"
	"log"
	"net/http"
)

func UserExists(db *sql.DB, id string, w http.ResponseWriter) bool {
	// check if selected user exists
	var exists bool
	checkExistsErr := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)`, id).Scan(&exists)
	if checkExistsErr != nil {
		log.Println("Error checking if message receiver exists: ", checkExistsErr)
		http.Error(w, "An error occured, please check back later", http.StatusInternalServerError)
		return false
	}
	return exists
}
