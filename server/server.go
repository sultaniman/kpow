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
//go:embed templates/*
var resources embed.FS

func CreateServer(config *config.Config) (*echo.Echo, error) {
	app := echo.New()
	app.HideBanner = true
	templates, err := template.ParseFS(resources, "templates/*.html")
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
	app.Match(allowedFormMethods, "/", RenderForm)

	return app, nil
}
