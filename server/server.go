package server

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
	"github.com/sultaniman/kpow/config"
)

// go:embed public/kpow.min.css
var css string

// go:embed form.html
var formTemplate string

func CreateServer(config *config.Config) *fiber.App {
	app := fiber.New(fiber.Config{})
	log.Debug().Msg("Server instance")
	return app
}
