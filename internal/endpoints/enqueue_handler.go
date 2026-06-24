package endpoints

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AliKefall/Somnambulist/internal/services/matchmaking"
	"github.com/google/uuid"
)

func (cfg *Config) EnqueueHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		TimeControl string `json:"time_control"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "enqueue_error", "Invalid request body", "", err)
		return
	}

	userID := ctx.Value(UserIDKey).(uuid.UUID)
	if userID == uuid.Nil {
		RespondWithError(w, http.StatusUnauthorized, "context_error", "Missing context key", "", nil)
		return
	}

	user, err := cfg.Queries.GetUserByID(ctx, userID)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "database_error", "User not found", "", nil)
		return
	}

	entry := matchmaking.QueueEntry{
		UserID:      user.ID,
		Username:    user.Username,
		Rating:      1500, // This is only a placeholder it will change when I apply a rating system
		TimeControl: req.TimeControl,
		JoinedAt:    time.Now().UTC(),
	}

	cfg.MatchmakingService.EnqueuePlayer(ctx, entry)

	RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Player in queue",
	})
}
