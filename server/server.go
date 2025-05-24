package server

import (
	"github.com/labstack/echo/v4"
	"github.com/sultaniman/kpow/config"
)

// go:embed public/kpow.min.css
var css string

// go:embed form.html
var formTemplate string

func CreateServer(config *config.Config) *echo.Echo {
	app := echo.New()
	app.HideBanner = true
	return app
}
