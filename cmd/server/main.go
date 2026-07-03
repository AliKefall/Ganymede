package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	config := NewServer()
	conn, deps := bootstrapServer(config)
	defer conn.Close()

	go deps.hub.Run(context.Background())

	router := buildRouter(config, deps)
	srv := http.Server{
		Addr:    ":" + config.Port,
		Handler: router,
	}

	go waitForShutdown(&srv)

	log.Printf("Server listening on: %s", config.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}

}

// Graceful shotdown function
func waitForShutdown(srv *http.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop // This channel only works if it is in the main function

	log.Println("Shutting down the server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Shutdown error: %v", err)

		_ = srv.Close()
	}

	signal.Stop(stop)
	log.Println("Server stopped cleanly")

}
