package ratelimiter

import (
	"context"
	_ "embed"
	"time"

	"github.com/redis/go-redis/v9"
)

//go:embed lua/token_bucket.lua
var tokenBucketLua string

type RedisTokenBucketLimiter struct {
	client     *redis.Client
	capacity   float64
	refillRate float64
	ctx        context.Context
	luaScript  *redis.Script
}

func NewRedisTokenBucketLimiter(client *redis.Client, capacity int, refillPerSec float64) (*RedisTokenBucketLimiter, error) {
	script := redis.NewScript(tokenBucketLua)

	return &RedisTokenBucketLimiter{
		client:     client,
		capacity:   float64(capacity),
		refillRate: refillPerSec,
		ctx:        context.Background(),
		luaScript:  script,
	}, nil
}

func (l *RedisTokenBucketLimiter) Allow(key string, now time.Time) (bool, error) {
	redisKey := "rate_token" + key
	result, err := l.luaScript.Run(
		l.ctx,
		l.client,
		[]string{redisKey},
		l.capacity,
		now.Unix(),
		1,
	).Int()

	if err != nil {
		return false, err
	}
	return result == 1, nil
}
