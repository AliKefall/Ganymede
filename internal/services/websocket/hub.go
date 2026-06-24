package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/AliKefall/Somnambulist/internal/database"
	"github.com/AliKefall/Somnambulist/internal/services/observability"
	"github.com/google/uuid"
)

const (
	//Inbound
	TypeSendMessage = "send_message"

	//Outbound
	TypeNewMessage    = "new_message"
	TypeError         = "error"

)

var (
	ErrRecipientRequired = errors.New("recipient is required")
	ErrRecipientOffline  = errors.New("recipient is offline")
)

type Message struct {
	ID        string          `json:"id,omitempty"`
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	CreatedAt time.Time       `json:"created_at"`
}

type SendMessagePayload struct {
	RecipientID    string `json:"recipient_id"`
	ConversationID string `json:"conversation_id,omitempty"`
	Content        string `json:"content"`
}

type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type inbound struct {
	sender *Client
	msg    Message
}

type NewMessagePayload struct {
	SenderID       string `json:"sender_id"`
	SenderUsername string `json:"sender_username"`
	ConversationID string `json:"conversation_id,omitempty"`
	Content        string `json:"content"`
}

// atomic alignment for 32 bit system. Active always must come first
type Hub struct {
	active             int64
	mu                 sync.RWMutex
	clients            map[*Client]bool
	users              map[string]map[*Client]bool
	queries            *database.Queries
	metrics            *observability.Metrics
	register           chan *Client
	unregister         chan *Client

	inbound            chan inbound
}

func NewHub(queries *database.Queries, metrics *observability.Metrics) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		users:      make(map[string]map[*Client]bool),
		queries:    queries,
		metrics:    metrics,
		register:   make(chan *Client, 256),
		unregister: make(chan *Client, 256),
		inbound:    make(chan inbound, 1024),
	}
}

func (h *Hub) Register(c *Client) {
	h.register <- c
}

func (h *Hub) Unregister(c *Client) {
	h.unregister <- c
}

func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case c := <-h.register:
			h.registerClient(c)
		case c := <-h.unregister:
			h.unregisterClient(c)
		case in := <-h.inbound:
			h.dispatch(in.sender, in.msg)
		}
	}
}

func (h *Hub) HandleMessage(sender *Client, msg Message) {
	if h == nil || sender == nil {
		return
	}
	if msg.ID == "" {
		msg.ID = uuid.NewString()
	}
	if msg.CreatedAt.IsZero() {
		msg.CreatedAt = time.Now().UTC()
	}
	select {
	case h.inbound <- inbound{sender: sender, msg: msg}:
	default:
		sender.sendError("hub_backpressure", "websocket hub is busy")
	}
}

func (h *Hub) dispatch(sender *Client, msg Message) {
	switch msg.Type {
	case TypeSendMessage:
		h.handleSendMessage(sender, msg)
	default:
		sender.sendError("unknown_type", "unsupported message type: "+msg.Type)
	}
}

func (h *Hub) handleSendMessage(sender *Client, msg Message) {
	var p SendMessagePayload
	if err := json.Unmarshal(msg.Payload, &p); err != nil {
		sender.sendError("invalid_payload", "send_message payload could not be parsed")
		return
	}
	if p.Content == "" {
		sender.sendError("empty_content", "message content cannot be empty")
		return
	}
	outPayload := NewMessagePayload{
		SenderID:       sender.UserID,
		SenderUsername: sender.Username,
		ConversationID: p.ConversationID,
		Content:        p.Content,
	}
	out, err := newMessage(TypeNewMessage, outPayload)
	if err != nil {
		sender.sendError("internal_error", "Could not build message")
		return
	}

	if err := h.deliver(p.RecipientID, out); err != nil {
		sender.sendError("recipient_offline", ErrRecipientOffline.Error())
		return
	}

	// Sender echo so they can see their own message too.
	_ = sender.writeJSON(out)

}

func newMessage(msgType string, payload any) (Message, error) {
	raw, err := json.Marshal(payload)
	if err != nil {
		return Message{}, err
	}

	return Message{
		ID:        uuid.NewString(),
		Type:      msgType,
		Payload:   raw,
		CreatedAt: time.Now().UTC(),
	}, nil
}

func (h *Hub) deliver(userID string, msg Message) error {
	if userID == "" {
		return ErrRecipientRequired
	}

	h.mu.RLock()
	conns := h.users[userID]
	if len(conns) == 0 {
		h.mu.RUnlock()
		return ErrRecipientOffline
	}

	targets := make([]*Client, 0, len(conns))
	for c := range conns {
		targets = append(targets, c)
	}
	h.mu.RUnlock()

	raw, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	for _, c := range targets {
		if err := c.sendRaw(raw); err != nil {
			c.Close()
		}
	}

	if h.metrics != nil {
		h.metrics.ObserveWSMessage("out", msg.Type)
	}
	return nil
}

// ----------------------------------------------

func (h *Hub) ActiveConnection() int64 {
	return atomic.LoadInt64(&h.active)
}

func (h *Hub) registerClient(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[c] = true
	if h.users[c.UserID] == nil {
		h.users[c.UserID] = make(map[*Client]bool)
	}

	h.users[c.UserID][c] = true
	atomic.AddInt64(&h.active, 1)
}

func (h *Hub) unregisterClient(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[c]; !ok {
		return
	}
	delete(h.clients, c)

	if conns := h.users[c.UserID]; conns != nil {
		delete(conns, c)
		if len(conns) == 0 {
			delete(h.users, c.UserID)
		}
	}

	atomic.AddInt64(&h.active, -1)
}
