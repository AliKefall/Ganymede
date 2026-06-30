package endpoints

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/AliKefall/Somnambulist/internal/auth"
	"github.com/AliKefall/Somnambulist/internal/database"
	"github.com/google/uuid"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *Config) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var req LoginRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "login_error", "Invalid request body", "", err)
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if req.Email == "" || req.Password == "" {
		RespondWithError(w, http.StatusBadRequest, "login_error", "Email and password are required", "", nil)
		return
	}

	user, err := cfg.Queries.GetUserByEmail(ctx, req.Email)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "login_error", "Invalid credentials", "", err)
		return
	}

	ok, err := cfg.Hasher.Verify(user.Password, req.Password)
	if !ok || err != nil {
		RespondWithError(w, http.StatusUnauthorized, "login_error", "Invalid credentials", "", err)
		return
	}

	now := time.Now()
	refreshTTL := 7 * 24 * time.Hour
	maxSessionTTL := 30 * 24 * time.Hour

	refreshExpires := now.Add(refreshTTL)
	maxExpires := now.Add(maxSessionTTL)

	sessionID := uuid.New()

	refreshToken, err := auth.MakeRefreshToken()

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "refresh_error", "Failed to create refresh token", "", err)
		return
	}

	hash := sha256.Sum256([]byte(refreshToken))
	refreshHash := hex.EncodeToString(hash[:])

	_, err = cfg.Queries.CreateSession(ctx, database.CreateSessionParams{
		ID:               sessionID,
		UserID:           user.ID,
		RefreshTokenHash: refreshHash,
		UserAgent: sql.NullString{
			String: r.UserAgent(),
			Valid:  true,
		},
		IpAddress: sql.NullString{
			String: r.RemoteAddr,
			Valid:  true,
		},
		CreatedAt:    now,
		ExpiresAt:    refreshExpires,
		MaxExpiresAt: maxExpires,
		LastUsedAt: sql.NullTime{
			Time:  now,
			Valid: true,
		},
	})

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "session_error", "Failed to create session", "", err)
		return
	}

	accessToken, err := cfg.JWT.Generate(user.ID.String(), sessionID.String())

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "jwt_error", "Failed to create access token", "", err)
		return
	}

	cfg.Redis.Set(ctx, "sess:"+refreshHash, user.ID.String(), refreshTTL)

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   shouldUserSecureCookie(r),
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Expires:  refreshExpires,
	})

	RespondWithJSON(w, http.StatusOK, map[string]any{
		"access_token": accessToken,
		"user": map[string]string{
			"user_id":  user.ID.String(),
			"username": user.Username,
			"email":    user.Email,
		},
	})
}
