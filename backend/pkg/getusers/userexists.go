package getusers

import (
	"database/sql"
	"log"
	"net/http"
)

type ReceiverData struct {
	Firstname string         `json:"firstname"`
	Lastname  string         `json:"lastname"`
	Avatar    sql.NullString `json:"avatar"`
}

func UserExists(db *sql.DB, id string, w http.ResponseWriter) (ReceiverData, bool) {
	// check if selected user exists and get their data
	var receiverData ReceiverData
	var exists bool

	checkExistsErr := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)`, id).Scan(&exists)
	if checkExistsErr != nil {
		log.Println("Error checking if message receiver exists: ", checkExistsErr)
		http.Error(w, "An error occured, please check back later", http.StatusInternalServerError)
		return receiverData, false
	}

	if !exists {
		return receiverData, false
	}

	queryErr := db.QueryRow(`
		SELECT fname, lname, imgurl
		FROM users 
		WHERE id = ?`, id).Scan(
		&receiverData.Firstname,
		&receiverData.Lastname,
		&receiverData.Avatar,
	)

	if queryErr != nil {
		log.Println("Error fetching receiver data: ", queryErr)
		http.Error(w, "An error occurred, please check back later", http.StatusInternalServerError)
		return receiverData, false
	}

	return receiverData, true
}
