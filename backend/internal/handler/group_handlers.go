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
	Service *service.GroupService
}

func (h *GroupHandler) GetGroups(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user := context.MustGetUser(r.Context())
	userID := user.ID

	// Check if filter parameter is present
	filter := r.URL.Query().Get("filter")

	var groups []model.Group
	var err error

	if filter == "my" {
		// Get only groups where user is a member
		groups, err = h.Service.GetUserGroups(userID)
	} else {
		// Get all groups (default behavior)
		groups, err = h.Service.GetAllGroups()
	}

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
	err = h.Service.RequestToJoinGroup(groupID, userID)
	if err != nil {
		log.Printf("Failed to create join request: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
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

// GetPendingJoinRequests handles GET /groups/:id/join-requests endpoint.
// It allows group creators to view pending join requests for their groups.
func (h *GroupHandler) GetPendingJoinRequests(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user := context.MustGetUser(r.Context())
	creatorUserID := user.ID

	log.Printf("[DEBUG] GetPendingJoinRequests - user ID: %s", creatorUserID)

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

	log.Printf("[DEBUG] GetPendingJoinRequests - group ID: %d, creator ID: %s", groupID, creatorUserID)

	// Call service to get pending join requests
	requests, err := h.Service.GetPendingJoinRequests(groupID, creatorUserID)
	if err != nil {
		log.Printf("Failed to get pending join requests: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("[DEBUG] GetPendingJoinRequests - found %d requests", len(requests))

	utils.RespondWithJSON(w, http.StatusOK, requests)
}

// RejectJoinRequest handles POST /groups/:id/join endpoint with action=reject.
// It allows group creators to reject pending join requests.
func (h *GroupHandler) RejectJoinRequest(w http.ResponseWriter, r *http.Request) {
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

	// Parse request body to get the user ID to reject
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

	// Call service to reject join request
	err = h.Service.RejectJoinRequest(groupID, req.UserID, creatorUserID)
	if err != nil {
		log.Printf("Failed to reject join request: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Join request rejected successfully",
	})
}

// CheckGroupMembership handles GET /groups/:id/membership endpoint.
// It checks if the current user is a member of the group.
func (h *GroupHandler) CheckGroupMembership(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
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

	// Check membership
	isMember, err := h.Service.IsUserGroupMember(groupID, userID)
	if err != nil {
		log.Printf("Failed to check group membership: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to check membership")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]bool{
		"is_member": isMember,
	})
}

// InviteUserToGroup handles POST /groups/:id/invite endpoint.
// It allows group members to invite users to join a group.
func (h *GroupHandler) InviteUserToGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user := context.MustGetUser(r.Context())
	inviterUserID := user.ID

	if inviterUserID == "0" {
		utils.RespondWithError(w, http.StatusUnauthorized, "User ID not found or is invalid")
		return
	}

	// Extract group ID from URL path
	groupID, err := extractGroupIDFromPath(r.URL.Path)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid group ID")
		return
	}

	// Parse request body to get the user ID to invite
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

	// Call service to invite user
	err = h.Service.InviteUserToGroup(groupID, req.UserID, inviterUserID)
	if err != nil {
		log.Printf("Failed to invite user: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "User invited successfully",
	})
}

// AcceptGroupInvite handles POST /groups/:id/invites/accept endpoint.
// It allows invited users to accept an invitation to join a group.
func (h *GroupHandler) AcceptGroupInvite(w http.ResponseWriter, r *http.Request) {
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

	// Call service to accept invitation
	err = h.Service.AcceptGroupInvite(groupID, userID)
	if err != nil {
		log.Printf("Failed to accept invitation: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Invitation accepted successfully",
	})
}

// RejectGroupInvite handles POST /groups/:id/invites/reject endpoint.
// It allows invited users to reject an invitation to join a group.
func (h *GroupHandler) RejectGroupInvite(w http.ResponseWriter, r *http.Request) {
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

	// Call service to reject invitation
	err = h.Service.RejectGroupInvite(groupID, userID)
	if err != nil {
		log.Printf("Failed to reject invitation: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Invitation rejected successfully",
	})
}

// GetUserGroupInvites handles GET /groups/invites endpoint.
// It retrieves all pending group invitations for the authenticated user.
func (h *GroupHandler) GetUserGroupInvites(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user := context.MustGetUser(r.Context())
	userID := user.ID

	if userID == "0" {
		utils.RespondWithError(w, http.StatusUnauthorized, "User ID not found or is invalid")
		return
	}

	// Call service to get user's group invitations
	invites, err := h.Service.GetUserGroupInvites(userID)
	if err != nil {
		log.Printf("Failed to get group invitations: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve invitations")
		return
	}

	if invites == nil {
		invites = []model.Group{}
	}

	utils.RespondWithJSON(w, http.StatusOK, invites)
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
