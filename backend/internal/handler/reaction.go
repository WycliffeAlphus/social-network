package handler

import (
	"backend/internal/context"
	"backend/internal/model"
	"backend/internal/repository"
	"database/sql"
	"encoding/json"
	"net/http"
)

type ReactionResponse struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	UserReaction string `json:"userReaction"` // "like", "dislike", or ""
	LikeCount   int    `json:"likeCount"`
	DislikeCount int   `json:"dislikeCount"`
}

func HandleReaction(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set content type for JSON response
		w.Header().Set("Content-Type", "application/json")
		
		// Get user from context
		user := context.MustGetUser(r.Context())
		
		// Get parameters
		postID := r.URL.Query().Get("post_id")
		reactionType := r.URL.Query().Get("reaction_type")
		
		// Validate input
		if postID == "" || reactionType == "" {
			sendErrorResponse(w, "Missing post_id or reaction_type", http.StatusBadRequest)
			return
		}
		
		if reactionType != "like" && reactionType != "dislike" {
			sendErrorResponse(w, "Invalid reaction_type. Must be 'like' or 'dislike'", http.StatusBadRequest)
			return
		}
		
		// Create reaction object
		reaction := model.Reaction{
			UserID: user.ID,
			PostID: postID,
			Type:   reactionType,
		}
		
		// Check if exact same reaction exists
		exists, err := repository.CheckIfReactionExist(&reaction, db)
		if err != nil {
			sendErrorResponse(w, "Failed to check reaction status", http.StatusInternalServerError)
			return
		}
		
		var userReaction string
		
		if exists {
			// Remove the reaction (toggle off)
			err := repository.DeleteReaction(&reaction, db)
			if err != nil {
				sendErrorResponse(w, "Failed to remove reaction", http.StatusInternalServerError)
				return
			}
			userReaction = "" // No reaction after deletion
		} else {
			// Check if user has a different reaction on this post
			hasReacted, err := repository.CheckIfUserAlreadyReacted(&reaction, db)
			if err != nil {
				sendErrorResponse(w, "Failed to check existing reactions", http.StatusInternalServerError)
				return
			}
			
			if hasReacted {
				// Update existing reaction to new type
				err := repository.UpdateReaction(&reaction, db)
				if err != nil {
					sendErrorResponse(w, "Failed to update reaction", http.StatusInternalServerError)
					return
				}
			} else {
				// Create new reaction
				err = repository.CreateReaction(&reaction, db)
				if err != nil {
					sendErrorResponse(w, "Failed to create reaction", http.StatusInternalServerError)
					return
				}
			}
			userReaction = reactionType
		}
		
		// Get updated reaction counts
		likes, dislikes, err := repository.GetReactionsCount(postID, db)
		if err != nil {
			sendErrorResponse(w, "Failed to get reaction counts", http.StatusInternalServerError)
			return
		}
		
		// Send successful response
		response := ReactionResponse{
			Success:      true,
			Message:      "Reaction updated successfully",
			UserReaction: userReaction,
			LikeCount:    likes,
			DislikeCount: dislikes,
		}
		
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	response := ReactionResponse{
		Success: false,
		Message: message,
	}
	
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}