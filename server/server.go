package server

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

func CreateServer(config *Config) *fiber.App {
	app := fiber.New(fiber.Config{})
	log.Debug().Msg("Server instance")
	return app
}
