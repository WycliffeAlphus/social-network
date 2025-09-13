package handler

import (
	"backend/internal/context"
	"backend/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type NotificationHandler struct {
	service *service.NotificationService
}

func NewNotificationHandler(s *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: s}
}

// GetNotifications handles the request to fetch a user's notifications.
func (h *NotificationHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	user := context.MustGetUser(r.Context())

	userID := user.ID

	notifications, err := h.service.GetByUserID(userID)
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
	user := context.MustGetUser(r.Context())

	userID := user.ID

	err := h.service.MarkAllAsRead(userID)
	if err != nil {
		http.Error(w, "Failed to mark notifications as read", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Notifications marked as read"})
}

// MarkNotificationAsRead handles marking a single notification as read.
func (h *NotificationHandler) MarkNotificationAsRead(w http.ResponseWriter, r *http.Request) {
	user := context.MustGetUser(r.Context())
	userID := user.ID

	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	notificationIDStr := parts[3]
	notificationID, err := strconv.Atoi(notificationIDStr)
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	err = h.service.MarkAsRead(notificationID, userID)
	if err != nil {
		http.Error(w, "Failed to mark notification as read", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Notification marked as read"})
}
