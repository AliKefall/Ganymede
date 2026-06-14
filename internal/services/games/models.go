package games

import (
	"github.com/AliKefall/Somnambulist/internal/services/websocket"
)

type GameHub struct {
	hub       *websocket.Hub
	MovesMade []string

}

type GameMovePayload struct {
	MatchID     string `json:"match_id"`
	MoveNumber  int    `json:"move_number"`
	PlayerID    string `json:"player_id"`
	SAN         string `json:"san"`
	UCI         string `json:"uci"`
	FENAfter    string `json:"fen_after"`
	WhiteTimeMS int64  `json:"white_time_ms"`
	BlackTimeMS int64  `json:"black_time_ms"`
}

type MatchFoundPayload struct {
	MatchID     string `json:"match_id"`
	WhiteID     string `json:"white_id"`
	BlackID     string `json:"black_id"`
	TimeControl string `json:"time_control"`
}
