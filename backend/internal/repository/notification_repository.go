package repository

import (
	"database/sql"

	"social-network/internal/model"
)

type NotificationRepository struct {
	DB *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{DB: db}
}

// Create inserts a new notification into the database.
func (r *NotificationRepository) Create(notification *model.Notification) error {
	// Implementation to be added
	return nil
}

// GetByUserID retrieves all notifications for a specific user.
func (r *NotificationRepository) GetByUserID(userID int) ([]*model.Notification, error) {
	// Implementation to be added
	return nil, nil
}

// MarkAsRead marks a single notification as read in the database.
func (r *NotificationRepository) MarkAsRead(notificationID int, userID int) error {
	// Implementation to be added
	return nil
}
