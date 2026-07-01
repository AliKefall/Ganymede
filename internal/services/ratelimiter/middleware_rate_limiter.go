package ratelimiter

import (
	"context"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

// RateLimiter is implemented by any backend capable of answering "is this
// key allowed right now". ctx should be used to bound how long the
// underlying call is allowed to take.
type RateLimiter interface {
	Allow(ctx context.Context, key string, now time.Time) (bool, error)
}

func ClientIP(r *http.Request) string {
	if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		for candidate := range strings.SplitSeq(forwardedFor, ",") {
			ip := strings.TrimSpace(candidate)
			if parsed := net.ParseIP(ip); parsed != nil {
				return parsed.String()
			}
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
}


const rateLimitTimeout = 200 * time.Millisecond

func MiddlewareRateLimiter(limiter RateLimiter) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := ClientIP(r) + ":" + r.URL.Path

			ctx, cancel := context.WithTimeout(r.Context(), rateLimitTimeout)
			defer cancel()

			ok, err := limiter.Allow(ctx, key, time.Now())
			if err != nil {
				log.Printf("rate limiter backend error for key=%q: %v", key, err)
				h.ServeHTTP(w, r)
				return
			}
			if !ok {
				w.Header().Set("Retry-After", "1")
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}

