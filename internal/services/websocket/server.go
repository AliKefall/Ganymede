package websocket

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/AliKefall/Somnambulist/internal/services/observability"
	"github.com/gorilla/websocket"
)

var (
	allowedOrigins = make(map[string]struct{})
)

func init() {
	for origin := range strings.SplitSeq(os.Getenv("WS_ALLOWED_ORIGINS"), ",") {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			allowedOrigins[origin] = struct{}{}
			// A little information in golang struct{}{} is actually 0 byte
			// That is why its been used in here.
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")

		if os.Getenv("APP_ENV") == "development" {
			return true
		}

		_, ok := allowedOrigins[origin]
		return ok
	},
}

func ServeWS(
	hub *Hub,
	metrics *observability.Metrics,
	w http.ResponseWriter,
	r *http.Request,
	userID string,
	username string,
	onConnected func(*Client),
) {
	if hub != nil && hub.ActiveConnection() > 20000 {
		http.Error(w, "Server is busy", http.StatusServiceUnavailable)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if metrics != nil {
			metrics.ObserveWSConnection("upgrade_failed")
		}
		log.Printf("ws upgrade failed for user=%s remote=%s err=%v", username, r.RemoteAddr, err)
		return
	}

	if metrics != nil {
		metrics.ObserveWSConnection("connected")
	}

	ctx := r.Context()
	client := NewClient(ctx, conn, hub, metrics, userID, username)

	select {
	case hub.register <- client:
	case <-time.After(3 * time.Second):
		log.Printf("ws register timeout user=%s", username)
		conn.Close()
		return
	case <-r.Context().Done():
		conn.Close()
		return

	}

	if onConnected != nil {
		onConnected(client)
	}

	go func() {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic in WritePump user=%s err=%v", username, rec)
			}
		}()
		client.WritePump()
	}()

	go func() {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic in ReadPump user=%s err=%v", username, rec)
			}
		}()
		client.ReadPump()
	}()

}
