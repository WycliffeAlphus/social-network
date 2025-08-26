package handler

import (
	"encoding/json"
	"net/http"
	"social-network/internal/context"
	"social-network/internal/service"
)

type NotificationHandler struct {
	service *service.NotificationService
}

func NewNotificationHandler(s *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: s}
}

// GetNotifications handles the request to fetch a user's notifications.
func (h *NotificationHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	user := context.GetUserFromContext(r.Context())
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	notifications, err := h.service.GetByUserID(user.ID)
	if err != nil {
		http.Error(w, "Failed to retrieve notifications", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if notifications == nil {
		w.Write([]byte("[]")) // Return empty JSON array if no notifications
		return
	}
	json.NewEncoder(w).Encode(notifications)
}

// MarkNotificationsAsRead handles marking all notifications as read.
func (h *NotificationHandler) MarkNotificationsAsRead(w http.ResponseWriter, r *http.Request) {
	user := context.GetUserFromContext(r.Context())
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err := h.service.MarkAllAsRead(user.ID)
	if err != nil {
		http.Error(w, "Failed to mark notifications as read", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Notifications marked as read"})
}
