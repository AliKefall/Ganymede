package websocket

import (
	"encoding/json"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/AliKefall/Somnambulist/internal/database"
	"github.com/AliKefall/Somnambulist/internal/services/observability"
	"github.com/google/uuid"
)

// Message type for signaling, expand this with causion.
const (
	MessageTypeDirectMessage = "direct_message"
	MessageTypeGameMove      = "game_move"
	MessageTypeMatchFound    = "match_found"
	MessageTypeMatchEnded    = "match_ended"
	MessageTypeError         = "error"
)

// You can write them down without these variables.
// But they are kinde useful.
var (
	ErrRecipientRequired = errors.New("Recipient is required")
	ErrRecipientOffline  = errors.New("Recipient is offline")
	ErrInvalidEnvelope   = errors.New("Invalid websocket envelope")
)

type Envelope struct {
	ID        string          `json:"id,omitempty"`
	Type      string          `json:"type"`
	Version   int             `json:"version"`
	Sender    *EnvelopeUser   `json:"sender,omitempty"`
	Recipient *EnvelopeUser   `json:"recipient,omitempty"`
	Payload   json.RawMessage `json:"payload,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
}

type EnvelopeUser struct {
	ID       string `json:"id"`
	Username string `json:"username,omitempty"`
}

type DirectMessagePayload struct {
	ConversationID string `json:"conversation_id,omitempty"`
	Content        string `json:"content"`
}



type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type inboundEnvelope struct {
	sender *Client
	env    Envelope
}

type Hub struct {
	Users      map[string]map[*Client]bool
	mu         sync.RWMutex
	Clients    map[*Client]bool
	Queries    *database.Queries
	Register   chan *Client
	Unregister chan *Client
	Inbound    chan inboundEnvelope
	Metrics    *observability.Metrics
	active     int64 // I will simplify this in production
}

func NewHub(queries *database.Queries, metrics *observability.Metrics) *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Users:      make(map[string]map[*Client]bool),
		Queries:    queries,
		Register:   make(chan *Client, 256),
		Unregister: make(chan *Client, 256),
		Inbound:    make(chan inboundEnvelope, 1024),
		Metrics:    metrics,
	}
}

func NewEnvelope(messageType string, sender, recipient *EnvelopeUser, payload any) (Envelope, error) {
	raw, err := json.Marshal(payload)
	if err != nil {
		return Envelope{}, err
	}

	return Envelope{
		ID:        uuid.NewString(),
		Type:      messageType,
		Version:   1,
		Sender:    sender,
		Recipient: recipient,
		Payload:   raw,
		CreatedAt: time.Now().UTC(),
	}, nil
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.registerClient(client)
		case client := <-h.Unregister:
			h.unregisterClient(client)
		case inbound := <-h.Inbound:
			h.dispatch(inbound.sender, inbound.env)
		}
	}
}

func (h *Hub) HandleMessage(sender *Client, env Envelope) {
	if h == nil || sender == nil {
		return
	}
	if env.Version == 0 {
		env.Version = 1
	}
	if env.CreatedAt.IsZero() {
		env.CreatedAt = time.Now().UTC()
	}
	if env.ID == "" {
		env.ID = uuid.NewString()
	}
	select {
	case h.Inbound <- inboundEnvelope{sender: sender, env: env}:
	default:
		sender.SendError("hub_backpressure", "Websocket hub is busy")
	}
}

func (h *Hub) dispatch(sender *Client, env Envelope) {
	senderUser := &EnvelopeUser{ID: sender.UserID, Username: sender.Username}
	env.Sender = senderUser

	switch env.Type {
	case MessageTypeDirectMessage:
		if env.Recipient == nil || env.Recipient.ID == "" {
			sender.SendError("recipient_required", ErrRecipientRequired.Error())
			return
		}
		if err := h.SendDirectEnvelope(env.Recipient.ID, env); err != nil {
			sender.SendError("recipient_offline", ErrRecipientOffline.Error())
			return
		}
		_ = sender.writeJSON(env)
	case MessageTypeGameMove:
		if env.Recipient == nil || env.Recipient.ID == "" {
			sender.SendError("recipient_required", ErrRecipientRequired.Error())
			return
		}
		if err := h.SendDirectEnvelope(env.Recipient.ID, env); err != nil {
			sender.SendError("recipient_offline", err.Error())
		}
	default:
		sender.SendError("unknown_type", "Unsupported websocket message type")

	}
}

func (h *Hub) SendDirectMessage(senderID, senderUsername, receiverID, content string) error {
	payload := DirectMessagePayload{Content: content}
	env, err := NewEnvelope(
		MessageTypeDirectMessage,
		&EnvelopeUser{ID: senderID, Username: senderUsername},
		&EnvelopeUser{ID: receiverID},
		payload,
	)
	if err != nil {
		return err
	}
	return h.SendDirectEnvelope(receiverID, env)
}

func (h *Hub) SendDirectEnvelope(receiverID string, env Envelope) error {
	if h == nil || receiverID == "" {
		return ErrRecipientRequired
	}
	h.mu.Lock()
	clients := h.Users[receiverID]
	if len(clients) == 0 {
		h.mu.RUnlock()
		return ErrRecipientRequired
	}
	localClients := make([]*Client, 0, len(clients))
	for client := range clients {
		localClients = append(localClients, client)
	}
	h.mu.RUnlock()

	payload, err := json.Marshal(env)
	if err != nil {
		return err
	}
	for _, client := range localClients {
		if err := client.sendRaw(payload); err != nil {
			client.Close()
		}
	}
	h.Metrics.ObserveWSMessage("out", env.Type)

	return nil
}

func (h *Hub) ActiveConnections() int64 {
	return atomic.LoadInt64(&h.active)
}

func (h *Hub) registerClient(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.Clients[c] = true
	if h.Users[c.UserID] == nil {
		h.Users[c.UserID] = make(map[*Client]bool)
	}
	h.Users[c.UserID][c] = true
	atomic.AddInt64(&h.active, 1)
}

func (h *Hub) unregisterClient(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.Clients[c]; !ok {
		return
	}
	delete(h.Clients, c)
	if userConns := h.Users[c.UserID]; userConns != nil {
		delete(userConns, c)
		if len(userConns) == 0 {
			delete(h.Users, c.UserID)
		}
	}

	atomic.AddInt64(&h.active, -1)

}
