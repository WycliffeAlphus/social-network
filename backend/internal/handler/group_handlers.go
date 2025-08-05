package handler

import (
	"encoding/json"

	"log"
	"net/http"

	"backend/internal/context"
	"backend/internal/model"
	"backend/internal/service"
	"backend/internal/utils"
)

// CreateGroupRequest matches the expected JSON payload for creating a group.
type CreateGroupRequest struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	PrivacySetting string `json:"privacy_setting"` // e.g., "public", "private", "secret"
}

// GroupHandler holds the business logic service for groups.
type GroupHandler struct {
	Service *service.GroupService
}

func (h *GroupHandler) GetGroups(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	groups, err := h.Service.GetAllGroups()
	if err != nil {
		log.Printf("Failed to retrieve groups: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve groups")
		return
	}

	if groups == nil {
		groups = []model.Group{}
	}

	utils.RespondWithJSON(w, http.StatusOK, groups)
}

// CreateGroup handles the POST /api/groups endpoint.
// It retrieves the authenticated user's ID from the request context.
func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	user := context.MustGetUser(r.Context())

	// extract the ID
	creatorID := user.ID

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

	// Call the service layer to handle the business logic and database operations
	newGroup, err := h.Service.CreateGroup(req.Title, req.Description, req.PrivacySetting, creatorID)
	if err != nil {
		log.Printf("Failed to create group via service: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create group: "+err.Error())
		return
	}

	// Respond with Success
	utils.RespondWithJSON(w, http.StatusCreated, newGroup)
}
