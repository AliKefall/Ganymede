package endpoints

import (
	"database/sql"

	"github.com/AliKefall/Somnambulist/internal/auth"
	"github.com/AliKefall/Somnambulist/internal/database"
	"github.com/AliKefall/Somnambulist/internal/services/matchmaking"
	ratelimite "github.com/AliKefall/Somnambulist/internal/services/ratelimiter"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	DB      *sql.DB
	Queries *database.Queries
	JWT     *auth.JWTManager
	Hasher  *auth.PasswordHasher
	Redis   *redis.Client
	RateLimiter ratelimite.RateLimiter
	MatchmakingService matchmaking.Service
}

func NewConfig(dbconn *sql.DB, queries *database.Queries, jwt *auth.JWTManager, hasher *auth.PasswordHasher, redisClient *redis.Client) *Config {
	return &Config{
		DB:      dbconn,
		Queries: queries,
		JWT:     jwt,
		Hasher:  hasher,
		Redis:   redisClient,
	}
}
