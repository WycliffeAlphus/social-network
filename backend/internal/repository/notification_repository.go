package repository

import (
	"backend/internal/model"
	"database/sql"
)

type NotificationRepository struct {
	DB *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{DB: db}
}

// Create inserts a new notification into the database.
func (r *NotificationRepository) Create(notification *model.Notification) error {
	query := `
		INSERT INTO notifications (user_id, actor_id, type, content_id, post_id, message)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := r.DB.Exec(query, notification.UserID, notification.ActorID, notification.Type, notification.ContentID, notification.PostID, notification.Message)
	return err
}

// GetByUserID retrieves all notifications for a specific user, ordered by most recent.
func (r *NotificationRepository) GetByUserID(userID string) ([]*model.Notification, error) {
	query := `
		SELECT id, user_id, actor_id, type, content_id, post_id, message, is_read, created_at
		FROM notifications
		WHERE user_id = ?
		ORDER BY created_at DESC
	`
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*model.Notification
	for rows.Next() {
		var n model.Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.ActorID, &n.Type, &n.ContentID, &n.PostID, &n.Message, &n.IsRead, &n.CreatedAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, &n)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return notifications, nil
}

// MarkAsRead marks a single notification as read in the database.
func (r *NotificationRepository) MarkAsRead(notificationID int, userID string) error {
	query := `UPDATE notifications SET is_read = 1 WHERE id = ? AND user_id = ?`
	_, err := r.DB.Exec(query, notificationID, userID)
	return err
}

// MarkAllAsRead marks all notifications for a user as read.
func (r *NotificationRepository) MarkAllAsRead(userID string) error {
	query := `UPDATE notifications SET is_read = 1 WHERE user_id = ?`
	_, err := r.DB.Exec(query, userID)
	return err
}
