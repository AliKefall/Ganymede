package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/AliKefall/Somnambulist/internal/auth"
	"github.com/AliKefall/Somnambulist/internal/database"
	"github.com/AliKefall/Somnambulist/internal/endpoints"
	"github.com/AliKefall/Somnambulist/internal/services/chat"
	"github.com/AliKefall/Somnambulist/internal/services/observability"
	"github.com/AliKefall/Somnambulist/internal/services/websocket"
	"github.com/AliKefall/Somnambulist/internal/services/websocket/handlers"
	_ "github.com/jackc/pgx/v5/stdlib"
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

	redisClient := NewRedisClient(config.RedisURL)

	// Services
	chatService := chat.NewService(conn, queries)

	// Websocket Hub
	hub := websocket.NewHub(
		queries,
		metrics,
		chatService,
	)

	// Websocket event handlers
	wsHandler := &handlers.Config{
		Hub:         hub,
		Queries:     queries,
		RedisClient: redisClient,
	}

	hub.SetEventHandler(wsHandler)

	// HTTP handlers
	handler := endpoints.NewConfig(
		hub,
		conn,
		queries,
		chatService,
		auth.NewJWTManager(
			config.JWTSecret,
			15*time.Minute,
		),
		auth.NewPasswordHasher(),
		redisClient,
	)

	return conn, serverDependencies{
		config:  handler,
		queries: queries,
		hub:     hub,
		metrics: metrics,
		redis:   redisClient,
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

	var opt *redis.Options
	var err error

	if redisURL != "" {
		opt, err = redis.ParseURL(redisURL)
		if err != nil {
			log.Fatalf("invalid REDIS_URL: %v", err)
		}
	} else {
		opt = &redis.Options{
			Addr:         "localhost:6379",
			Password:     "",
			DB:           0,
			PoolSize:     10,
			MaxIdleConns: 5,
		}
	}

	rdb := redis.NewClient(opt)

	// optional but recommended: early fail
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("redis connection failed: %v", err)
	}

	return rdb
}
