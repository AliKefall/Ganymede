package endpoints

import (
	"bufio"
	"context"
	"database/sql"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/AliKefall/Somnambulist/internal/auth"
	"github.com/AliKefall/Somnambulist/internal/database"
	"github.com/google/uuid"
)

// NOTE: When working with contexts always create you own type
// If we ever need to add another one with the same value it will be overwritten.
type contextKey string

const (
	UserIDKey    contextKey = "user_id"
	SessionIDKey contextKey = "session_id"
)

func (cfg *Config) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// TOKEN EXTRACTION

		token, err := auth.GetBearer(r.Header)
		if err != nil {

			// fallback for websocket clients
			token = r.URL.Query().Get("token")

			if token == "" {
				RespondWithError(
					w,
					http.StatusUnauthorized,
					"token_error",
					"missing token",
					"",
					nil,
				)
				return
			}
		}

		ctx := r.Context()

		// JWT VERIFY

		claims, err := cfg.JWT.Verify(token)
		if err != nil {
			RespondWithError(
				w,
				http.StatusUnauthorized,
				"token_error",
				"invalid token",
				"",
				nil,
			)
			return
		}

		// REDIS BLACKLIST CHECK

		blackKey := "bl:" + token

		exists, err := cfg.Redis.Exists(ctx, blackKey).Result()
		if err == nil && exists == 1 {
			RespondWithError(
				w,
				http.StatusUnauthorized,
				"token_error",
				"token revoked",
				"",
				nil,
			)
			return
		}

		// SESSION VALIDATIO
		sessionUUID, err := uuid.Parse(claims.SessionID)
		if err != nil {
			RespondWithError(
				w,
				http.StatusUnauthorized,
				"session_error",
				"invalid session id",
				"",
				err,
			)
			return
		}

		session, err := cfg.Queries.GetSessionByID(ctx, sessionUUID)
		if err != nil {
			RespondWithError(
				w,
				http.StatusUnauthorized,
				"session_error",
				"session not found",
				"",
				err,
			)
			return
		}

		now := time.Now().UTC()

		if session.RevokedAt.Valid {
			RespondWithError(
				w,
				http.StatusUnauthorized,
				"session_error",
				"session revoked",
				"",
				nil,
			)
			return
		}

		if now.After(session.ExpiresAt) {
			RespondWithError(
				w,
				http.StatusUnauthorized,
				"session_error",
				"session expired",
				"",
				nil,
			)
			return
		}

		if now.After(session.MaxExpiresAt) {
			RespondWithError(
				w,
				http.StatusUnauthorized,
				"session_error",
				"session lifetime exceeded",
				"",
				nil,
			)
			return
		}

		// ASYNC LAST USED UPDATE

		go func(sessionID uuid.UUID) {

			ctx, cancel := context.WithTimeout(
				context.Background(),
				3*time.Second,
			)
			defer cancel()

			_ = cfg.Queries.UpdateSessionLastUsed(
				ctx,
				database.UpdateSessionLastUsedParams{
					ID: sessionID,
					LastUsedAt: sql.NullTime{
						Valid: true,
						Time:  now,
					},
				},
			)

		}(session.ID)

		// CONTEXT INJECTION

		ctx = context.WithValue(ctx, UserIDKey, session.UserID)
		ctx = context.WithValue(ctx, SessionIDKey, session.ID)

		// RESPONSE WRAPPER

		rw := &responseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// Websocket will use this middleware that is the very reason
// why we need to write our own hijacker and response writer.
func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {

	hijacker, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New(
			"underlying ResponseWriter does not support hijacking",
		)
	}

	return hijacker.Hijack()
}

// streaming support
func (rw *responseWriter) Flush() {

	flusher, ok := rw.ResponseWriter.(http.Flusher)
	if !ok {
		return
	}

	flusher.Flush()
}

// HTTP/2 server push support
func (rw *responseWriter) Push(
	target string,
	opts *http.PushOptions,
) error {

	pusher, ok := rw.ResponseWriter.(http.Pusher)
	if !ok {
		return http.ErrNotSupported
	}

	return pusher.Push(target, opts)
}
