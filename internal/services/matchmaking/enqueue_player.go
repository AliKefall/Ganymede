package matchmaking

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	DefaultQueueTTL       = 5 * time.Minute
	DefaultPollInterval   = 75
	DefaultRatingWindow   = 25
	DefaultWindowGrowth   = 25
	DefaultMaxWindow      = 500
	DefaultMatchBatchSize = 25
	DefaultRating = 1500
)

var (
	//go:embed lua/*.lua
	luaScripts embed.FS

	ErrAlreadyQueued = errors.New("player is already queued")
	ErrNotQueued     = errors.New("player is not queued")
)

type Service struct {
	Redis        redis.Cmdable
	QueueTTL     time.Duration
	RatingWindow int
	WindowGrowth int
	MaxWindow    int
}

func NewService(redisClient redis.Cmdable) *Service {
	return &Service{
		Redis:        redisClient,
		QueueTTL:     DefaultQueueTTL,
		RatingWindow: DefaultRatingWindow,
		WindowGrowth: DefaultWindowGrowth,
		MaxWindow:    DefaultMaxWindow,
	}
}

func (s *Service) EnqueuePlayer(ctx context.Context, entry QueueEntry) error {
	if s == nil || s.Redis == nil {
		return errors.New("Matchmaking redis client is nil")
	}

	if entry.UserID == uuid.Nil {
		return errors.New("user_id is required")
	}

	if entry.Rating <= 0 {
		entry.Rating = int(DefaultRating)
	}

	if entry.JoinedAt.IsZero() {
		entry.JoinedAt = time.Now().UTC()
	}

	if entry.TimeControl == "" {
		entry.TimeControl = "rapid"
	}

	script, err := luaScripts.ReadFile("lua/enqueue.lua")
	if err != nil {
		return err
	}
	result, err := s.Redis.Eval(ctx, string(script), []string{
		queueKey(entry.TimeControl),
		userKey(entry.UserID),
	}, entry.UserID.String(), entry.Username, entry.Rating, entry.JoinedAt.Unix(), entry.TimeControl, int(ttl(s.QueueTTL).Seconds())).Result()
	if err != nil {
		if err.Error() == "already_queued" || err.Error() == "ERR already_queued" {
			return ErrAlreadyQueued
		}
		return err
	}
	if fmt.Sprint(result) != "ok" {
		return fmt.Errorf("unexpected enqueue result: %v", result)
	}
	return nil

}

func ttl(value time.Duration) time.Duration {
	if value <= 0 {
		return DefaultQueueTTL
	}
	return value
}

func queueKey(timeControl string) string {
	if timeControl == "" {
		timeControl = "rapid"
	}
	return "matchmaking:queue:" + timeControl
}

func userKey(userID uuid.UUID) string {
	return "matchmaking:user:" + userID.String()
}

// This is kinda bad code. But it gets the job done for now
// NOTE: Change this in production
func parseLuaString(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	default:
		return fmt.Sprint(v)
	}
}

func parseLuaInt(value any) (int, error) {
	switch v := value.(type) {
	case int:
		return v, nil
	case []byte:
		return strconv.Atoi(string(v))
	case string:
		return strconv.Atoi(v)
	case int64:
		return int(v), nil
	default:
		return 0, fmt.Errorf("unsupported integer value %T", value)
	}

}
