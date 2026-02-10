package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"devops/backend/internal/api"
	"devops/backend/internal/config"
	"devops/backend/internal/store"
)

func main() {
	addr := ":8080"
	if v := os.Getenv("PORT"); v != "" {
		addr = ":" + v
	}

	staticDir := os.Getenv("STATIC_DIR")
	hostsFile := os.Getenv("HOSTS_FILE")
	security, err := config.ParseSecurity(os.Getenv("API_KEYS"))
	if err != nil {
		log.Fatalf("invalid API_KEYS: %v", err)
	}

	hosts, err := config.LoadHosts(hostsFile)
	if err != nil {
		log.Fatalf("failed to load hosts from %s: %v", hostsFile, err)
	}

	handler := api.New(store.NewWithHosts(hosts), staticDir, security).Router()
	srv := &http.Server{Addr: addr, Handler: handler, ReadHeaderTimeout: 10 * time.Second}

	log.Printf("starting API server on %s", addr)
	if staticDir != "" {
		log.Printf("serving static UI from %s", staticDir)
	}
	if hostsFile != "" {
		log.Printf("loading KVM hosts inventory from %s", hostsFile)
	}
	if len(security.APIKeys) == 0 {
		log.Printf("WARNING: API auth disabled (API_KEYS not set)")
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Printf("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown error: %v", err)
	}
}
