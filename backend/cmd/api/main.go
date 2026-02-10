package main

import (
	"log"
	"net/http"
	"os"

	"devops/backend/internal/api"
	"devops/backend/internal/store"
)

func main() {
	addr := ":8080"
	if v := os.Getenv("PORT"); v != "" {
		addr = ":" + v
	}

	staticDir := os.Getenv("STATIC_DIR")
	s := api.New(store.New(), staticDir)
	log.Printf("starting API server on %s", addr)
	if staticDir != "" {
		log.Printf("serving static UI from %s", staticDir)
	}

	if err := http.ListenAndServe(addr, s.Router()); err != nil {
		log.Fatal(err)
	}
}
