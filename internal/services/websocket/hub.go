package websocket

import (
	"sync/atomic"

	"github.com/AliKefall/Somnambulist/internal/database"
	"github.com/AliKefall/Somnambulist/internal/services/observability"
)

type Hub struct {
	CLients    map[*Client]bool
	Users      map[string]map[*Client]bool
	Queries    *database.Queries
	Register   chan *Client
	Unregister chan *Client
	Metrics    observability.Metrics
	active     int64
}

type Message struct {
	Type           string `json:"type"`
	To             string `json:"to"`
	SenderID       string `json:"sender_id"`
	SenderUsername string `json:"sender_username"`
	Content        string `json:"content"`
	CreatedAt      int64  `json:"created_at"`
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.registerClient(client)
		case client := <-h.Unregister:
			h.unregisterClient(client)
		case env := h.Inbound:
			h.dispatch(env.sender, env.msg)
		}
	}
}

func (h *Hub) ActiveConnections() int64 {
	return atomic.LoadInt64(&h.active)
}

func (h *Hub) registerClient(c *Client) {
	h.CLients[c] = true
	if h.Users[c.Username] == nil {
		h.Users[c.Username] = make(map[*Client]bool)
	}
	h.Users[c.Username][c] = true
	atomic.AddInt64(&h.active, 1)
}

func (h *Hub) unregisterClient(c *Client) {
	delete(h.CLients, c)
	if userConns := h.Users[c.Username]; userConns != nil {
		delete(userConns, c)
		if len(userConns) == 0 {
			delete(h.Users, c.Username)
		}

	}
	atomic.AddInt64(&h.active, -1)
}

func (h *Hub) SendDirectMessage(senderID string, receiverID string) {

}
