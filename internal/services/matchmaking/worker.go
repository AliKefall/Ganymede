package matchmaking

import (
	"context"
	"errors"
	"log"
	"math/rand/v2"
	"time"

	"github.com/AliKefall/Somnambulist/internal/services/websocket"
	"github.com/google/uuid"
)

type MatchHandler func(context.Context, Match) error

type Worker struct {
	Service      *Service
	Hub          *websocket.Hub
	TimeControls []string
	PollInterval time.Duration
	BatchSize    int
	OnMatch      MatchHandler
}

func (w *Worker) Run(ctx context.Context) error {
	if w == nil || w.Service == nil {
		return errors.New("Matchmaking worker service is nil")
	}
	interval := w.PollInterval
	if interval <= 0 {
		interval = DefaultPollInterval
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			for _, timeControl := range w.timeControls() {
				matches, err := w.Service.FindMatches(ctx, timeControl, w.batchSize())
				if err != nil {
					log.Printf("Matchmaking tick failed time_control=%s err=%v", timeControl, err)
					continue
				}
				for _, match := range matches {
					if w.OnMatch != nil {
						if err := w.OnMatch(ctx, match); err != nil {
							log.Printf("Matchmaking match handler failed match_id=%s err=%v", match.ID, err)
						}
					}
					//Notify matches
				}
			}
		}
	}
}

func (s *Service) FindMatches(
	ctx context.Context,
	timeControl string,
	limit int,
) ([]Match, error) {

	if s == nil || s.Redis == nil {
		return nil, errors.New("redis nil")
	}

	script, err := luaScripts.ReadFile("lua/matchmake.lua")
	if err != nil {
		return nil, err
	}

	result, err := s.Redis.Eval(
		ctx,
		string(script),
		[]string{queueKey(timeControl)},
		time.Now().UTC().Unix(),
		s.RatingWindow,
		s.WindowGrowth,
		s.MaxWindow,
		limit,
	).Result()

	if err != nil {
		return nil, err
	}

	rows, ok := result.([]any)
	if !ok || len(rows) == 0 {
		return nil, nil
	}

	matches := make([]Match, 0, len(rows))

	for _, row := range rows {
		fields, ok := row.([]any)
		if !ok || len(fields) < 7 {
			continue
		}

		matchRaw := parseLuaString(fields[0])

		p1ID, err := uuid.Parse(parseLuaString(fields[1]))
		if err != nil {
			continue
		}

		p2ID, err := uuid.Parse(parseLuaString(fields[2]))
		if err != nil {
			continue
		}

		p1Rating, err := parseLuaInt(fields[3])
		if err != nil {
			continue
		}

		p2Rating, err := parseLuaInt(fields[4])
		if err != nil {
			continue
		}

		p1Username := parseLuaString(fields[5])
		p2Username := parseLuaString(fields[6])

		matchID, err := uuid.Parse(matchRaw)
		if err != nil {
			matchID = uuid.NewSHA1(
				uuid.NameSpaceURL,
				[]byte(matchRaw),
			)
		}

		whiteID := p1ID
		blackID := p2ID

		whiteRating := p1Rating
		blackRating := p2Rating

		whiteUsername := p1Username
		blackUsername := p2Username

		// randomize colors
		if rand.IntN(2) == 1 {
			whiteID, blackID = blackID, whiteID
			whiteRating, blackRating = blackRating, whiteRating
			whiteUsername, blackUsername = blackUsername, whiteUsername
		}

		matches = append(matches, Match{
			ID: matchID,

			WhiteID: whiteID,
			BlackID: blackID,

			WhiteUsername: whiteUsername,
			BlackUsername: blackUsername,

			WhiteRating: whiteRating,
			BlackRating: blackRating,

			TimeControl: timeControl,
			CreatedAt:   time.Now().UTC(),
		})
	}

	return matches, nil
}

func (w *Worker) timeControls() []string {
	if len(w.TimeControls) == 0 {
		return []string{"rapid"}
	}
	return w.TimeControls
}

func (w *Worker) batchSize() int {
	if w.BatchSize <= 0 {
		return DefaultMatchBatchSize
	}
	return w.BatchSize
}

