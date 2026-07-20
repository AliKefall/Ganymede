package endpoints

import (
	"database/sql"

	"github.com/AliKefall/Somnambulist/internal/auth"
	"github.com/AliKefall/Somnambulist/internal/database"
	"github.com/AliKefall/Somnambulist/internal/services/chat"
	"github.com/AliKefall/Somnambulist/internal/services/matchmaking"
	ratelimite "github.com/AliKefall/Somnambulist/internal/services/ratelimiter"
	"github.com/AliKefall/Somnambulist/internal/services/websocket"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	WS                 *websocket.Hub
	DB                 *sql.DB
	Queries            *database.Queries
	JWT                *auth.JWTManager
	Hasher             *auth.PasswordHasher
	Redis              *redis.Client
	RateLimiter        ratelimite.RateLimiter
	MatchmakingService matchmaking.Service
	Chat               *chat.Service
}

func NewConfig(
    ws *websocket.Hub,
    dbconn *sql.DB,
    queries *database.Queries,
    chatService *chat.Service,
    jwt *auth.JWTManager,
    hasher *auth.PasswordHasher,
    redisClient *redis.Client,
) *Config {
    return &Config{
        WS:      ws,
        DB:      dbconn,
        Queries: queries,
        JWT:     jwt,
        Hasher:  hasher,
        Redis:   redisClient,
        Chat:    chatService,
    }
}
