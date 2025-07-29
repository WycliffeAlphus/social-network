package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"backend/internal/utils"
)

type CreateGroupRequest struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	PrivacySetting string `json:"privacy_setting"` // e.g., "public", "private", "secret"
}

// CreateGroup handles the POST /groups endpoint
func CreateGroup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		creatorID := uint(1) // Placeholder: Replace with actual authenticated user ID

		// This is for local testing convenience and should be handled by user registration.
		var userCount int
		err := db.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", creatorID).Scan(&userCount)
		if err != nil || userCount == 0 {
			log.Printf("Creator user with ID %d does not exist. Creating dummy user.", creatorID)

			hashedPassword, err := utils.HashPassword("password123")
			if err != nil {
				utils.RespondWithError(w, http.StatusInternalServerError, "Failed to hash password for dummy user")
				return
			}
			_, err = db.Exec("INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)",
				"dummyuser", "dummy@example.com", hashedPassword)
			if err != nil {
				log.Printf("Failed to create dummy user: %v", err)
				utils.RespondWithError(w, http.StatusInternalServerError, "Failed to ensure creator user exists")
				return
			}
		}

		var req CreateGroupRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
			return
		}

		if req.Title == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Group title is required")
			return
		}
		if req.PrivacySetting == "" {
			req.PrivacySetting = "private"
		}

		allowedPrivacy := map[string]bool{"public": true, "private": true, "secret": true}
		if !allowedPrivacy[req.PrivacySetting] {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid privacy setting. Must be 'public', 'private', or 'secret'.")
			return
		}

		tx, err := db.Begin()
		if err != nil {
			log.Printf("Failed to begin transaction: %v", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create group")
			return
		}
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
				panic(r) // re-throw panic after Rollback
			} else if err != nil { // err will be non-nil if any DB operation failed before Commit
				tx.Rollback()
			}
		}()

		// Insert into 'groups' table
		stmt, err := tx.Prepare(`INSERT INTO groups (title, description, creator_id, privacy_setting) VALUES (?, ?, ?, ?)`)
		if err != nil {
			log.Printf("Failed to prepare group insert statement: %v", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create group")
			return
		}
		defer stmt.Close()

		res, err := stmt.Exec(req.Title, req.Description, creatorID, req.PrivacySetting)
		if err != nil {
			log.Printf("Failed to execute group insert: %v", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create group")
			return
		}

		groupID, err := res.LastInsertId()
		if err != nil {
			log.Printf("Failed to get last insert ID for group: %v", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve group ID")
			return
		}

		// Insert into 'group_members' table (add creator as admin)
		memberStmt, err := tx.Prepare(`INSERT INTO group_members (group_id, user_id, role, status) VALUES (?, ?, ?, ?)`)
		if err != nil {
			log.Printf("Failed to prepare group member insert statement: %v", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to add creator to group")
			return
		}
		defer memberStmt.Close()

		_, err = memberStmt.Exec(groupID, creatorID, "admin", "active")
		if err != nil {
			log.Printf("Failed to execute group member insert: %v", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to add creator to group")
			return
		}

		// Commit the transaction
		if err = tx.Commit(); err != nil {
			log.Printf("Failed to commit transaction: %v", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to finalize group creation")
			return
		}

		// Respond with Success
		// Retrieve the newly created group data to return in the response
		var newGroup models.Group
		row := db.QueryRow("SELECT id, title, description, creator_id, privacy_setting, created_at, updated_at FROM groups WHERE id = ?", groupID)
		if err := row.Scan(
			&newGroup.ID,
			&newGroup.Title,
			&newGroup.Description,
			&newGroup.CreatorID,
			&newGroup.PrivacySetting,
			&newGroup.CreatedAt,
			&newGroup.UpdatedAt,
		); err != nil {
			log.Printf("Failed to retrieve newly created group: %v", err)
			utils.RespondWithJSON(w, http.StatusCreated, map[string]any{"id": groupID, "message": "Group created successfully, but failed to fetch details"})
			return
		}

		utils.RespondWithJSON(w, http.StatusCreated, newGroup)
	}
}
