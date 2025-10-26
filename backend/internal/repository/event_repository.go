package repository

import (
	"database/sql"

	"backend/internal/model"
)

type EventRepository struct {
	DB *sql.DB
}

// NewEventRepository creates and returns a new instance of EventRepository.
func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{DB: db}
}

// CreateEvent inserts a new event into the database.
func (r *EventRepository) CreateEvent(event *model.Event) (uint, error) {
	stmt, err := r.DB.Prepare(`
		INSERT INTO events (group_id, creator_id, title, description, event_time)
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(event.GroupID, event.CreatorID, event.Title, event.Description, event.EventTime)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// GetEventByID retrieves an event by its ID.
func (r *EventRepository) GetEventByID(eventID uint) (*model.Event, error) {
	var event model.Event
	err := r.DB.QueryRow(`
		SELECT id, group_id, creator_id, title, description, event_time, created_at, updated_at
		FROM events
		WHERE id = ? AND deleted_at IS NULL
	`, eventID).Scan(&event.ID, &event.GroupID, &event.CreatorID, &event.Title, &event.Description, &event.EventTime, &event.CreatedAt, &event.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Event not found
		}
		return nil, err
	}
	return &event, nil
}

// GetGroupEvents retrieves all events for a group.
func (r *EventRepository) GetGroupEvents(groupID uint) ([]model.Event, error) {
	rows, err := r.DB.Query(`
		SELECT id, group_id, creator_id, title, description, event_time, created_at, updated_at
		FROM events
		WHERE group_id = ? AND deleted_at IS NULL
		ORDER BY event_time ASC
	`, groupID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []model.Event
	for rows.Next() {
		var event model.Event
		err := rows.Scan(&event.ID, &event.GroupID, &event.CreatorID, &event.Title, &event.Description, &event.EventTime, &event.CreatedAt, &event.UpdatedAt)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

// CreateOrUpdateEventResponse creates or updates a user's response to an event.
func (r *EventRepository) CreateOrUpdateEventResponse(eventID uint, userID string, response string) error {
	_, err := r.DB.Exec(`
		INSERT INTO event_responses (event_id, user_id, response)
		VALUES (?, ?, ?)
		ON CONFLICT(event_id, user_id) DO UPDATE SET response = ?, updated_at = CURRENT_TIMESTAMP
	`, eventID, userID, response, response)
	return err
}

// GetEventResponses retrieves all responses for an event.
func (r *EventRepository) GetEventResponses(eventID uint) (map[string]int, error) {
	rows, err := r.DB.Query(`
		SELECT response, COUNT(*) as count
		FROM event_responses
		WHERE event_id = ? AND deleted_at IS NULL
		GROUP BY response
	`, eventID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := map[string]int{
		"going":     0,
		"not_going": 0,
	}

	for rows.Next() {
		var response string
		var count int
		err := rows.Scan(&response, &count)
		if err != nil {
			return nil, err
		}
		counts[response] = count
	}

	return counts, nil
}

// GetUserEventResponse retrieves a user's response to an event.
func (r *EventRepository) GetUserEventResponse(eventID uint, userID string) (string, error) {
	var response string
	err := r.DB.QueryRow(`
		SELECT response
		FROM event_responses
		WHERE event_id = ? AND user_id = ? AND deleted_at IS NULL
	`, eventID, userID).Scan(&response)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // No response found
		}
		return "", err
	}
	return response, nil
}

// GetGroupEventsWithResponses retrieves all events for a group with response counts and user's response.
func (r *EventRepository) GetGroupEventsWithResponses(groupID uint, userID string) ([]model.EventWithResponses, error) {
	rows, err := r.DB.Query(`
		SELECT
			e.id, e.group_id, e.creator_id, e.title, e.description, e.event_time, e.created_at, e.updated_at,
			COALESCE(SUM(CASE WHEN er.response = 'going' THEN 1 ELSE 0 END), 0) as going_count,
			COALESCE(SUM(CASE WHEN er.response = 'not_going' THEN 1 ELSE 0 END), 0) as not_going_count
		FROM events e
		LEFT JOIN event_responses er ON e.id = er.event_id AND er.deleted_at IS NULL
		WHERE e.group_id = ? AND e.deleted_at IS NULL
		GROUP BY e.id, e.group_id, e.creator_id, e.title, e.description, e.event_time, e.created_at, e.updated_at
		ORDER BY e.event_time ASC
	`, groupID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []model.EventWithResponses
	for rows.Next() {
		var event model.EventWithResponses
		err := rows.Scan(
			&event.ID, &event.GroupID, &event.CreatorID, &event.Title, &event.Description,
			&event.EventTime, &event.CreatedAt, &event.UpdatedAt,
			&event.GoingCount, &event.NotGoingCount,
		)
		if err != nil {
			return nil, err
		}

		// Get user's response for this event
		userResponse, err := r.GetUserEventResponse(event.ID, userID)
		if err != nil {
			return nil, err
		}
		event.UserResponse = userResponse

		events = append(events, event)
	}

	return events, nil
}

// DeleteEvent soft deletes an event.
func (r *EventRepository) DeleteEvent(eventID uint) error {
	_, err := r.DB.Exec(`
		UPDATE events
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = ? AND deleted_at IS NULL
	`, eventID)
	return err
}
