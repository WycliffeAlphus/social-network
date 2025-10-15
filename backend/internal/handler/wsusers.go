package handler

import (
	"backend/internal/context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"
)

type Envelope struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// UserStatus represents a user's information including their online status
type UserStatus struct {
	ID        string `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Status    string `json:"status"`
}

// getUserStatuses retrieves and sorts user statuses for a given user
func getUserStatuses(db *sql.DB, requestedUserID string) ([]UserStatus, error) {
	// get users that the requested user has chatted with their last message timestamps
	rows, err := db.Query(`
		SELECT 
            u.id, 
            u.fname, 
            u.lname,
            MAX(m.created_at) as last_message_time
        FROM users u
        JOIN messages m ON (u.id = m.sender_id OR u.id = m.receiver_id)
        WHERE (m.sender_id = ? OR m.receiver_id = ?) 
        AND u.id != ?
        GROUP BY u.id, u.fname, u.lname`,
		requestedUserID, requestedUserID, requestedUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	type userWithTime struct {
		id        string
		firstname string
		lastname  string
		lastTime  time.Time
	}

	var userStatuses []userWithTime
	for rows.Next() {
		var id, fname, lname string
		var lastTimeStr sql.NullString // NULL case when no messages exist

		if err := rows.Scan(&id, &fname, &lname, &lastTimeStr); err != nil {
			return nil, fmt.Errorf("failed to scan user statuses row: %w", err)
		}

		var lastTime time.Time
		if lastTimeStr.Valid {
			lastTime, _ = time.Parse("2006-01-02 15:04:05", lastTimeStr.String)
		}

		userStatuses = append(userStatuses, userWithTime{
			id:        id,
			firstname: fname,
			lastname:  lname,
			lastTime:  lastTime, // converts NullTime to time.Time, if null this will be a zero time
		})
	}
	fmt.Println("userStatuses are: ", userStatuses)

	// sort by last message time (newest first), then alphabetically
	sort.Slice(userStatuses, func(i, j int) bool {
		ti := userStatuses[i].lastTime
		tj := userStatuses[j].lastTime

		// if both have non-zero times, sort by time descending
		if !ti.IsZero() && !tj.IsZero() && !ti.Equal(tj) {
			return ti.After(tj) // newest message first
		}

		// if both times are zero or equal, sort by full name
		nameI := userStatuses[i].firstname + " " + userStatuses[i].lastname
		nameJ := userStatuses[j].firstname + " " + userStatuses[j].lastname
		return nameI < nameJ
	})

	// convert to final status list
	var result []UserStatus
	for _, user := range userStatuses {
		status := "offline"
		if _, ok := users[user.id]; ok {
			status = "online"
		}
		result = append(result, UserStatus{
			ID:        user.id,
			Firstname: user.firstname,
			Lastname:  user.lastname,
			Status:    status,
		})
	}

	return result, nil
}

func HandleUserStatuses(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed. Use GET for profile.", http.StatusMethodNotAllowed)
			return
		}

		currentUserID := context.MustGetUser(r.Context()).ID

		result, err := getUserStatuses(db, currentUserID)
		if err != nil {
			log.Println("error getting user statuses:", err)
			http.Error(w, "Failed to get user statuses", http.StatusInternalServerError)
			return
		}

		fmt.Println("result is: ", result)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}

func BroadcastUserList(db *sql.DB) {
	for connectedUserID, conn := range users {
		result, err := getUserStatuses(db, connectedUserID)
		if err != nil {
			log.Println("Failed to get user statuses:", err)
			continue
		}

		payload := Envelope{
			Type: "userlist",
			Data: result,
		}

		if err := conn.WriteJSON(payload); err != nil {
			log.Println("Failed to send user list update:", err)
			conn.Close()
			delete(users, connectedUserID)
		}
	}
}
