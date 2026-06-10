package endpoints

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/AliKefall/Somnambulist/internal/database"
)

const ttl = 7 * 24 * time.Hour

func (cfg *Config) HandlerLogout(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	cookie, err := r.Cookie("refresh_token")

	clearRefreshToken(w, r)

	if err != nil || cookie.Value == "" {
		RespondWithJSON(w, http.StatusOK, map[string]string{
			"message": "logged out",
		})
		return
	}

	hash := sha256.Sum256([]byte(cookie.Value))
	tokenHash := hex.EncodeToString(hash[:])

	err = cfg.Queries.RevokeSessionByTokenHash(ctx, database.RevokeSessionByTokenHashParams{
		RefreshTokenHash: tokenHash,
		RevokedAt: sql.NullTime{
			Time: time.Now(),
			Valid: true,
		},
	})

	if err != nil && err != sql.ErrNoRows{
		RespondWithError(w, http.StatusInternalServerError, "logout_error", "logout failed", "", err)
		return
	}

	cfg.Redis.Set(ctx, "bl:" + tokenHash, "1", ttl)
	cfg.Redis.Del(ctx, "sess:"+tokenHash)

	RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "logged out",
	})
}
