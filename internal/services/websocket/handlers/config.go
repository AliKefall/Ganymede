package handlers

import (
	"github.com/AliKefall/Somnambulist/internal/database"
	"github.com/AliKefall/Somnambulist/internal/services/websocket"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Hub         *websocket.Hub
	Queries     *database.Queries
	RedisClient *redis.Client
}
