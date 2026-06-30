package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *Config) DequeueHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := ctx.Value(UserIDKey).(uuid.UUID)
	if userID == uuid.Nil {
		RespondWithError(w, http.StatusUnauthorized, "context_error", "Missing userid", "", nil)
		return
	}

	var req struct {
		TimeControl string `json:"time_control"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "dequeue_error", "Json body could not be parsed", "", err)
		return
	}

	cfg.MatchmakingService.DequeuePlayer(ctx, userID, req.TimeControl)

	RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Dequeue successful",
	})
}
