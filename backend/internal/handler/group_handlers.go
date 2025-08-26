package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

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
	Service             *service.GroupService
	NotificationService *service.NotificationService
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

// JoinGroupRequest handles POST /groups/:id/join endpoint.
// It allows users to request to join a group.
func (h *GroupHandler) JoinGroupRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user := context.MustGetUser(r.Context())
	userID := user.ID

	if userID == "0" {
		utils.RespondWithError(w, http.StatusUnauthorized, "User ID not found or is invalid")
		return
	}

	// Extract group ID from URL path
	groupID, err := extractGroupIDFromPath(r.URL.Path)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid group ID")
		return
	}

	// Call service to create join request
	group, err := h.Service.RequestToJoinGroup(groupID, userID)
	if err != nil {
		log.Printf("Failed to create join request: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Trigger notification for the group owner
	actorID, _ := strconv.Atoi(userID)
	groupOwnerID, _ := strconv.Atoi(group.CreatorID)
	if err := h.NotificationService.CreateGroupJoinRequestNotification(actorID, groupOwnerID, int(groupID)); err != nil {
		log.Printf("Failed to create group join request notification: %v", err)
		// Do not block response to user for notification failure
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Join request sent successfully",
	})
}

// AcceptJoinRequest handles POST /groups/:id/join endpoint with action=accept.
// It allows group creators to accept pending join requests.
func (h *GroupHandler) AcceptJoinRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user := context.MustGetUser(r.Context())
	creatorUserID := user.ID

	if creatorUserID == "0" {
		utils.RespondWithError(w, http.StatusUnauthorized, "User ID not found or is invalid")
		return
	}

	// Extract group ID from URL path
	groupID, err := extractGroupIDFromPath(r.URL.Path)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid group ID")
		return
	}

	// Parse request body to get the user ID to accept
	var req struct {
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	if req.UserID == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	// Call service to accept join request
	err = h.Service.AcceptJoinRequest(groupID, req.UserID, creatorUserID)
	if err != nil {
		log.Printf("Failed to accept join request: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Join request accepted successfully",
	})
}

// extractGroupIDFromPath extracts the group ID from URL paths like /groups/123/join
func extractGroupIDFromPath(path string) (uint, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 2 {
		return 0, fmt.Errorf("invalid path format")
	}

	// Find the groups segment and get the next part as ID
	for i, part := range parts {
		if part == "groups" && i+1 < len(parts) {
			id, err := strconv.ParseUint(parts[i+1], 10, 32)
			if err != nil {
				return 0, fmt.Errorf("invalid group ID format")
			}
			return uint(id), nil
		}
	}

	return 0, fmt.Errorf("group ID not found in path")
}

// InviteUserToGroup handles POST /groups/:id/invite
func (h *GroupHandler) InviteUserToGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	inviter := context.MustGetUser(r.Context())
	inviterID := inviter.ID

	groupID, err := extractGroupIDFromPath(r.URL.Path)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid group ID")
		return
	}

	var req struct {
		TargetUserID string `json:"target_user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.Service.InviteUserToGroup(groupID, inviterID, req.TargetUserID); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Trigger notification for the invited user
	inviterIDInt, _ := strconv.Atoi(inviterID)
	targetUserIDInt, _ := strconv.Atoi(req.TargetUserID)
	if err := h.NotificationService.CreateGroupInviteNotification(inviterIDInt, targetUserIDInt, int(groupID)); err != nil {
		log.Printf("Failed to create group invite notification: %v", err)
		// Do not block response to user for notification failure
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Invitation sent successfully"})
}

// CreateEvent handles POST /groups/:id/events
func (h *GroupHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	creator := context.MustGetUser(r.Context())
	creatorID, _ := strconv.Atoi(creator.ID)

	groupID, err := extractGroupIDFromPath(r.URL.Path)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid group ID")
		return
	}

	var event model.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	event.CreatorID = creatorID
	event.GroupID = int(groupID)

	createdEvent, err := h.Service.CreateEvent(&event)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Trigger notification for group members
	if err := h.NotificationService.CreateGroupEventNotification(creatorID, int(groupID), createdEvent.ID); err != nil {
		log.Printf("Failed to create group event notification: %v", err)
	}

	utils.RespondWithJSON(w, http.StatusCreated, createdEvent)
}
