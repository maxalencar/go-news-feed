package news

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"go-news-feed/pkg/model"
)

type endpoint struct {
	service Service
}

// newEndpoint - constructor
func newEndpoint(service Service) *endpoint {
	return &endpoint{
		service: service,
	}
}

func (e *endpoint) init() *echo.Echo {
	// Echo instance
	ei := echo.New()

	// Middlewares
	ei.Use(middleware.Logger())
	ei.Use(middleware.Recover())
	ei.Use(middleware.Secure())

	// Custom middleware to sanitize request
	ei.Use(Sanitize())

	// Routes
	ei.GET("/find", e.find)
	ei.GET("/load", e.load)

	return ei
}

func (e endpoint) find(c echo.Context) error {
	var fr model.FindRequest
	if err := c.Bind(&fr); err != nil {
		return echo.ErrBadRequest
	}

	response, err := e.service.Find(c.Request().Context(), fr)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
}

func (e endpoint) load(c echo.Context) error {
	response, err := e.service.Load(c.Request().Context(), c.QueryParam("feedUrl"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
}
