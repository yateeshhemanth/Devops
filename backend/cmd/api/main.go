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
	s := api.New(store.New())
	log.Printf("starting API server on %s", addr)
	if err := http.ListenAndServe(addr, s.Router()); err != nil {
		log.Fatal(err)
	}
}
