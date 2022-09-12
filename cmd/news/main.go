package main

import (
	"log"

	"go-news-feed/internal/news"
)

func main() {
	srv := news.NewServer()

	// Init server's dependencies
	if err := srv.Init(); err != nil {
		log.Fatalf("error initialising server. err: %v", err)
	}

	// Start server
	if err := srv.Start(); err != nil {
		log.Fatalf("error starting server. err: %v", err)
	}
}
