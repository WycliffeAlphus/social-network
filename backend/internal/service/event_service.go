package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"fmt"
)

type EventService struct {
	EventRepo *repository.EventRepository
	GroupRepo *repository.GroupRepository
}

// NewEventService creates and returns a new instance of EventService.
func NewEventService(eventRepo *repository.EventRepository, groupRepo *repository.GroupRepository) *EventService {
	return &EventService{
		EventRepo: eventRepo,
		GroupRepo: groupRepo,
	}
}

// CreateEvent creates a new event for a group.
func (s *EventService) CreateEvent(event *model.Event, userID string) (*model.Event, error) {
	// Verify that the user is a member of the group
	isMember, status, err := s.GroupRepo.CheckUserMembership(event.GroupID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember || status != "active" {
		return nil, fmt.Errorf("only group members can create events")
	}

	// Create the event
	eventID, err := s.EventRepo.CreateEvent(event)
	if err != nil {
		return nil, err
	}
	event.ID = eventID

	return event, nil
}

// GetGroupEvents retrieves all events for a group with response counts.
func (s *EventService) GetGroupEvents(groupID uint, userID string) ([]model.EventWithResponses, error) {
	// Verify that the user is a member of the group
	isMember, status, err := s.GroupRepo.CheckUserMembership(groupID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember || status != "active" {
		return nil, fmt.Errorf("only group members can view events")
	}

	// Get events with response counts
	return s.EventRepo.GetGroupEventsWithResponses(groupID, userID)
}

// RespondToEvent allows a user to respond to an event (going/not_going).
func (s *EventService) RespondToEvent(eventID uint, userID string, response string) error {
	// Validate response
	if response != "going" && response != "not_going" {
		return fmt.Errorf("invalid response. Must be 'going' or 'not_going'")
	}

	// Get event to check group membership
	event, err := s.EventRepo.GetEventByID(eventID)
	if err != nil {
		return err
	}
	if event == nil {
		return fmt.Errorf("event not found")
	}

	// Verify that the user is a member of the group
	isMember, status, err := s.GroupRepo.CheckUserMembership(event.GroupID, userID)
	if err != nil {
		return err
	}
	if !isMember || status != "active" {
		return fmt.Errorf("only group members can respond to events")
	}

	// Create or update the response
	return s.EventRepo.CreateOrUpdateEventResponse(eventID, userID, response)
}

// GetEventDetails retrieves event details with response counts.
func (s *EventService) GetEventDetails(eventID uint, userID string) (*model.EventWithResponses, error) {
	// Get event
	event, err := s.EventRepo.GetEventByID(eventID)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, fmt.Errorf("event not found")
	}

	// Verify that the user is a member of the group
	isMember, status, err := s.GroupRepo.CheckUserMembership(event.GroupID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember || status != "active" {
		return nil, fmt.Errorf("only group members can view event details")
	}

	// Get response counts
	counts, err := s.EventRepo.GetEventResponses(eventID)
	if err != nil {
		return nil, err
	}

	// Get user's response
	userResponse, err := s.EventRepo.GetUserEventResponse(eventID, userID)
	if err != nil {
		return nil, err
	}

	eventWithResponses := &model.EventWithResponses{
		Event:         *event,
		GoingCount:    counts["going"],
		NotGoingCount: counts["not_going"],
		UserResponse:  userResponse,
	}

	return eventWithResponses, nil
}

// DeleteEvent deletes an event (only creator can delete).
func (s *EventService) DeleteEvent(eventID uint, userID string) error {
	// Get event to check creator
	event, err := s.EventRepo.GetEventByID(eventID)
	if err != nil {
		return err
	}
	if event == nil {
		return fmt.Errorf("event not found")
	}

	// Verify that the user is the creator of the event
	if event.CreatorID != userID {
		return fmt.Errorf("only the event creator can delete the event")
	}

	// Delete the event
	return s.EventRepo.DeleteEvent(eventID)
}
