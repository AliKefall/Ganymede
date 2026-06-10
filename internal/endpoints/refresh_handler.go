package endpoints

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"github.com/AliKefall/Somnambulist/internal/auth"
	"github.com/AliKefall/Somnambulist/internal/database"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func (cfg *Config) HandlerRefresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	now := time.Now()

	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "refresh_error", "Missing refresh token", "", nil)
		return
	}

	hash := sha256.Sum256([]byte(cookie.Value))
	refreshHash := hex.EncodeToString(hash[:])

	blacklisted, err := cfg.Redis.Get(ctx, "bl:"+refreshHash).Result()
	if err != nil || errors.Is(err, redis.Nil) {
		RespondWithError(w, http.StatusInternalServerError, "redis_error", "Redis read failed", "", nil)
		return
	}
	if blacklisted == "1" {
		clearRefreshToken(w, r)
		RespondWithError(w, http.StatusUnauthorized, "token_reuse", "Token reuse detected", "", nil)
		return
	}

	tx, err := cfg.DB.BeginTx(ctx, nil)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "transaction_error", "Transaction failed on refresh endpoint", "", err)
		return
	}

	defer tx.Rollback()

	qtx := cfg.Queries.WithTx(tx)

	session, err := qtx.GetSessionForUpdateByTokenHash(ctx, refreshHash)

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "refresh_error", "Invalid refresh token", "", nil)
		return
	}

	if session.RevokedAt.Valid {
		_ = qtx.DeleteSessionsByUserID(ctx, session.UserID)
		clearRefreshToken(w, r)
		_ = tx.Commit()

		RespondWithError(w, http.StatusUnauthorized, "session_expires", "Session expired", "", nil)
		return
	}
	if now.After(session.MaxExpiresAt) {
		_ = qtx.DeleteSessionsByUserID(ctx, session.UserID)
		clearRefreshToken(w, r)
		_ = tx.Commit()

		RespondWithError(w, http.StatusUnauthorized, "session_exceed", "Session lifetime exceeded", "", nil)
		return
	}

	_ = qtx.RevokeSession(ctx, database.RevokeSessionParams{
		ID:        session.ID,
		RevokedAt: sql.NullTime{Time: now, Valid: true},
	})

	newRefresh, err := auth.MakeRefreshToken()

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "token_error", "Token generation failed", "", err)
		return
	}

	newHash := sha256.Sum256([]byte(newRefresh))
	newRefreshHash := hex.EncodeToString(newHash[:])

	newSessionID := uuid.New()
	refreshTTL := 7 * 24 * time.Hour

	newExpires := now.Add(refreshTTL)

	_, err = qtx.CreateSession(ctx, database.CreateSessionParams{
		ID:               newSessionID,
		UserID:           session.UserID,
		RefreshTokenHash: newRefreshHash,

		UserAgent: session.UserAgent,
		IpAddress: session.IpAddress,
		CreatedAt: now,
		ExpiresAt: session.MaxExpiresAt,
		LastUsedAt: sql.NullTime{
			Time:  now,
			Valid: true,
		},
	})

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "session_error", "Failed to create session", "", err)
		return
	}

	accessToken, err := cfg.JWT.Generate(session.UserID.String(), newSessionID.String())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "jwt_error", "JWT Generation failed", "", err)
		return
	}

	if err := tx.Commit(); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "transaction_error", "Commit failed", "", err)
		return
	}

	ttl := time.Until(newExpires)

	cfg.Redis.Set(ctx, "bl:"+refreshHash, "1", ttl)
	cfg.Redis.Del(ctx, "sess:"+refreshHash)
	cfg.Redis.Set(ctx, "sess:"+newRefreshHash, session.UserID.String(), ttl)

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefresh,
		HttpOnly: true,
		Secure:   shouldUserSecureCookie(r),
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Expires:  newExpires,
	})

	RespondWithJSON(w, http.StatusOK, map[string]string{
		"access_token": accessToken,
	})
}

func clearRefreshToken(w http.ResponseWriter, r *http.Request) {
	secureCookie := shouldUserSecureCookie(r)
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   secureCookie,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}
