package endpoints

import (
	"time"

	"github.com/AliKefall/Somnambulist/internal/database"
)

type ChatMessageResponse struct {
	ID             string     `json:"id"`
	ConversationID string     `json:"conversation_id"`
	SenderID       string     `json:"sender_id"`
	Content        string     `json:"content"`
	CreatedAt      time.Time  `json:"created_at"`
	EditedAt       *time.Time `json:"edited_at,omitempty"`
}

type ConversationResponse struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

func toChatMessageResponse(m database.Message) ChatMessageResponse {
	return ChatMessageResponse{
		ID:             m.ID.String(),
		ConversationID: m.ConversationID.String(),
		SenderID:       m.SenderID.String(),
		Content:        m.Content,
		CreatedAt:      m.CreatedAt,
		EditedAt:       &m.EditedAt.Time,
	}
}

func toConversationResponse(c database.Conversation) ConversationResponse{
	return ConversationResponse{
		ID: c.ID.String(),
		Type: string(c.Type),
		CreatedAt: c.CreatedAt,
	}
}
