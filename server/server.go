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

func CreateServer(config *config.Config) (*echo.Echo, error) {
	app := echo.New()
	app.HideBanner = true
	handler := NewHandler(config)
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
	allowedFormMethods := []string{"GET", "POST"}
	app.Match(
		allowedFormMethods,
		"/",
		handler.RenderForm,
		middleware.CSRFWithConfig(middleware.CSRFConfig{
			TokenLookup: "form:csrf",
			ContextKey:  "csrfToken",
		}),
	)

	return app, nil
}
