package games

type EnqueuePlayer struct {
	UserID       string `json:"user_id"`
	Username     string `json:"username"`
	Rating       string `json:"rating"`
	Joined_at    string `json:"joined_at"`
	Time_control string `json:"time_control"`
}

type DequeuePlayer struct {
	UserID string `json:"user_id"`
}

type MatchFound struct {
}

type GameFoundPayload struct {
	GameID      string
	WhiteID     string
	BlackID     string
	TimeControl string
}
