package main

import (
	"log"
	"net/http"
	"os"

	"github.com/sffinn/restless/internal/api"
	"github.com/sffinn/restless/internal/datastore"
)

func main() {
	store := datastore.NewStore()

	// Seed with a couple of items so the API is immediately useful.
	store.Create("Deploy to DigitalOcean App Platform")
	store.Create("Write clean REST API in Go")

	handler := api.New(store)

	// App Platform injects the PORT environment variable.
	// Default to 8080 for local development.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := "0.0.0.0:" + port
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}
