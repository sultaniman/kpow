package server

import (
	"embed"
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/sultaniman/kpow/config"
	"golang.org/x/time/rate"
)

//go:embed public/*
//go:embed templates/* templates/icons/*
var resources embed.FS

func CreateServer(conf *config.Config) (*echo.Echo, error) {
	app := echo.New()
	app.HideBanner = true
	handler, err := NewHandler(conf)
	if err != nil {
		return nil, err
	}

	templates, err := template.ParseFS(resources, "templates/*.html", "templates/icons/*.svg")
	if err != nil {
		return nil, err
	}

	app.Renderer = &Template{
		templates: templates,
	}

	app.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:       "public",
		HTML5:      false,
		Filesystem: http.FS(resources),
	}))

	var formMiddlewares = []echo.MiddlewareFunc{
		middleware.BodyLimitWithConfig(middleware.BodyLimitConfig{
			Limit: fmt.Sprintf("%dB", conf.Server.MessageSize),
		}),
		middleware.CSRFWithConfig(middleware.CSRFConfig{
			TokenLookup: "form:csrf",
			ContextKey:  "csrfToken",
		}),
	}
	// If it is set to 0 or less then we don't enable rate limiting
	if conf.RateLimiter != nil && conf.RateLimiter.RPM > 0 {
		log.
			Info().
			Int("rpm", conf.RateLimiter.RPM).
			Int("burst", conf.RateLimiter.Burst).
			Int("cooldown", conf.RateLimiter.CooldownSeconds).
			Msg("Rate limiting enabled")

		formMiddlewares = append(formMiddlewares, middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
			Skipper: middleware.DefaultSkipper,
			Store: middleware.NewRateLimiterMemoryStoreWithConfig(
				middleware.RateLimiterMemoryStoreConfig{
					Rate:      rate.Limit(conf.RateLimiter.RPM),                // requests per minute
					Burst:     conf.RateLimiter.Burst,                          // max burst
					ExpiresIn: time.Duration(conf.RateLimiter.CooldownSeconds), // keep IP in memory for cooldown
				},
			),
			IdentifierExtractor: func(c echo.Context) (string, error) {
				return c.RealIP(), nil
			},
			ErrorHandler: func(c echo.Context, err error) error {
				return &echo.HTTPError{
					Code:     http.StatusForbidden,
					Message:  "Unable to read real ip",
					Internal: err,
				}
			},
			DenyHandler: func(context echo.Context, identifier string, err error) error {
				return &echo.HTTPError{
					Code:     http.StatusTooManyRequests,
					Message:  "Too many requests",
					Internal: err,
				}
			},
		}))
	}

	app.Match(
		[]string{"GET", "POST"},
		"/",
		handler.RenderForm,
		formMiddlewares...,
	)

	app.HTTPErrorHandler = errorHandler

	return app, nil
}
