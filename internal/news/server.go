package news

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
)

type Server struct {
	echo   *echo.Echo
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

	s.echo = endpoint.init()

	return nil
}

func (s *Server) Start() error {
	return s.echo.Start(fmt.Sprintf(":%s", s.config.Server.Port))
}
