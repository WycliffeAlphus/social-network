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
	// userID, ok := r.Context().Value("userID").(int)
	// if !ok {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }

	// notifications, err := h.service.GetByUserID(userID)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(notifications)
}

// MarkNotificationsAsRead handles marking notifications as read.
func (h *NotificationHandler) MarkNotificationsAsRead(w http.ResponseWriter, r *http.Request) {
	// Implementation to be added
	w.WriteHeader(http.StatusNotImplemented)
}
