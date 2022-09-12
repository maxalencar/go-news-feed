package news

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	// SanitizeConfig defines the config for sanitize middleware.
	SanitizeConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper middleware.Skipper

		// Sanitizer defines a funciton to sanitize the params or a target param.
		// Optional. Default values are: ["*", "/", "\\", "$", "^"].
		Sanitizer func(string) string

		// TargetPathParam defines which param to look for
		TargetPathParam string

		// TargetQueryParam definies which param to look for
		TargetQueryParam string
	}
)

var (
	DefaultSanitizeConfig = SanitizeConfig{
		Skipper:   middleware.DefaultSkipper,
		Sanitizer: sanitizer,
	}

	DefaultSanitizeTokens = []string{"*", "/", "\\", "$", "^"}
)

func Sanitize() echo.MiddlewareFunc {
	return SanitizeWithConfig(DefaultSanitizeConfig)
}

func SanitizeWithConfig(config SanitizeConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultSanitizeConfig.Skipper
	}

	if config.Sanitizer == nil {
		config.Sanitizer = sanitizer
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			paramValues := c.ParamValues()

			for i, v := range c.ParamNames() {
				if config.TargetPathParam != "" && v != config.TargetPathParam {
					continue
				}

				paramValues[i] = config.Sanitizer(paramValues[i])
			}

			c.SetParamValues(paramValues...)

			for key, values := range c.QueryParams() {
				if config.TargetQueryParam != "" && values[0] != config.TargetQueryParam {
					continue
				}

				c.QueryParams().Set(key, config.Sanitizer(values[0]))
			}

			return next(c)
		}
	}
}

func sanitizer(value string) string {
	for _, token := range DefaultSanitizeTokens {
		value = strings.ReplaceAll(value, token, "")
	}

	return value
}
