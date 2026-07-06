package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/AliKefall/Somnambulist/internal/services/observability"
	"github.com/gorilla/websocket"
)

var errSendBufferFull = errors.New("Websocket send buffer is full")

const (
	maxMessageSize = 16 * 1024
	pongWait       = 75 * time.Second
	pingPeriod     = 25 * time.Second
	writeWait      = 10 * time.Second

	sendBufferSize = 256
)

type Client struct {
	UserID   string
	Username string
	Conn     *websocket.Conn

	Send    chan []byte
	Hub     *Hub
	Metrics *observability.Metrics

	ctx    context.Context
	cancel context.CancelFunc
}

func NewClient(
	conn *websocket.Conn,
	hub *Hub,
	metrics *observability.Metrics,
	userID,
	username string,
) *Client {
	cctx, cancel := context.WithCancel(context.Background())
	return &Client{
		UserID:   userID,
		Username: username,
		Conn:     conn,
		Send:     make(chan []byte, sendBufferSize),
		Hub:      hub,
		Metrics:  metrics,
		ctx:      cctx,
		cancel:   cancel,
	}
}

func (c *Client) writeJSON(v any) error {
	payload, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return c.sendRaw(payload)
}

func (c *Client) sendRaw(payload []byte) error {
	select {
	case c.Send <- payload:
		return nil
	default:
		return errSendBufferFull
	}
}

func (c *Client) Close() {
	c.cancel()
	_ = c.Conn.Close()
}

func (c *Client) ReadPump() {
	defer func() {
		if c.Hub != nil {
			c.Hub.unregister <- c
		}
		c.Close()
		if c.Metrics != nil {
			c.Metrics.DecWSConnections()
		}
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))

	c.Conn.SetPongHandler(func(_ string) error {
		return c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		if c.ctx.Err() != nil {
			return
		}

		var msg Message
		if err := c.Conn.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseGoingAway) {
				log.Printf("WS read error user=%s err=%v", c.Username, err)
			}
			return
		}
		c.Metrics.ObserveWSMessage("in", msg.Type)
		c.Hub.HandleMessage(c, msg)
	}

}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.Close()
	}()

	for {
		select {
		case <-c.ctx.Done():
			return
		case firstMsg, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, nil)
				return
			}
			if err := c.writeBatch(firstMsg); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) sendError(code, message string) {
	msg, err := newMessage(TypeError, ErrorPayload{Code: code, Message: message})
	if err != nil {
		return
	}
	_ = c.writeJSON(msg)
}

func (c *Client) writeBatch(first []byte) error {
	w, err := c.Conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}

	if _, err := w.Write(first); err != nil {
		_ = w.Close()
		return err
	}

	n := len(c.Send)
batchDrain:
	for range n {
		select {
		case msg := <-c.Send:
			if _, err := w.Write([]byte{'\n'}); err != nil {
				_ = w.Close()
				return err
			}
			if _, err := w.Write(msg); err != nil {
				_ = w.Close()
				return err
			}
		default:
			break batchDrain
		}
	}

	if err := w.Close(); err != nil {
		return err
	}

	if c.Metrics != nil {
		c.Metrics.ObserveWSMessage("out", "batched")
	}

	return nil
}
