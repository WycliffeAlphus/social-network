package handler

import (
	"backend/internal/context"
	"backend/internal/model"
	"backend/pkg/getusers"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func PrivateConversations(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUserId := context.MustGetUser(r.Context()).ID
		receiverId := r.URL.Query().Get("receiverId")

		exists := getusers.UserExists(db, receiverId, w)

		if receiverId == "" || !exists {
			http.Error(w, "Select a user to converse with", http.StatusBadRequest)
			return
		}

		query := `
        SELECT sender_id, receiver_id, content, created_at 
        FROM messages
        WHERE (sender_id = ? AND receiver_id = ?)
           OR (sender_id = ? AND receiver_id = ?)
        ORDER BY created_at DESC`

		rows, err := db.Query(query, currentUserId, receiverId, receiverId, currentUserId)
		if err != nil {
			log.Println("Failed to query private messages history: ", err)
			http.Error(w, "An error occured, please check back later", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var messages []model.Message
		for rows.Next() {
			var msg model.Message
			if err := rows.Scan(&msg.From, &msg.To, &msg.Content, &msg.Timestamp); err != nil {
				log.Println("Failed to scan messages details in conversation: %w", err)
			}
			messages = append(messages, msg)
		}

		// since the order was DESC for newest first, reverse them before sending
		for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
			messages[i], messages[j] = messages[j], messages[i]
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(messages)
	}
}
