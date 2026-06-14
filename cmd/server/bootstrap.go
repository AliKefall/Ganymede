package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/AliKefall/Somnambulist/internal/auth"
	"github.com/AliKefall/Somnambulist/internal/database"
	"github.com/AliKefall/Somnambulist/internal/endpoints"
	"github.com/AliKefall/Somnambulist/internal/services/observability"
	"github.com/AliKefall/Somnambulist/internal/services/websocket"
	"github.com/redis/go-redis/v9"
)

type serverDependencies struct {
	config  *endpoints.Config
	queries *database.Queries
	hub     *websocket.Hub
	metrics *observability.Metrics
	redis   *redis.Client
}

func bootstrapServer(config *ServerConfig) (*sql.DB, serverDependencies) {
	conn := mustOpenDatabase(config)
	queries := database.New(conn)

	metrics := observability.New()
	hub := websocket.NewHub(queries, metrics)
	redisClient := NewRedisClient(config.RedisURL)

	handler := endpoints.NewConfig(conn, queries, auth.NewJWTManager(config.JWTSecret, 15*time.Minute), auth.NewPasswordHasher(), redisClient)
	return conn, serverDependencies{
		config:  handler,
		queries: queries,
		hub:     hub,
		metrics: metrics,
		redis: redisClient,
	}
}

func mustOpenDatabase(config *ServerConfig) *sql.DB {
	dbURL := config.DBUrl
	conn, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatal("database connection error: ", err)
	}

	conn.SetMaxOpenConns(50)
	conn.SetMaxIdleConns(25)
	conn.SetConnMaxLifetime(5 * time.Minute)
	if err := conn.Ping(); err != nil {
		log.Fatal("database ping failed", err)
	}
	return conn
}

func NewRedisClient(redisURL string) *redis.Client {
	if redisURL != "" {
		opt, err := redis.ParseURL(redisURL)
		if err != nil {
			return redis.NewClient(opt)
		} else {
			log.Fatalf("Error creating redis client: %v", err)
		}
	}
	return redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		Password:     "",
		DB:           0,
		PoolSize:     10,
		MaxIdleConns: 5,
	})
}
