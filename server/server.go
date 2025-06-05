package server

import (
	"embed"
	"net/http"
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sultaniman/kpow/config"
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

	app.Use(middleware.BodyLimitWithConfig(middleware.BodyLimitConfig{
		Limit: "0.5KI", // 512 bytes
	}))

	app.Match(
		[]string{"GET", "POST"},
		"/",
		handler.RenderForm,
		middleware.CSRFWithConfig(middleware.CSRFConfig{
			TokenLookup: "form:csrf",
			ContextKey:  "csrfToken",
		}),
	)

	app.HTTPErrorHandler = errorHandler

	return app, nil
}
