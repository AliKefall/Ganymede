package matchmaking

import (
	"time"

	"github.com/google/uuid"
)

type QueueEntry struct {
	UserID      uuid.UUID `json:"user_id"`
	Rating      int       `json:"rating"`
	TimeControl string    `json:"time_control"`

	JoinedAt  time.Time `json:"joined_at"`
	MinRating int       `json:"min_rating"`
	MaxRating int       `json:"max_rating"`
}

type Match struct {
	ID      uuid.UUID `json:"id"`
	WhiteID uuid.UUID `json:"white_id"`
	BlackID uuid.UUID `json:"black_id"`

	WhiteRating int       `json:"white_rating"`
	BlackRating int       `json:"black_rating"`
	TimeControl string    `json:"time_control"`
	CreatedAt   time.Time `json:"created_at"`
}


