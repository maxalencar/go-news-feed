package news

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Server struct {
	mux    *http.ServeMux
	config Config
}

// NewServer - constructor
func NewServer() *Server {
	return &Server{}
}

func (s *Server) Init() error {
	ctx := context.Background()

	config, err := newConfig()
	if err != nil {
		return err
	}

	s.config = config

	repository, err := newRepository(ctx, config.MongoConfig)
	if err != nil {
		return err
	}

	service := newService(repository)
	endpoint := newEndpoint(service)

	s.mux = endpoint.init()

	return nil
}

func (s *Server) Start() error {
	// Start HTTP server
	addr := fmt.Sprintf(":%d", s.config.Server.Port)
	log.Printf("server listening on port %d...\n", s.config.Server.Port)

	server := &http.Server{
		Addr:              addr,
		ReadHeaderTimeout: 10 * time.Second,
		Handler:           s.mux,
	}

	return server.ListenAndServe()
}
