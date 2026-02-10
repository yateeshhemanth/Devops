package main

import (
	"log"
	"net/http"
	"os"

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

	hosts, err := config.LoadHosts(hostsFile)
	if err != nil {
		log.Fatalf("failed to load hosts from %s: %v", hostsFile, err)
	}

	s := api.New(store.NewWithHosts(hosts), staticDir)
	log.Printf("starting API server on %s", addr)
	if staticDir != "" {
		log.Printf("serving static UI from %s", staticDir)
	}
	if hostsFile != "" {
		log.Printf("loading KVM hosts inventory from %s", hostsFile)
	}

	if err := http.ListenAndServe(addr, s.Router()); err != nil {
		log.Fatal(err)
	}
}
