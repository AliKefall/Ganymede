package main

import (
	"log"
	"net/http"
	"net/url"

	"github.com/AliKefall/Somnambulist/internal/endpoints"
	"github.com/AliKefall/Somnambulist/internal/services/observability"
	"github.com/AliKefall/Somnambulist/internal/services/ratelimiter"
	"github.com/AliKefall/Somnambulist/internal/services/websocket"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
)

func buildRouter(config *ServerConfig, deps serverDependencies) chi.Router {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: config.AllowedOrigins,
		AllowOriginFunc: func(r *http.Request, origin string) bool {
			u, err := url.Parse(origin)
			if err != nil {
				return false
			}
			if (u.Hostname() == "localhost" || u.Hostname() == "127.0.0.1") && (u.Scheme == "http" || u.Scheme == "https") {
				return true
			}
			for _, allowed := range config.AllowedOrigins {
				if allowed == origin {
					return true
				}
			}
			return false
		},
		AllowedMethods:   []string{"GET", "OPTIONS", "POST", "DELETE", "PUT", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"X-Request-Id"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(securityHeaderMiddleware)
	r.Use(observability.MiddlewareHTTPMetrics(deps.metrics))

	rateLimiter, err := ratelimiter.NewRedisTokenBucketLimiter(deps.redis, 20, 10)
	if err != nil {
		log.Fatalf("Ratelimiter init failed: %v", err)
	}
	r.Use(ratelimiter.MiddlewareRateLimiter(rateLimiter))

	r.With(deps.config.AuthMiddleware).Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		uid, ok := r.Context().Value(endpoints.UserIDKey).(uuid.UUID)
		if !ok || uid == uuid.Nil {
			endpoints.RespondWithError(w, http.StatusUnauthorized, "websocket_error", "Invalid user context", "", nil)
			return
		}

		user, err := deps.queries.GetUserByID(r.Context(), uid)
		if err != nil {
			endpoints.RespondWithError(w, http.StatusUnauthorized, "database_error", "User not found", "", err)
			return
		}

		websocket.ServeWS(deps.hub, deps.metrics, w, r, uid.String(), user.Username, nil)
	})

	r.Group(func(pr chi.Router) {
		pr.Use(deps.config.AuthMiddleware)
		pr.Get("/friends", deps.config.HandleListFriends)
		pr.Get("/friends/requests", deps.config.HandleFriendRequests)
		pr.Post("/friends/requests", deps.config.HandleSendFriendRequest)
		pr.Post("/friends/requests/accept", deps.config.HandleAcceptFriendRequest)
		pr.Post("/friends/requests/reject", deps.config.HandleRejectFriendRequest)
	})

	r.Route("/auth", func(ar chi.Router) {
		ar.Post("/register", deps.config.HandlerRegister)
		ar.Post("/login", deps.config.HandlerLogin)
		ar.Post("/refresh", deps.config.HandlerRefresh)
		ar.Post("/logout", deps.config.HandlerLogout)
	})
	return r
}

func securityHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers := w.Header()
		headers.Set("X-Content-Type-Options", "nosniff")
		headers.Set("X-Frame-Options", "DENY")
		headers.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		headers.Set("Cross-Origin-Resource-Policy", "same-site")
		headers.Set("Cross-Origin-Opener-Policy", "same-origin")
		headers.Set("Cross-Origin-Embedder-Policy", "credentialless")
		if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
			headers.Set("Strict-Transport-Security", "max-age=31536000; includeSubdomains")
		}
		next.ServeHTTP(w, r)
	})
}
