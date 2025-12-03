package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"backend/internal/context"
	"backend/internal/model"
	"backend/internal/service"
	"backend/internal/utils"
)

// CreateEventRequest matches the expected JSON payload for creating an event.
type CreateEventRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	EventTime   string `json:"event_time"` // ISO 8601 format
}

// EventResponseRequest matches the expected JSON payload for responding to an event.
type EventResponseRequest struct {
	Response string `json:"response"` // 'going' or 'not_going'
}

// EventHandler holds the business logic service for events.
type EventHandler struct {
	Service *service.EventService
}

// CreateEvent handles POST /groups/:id/events endpoint.
// It allows group members to create events.
func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
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

	var req CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	if req.Title == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Event title is required")
		return
	}

	if req.EventTime == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Event time is required")
		return
	}

	// Parse event time
	eventTime, err := time.Parse(time.RFC3339, req.EventTime)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid event time format. Use ISO 8601 format (e.g., 2024-12-31T15:04:05Z)")
		return
	}

	event := &model.Event{
		GroupID:     groupID,
		CreatorID:   userID,
		Title:       req.Title,
		Description: req.Description,
		EventTime:   eventTime,
	}

	// Call the service layer to create the event
	newEvent, err := h.Service.CreateEvent(event, userID)
	if err != nil {
		log.Printf("Failed to create event: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, newEvent)
}

// GetGroupEvents handles GET /groups/:id/events endpoint.
// It retrieves all events for a group with response counts.
func (h *EventHandler) GetGroupEvents(w http.ResponseWriter, r *http.Request) {
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

	// Get events
	events, err := h.Service.GetGroupEvents(groupID, userID)
	if err != nil {
		log.Printf("Failed to get events: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if events == nil {
		events = []model.EventWithResponses{}
	}

	utils.RespondWithJSON(w, http.StatusOK, events)
}

// RespondToEvent handles POST /events/:id/respond endpoint.
// It allows group members to respond to an event (going/not_going).
func (h *EventHandler) RespondToEvent(w http.ResponseWriter, r *http.Request) {
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

	// Extract event ID from URL path
	eventID, err := extractEventIDFromPath(r.URL.Path)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	var req EventResponseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	if req.Response == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Response is required")
		return
	}

	// Call service to respond to event
	err = h.Service.RespondToEvent(eventID, userID, req.Response)
	if err != nil {
		log.Printf("Failed to respond to event: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Response recorded successfully",
	})
}

// GetEventDetails handles GET /events/:id endpoint.
// It retrieves event details with response counts.
func (h *EventHandler) GetEventDetails(w http.ResponseWriter, r *http.Request) {
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

	// Extract event ID from URL path
	eventID, err := extractEventIDFromPath(r.URL.Path)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	// Get event details
	event, err := h.Service.GetEventDetails(eventID, userID)
	if err != nil {
		log.Printf("Failed to get event details: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, event)
}

// DeleteEvent handles DELETE /events/:id endpoint.
// It allows event creators to delete events.
func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user := context.MustGetUser(r.Context())
	userID := user.ID

	if userID == "0" {
		utils.RespondWithError(w, http.StatusUnauthorized, "User ID not found or is invalid")
		return
	}

	// Extract event ID from URL path
	eventID, err := extractEventIDFromPath(r.URL.Path)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	// Call service to delete event
	err = h.Service.DeleteEvent(eventID, userID)
	if err != nil {
		log.Printf("Failed to delete event: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Event deleted successfully",
	})
}

// extractEventIDFromPath extracts the event ID from URL paths like /events/123/respond
func extractEventIDFromPath(path string) (uint, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 2 {
		return 0, ErrInvalidPath
	}

	// Find the events segment and get the next part as ID
	for i, part := range parts {
		if part == "events" && i+1 < len(parts) {
			id, err := strconv.ParseUint(parts[i+1], 10, 32)
			if err != nil {
				return 0, err
			}
			return uint(id), nil
		}
	}

	return 0, ErrEventIDNotFound
}

var (
	ErrInvalidPath = &customError{message: "invalid path format"}
	ErrEventIDNotFound = &customError{message: "event ID not found in path"}
)

type customError struct {
	message string
}

func (e *customError) Error() string {
	return e.message
}
