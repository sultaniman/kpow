package server

import "github.com/sultaniman/kpow/config"

type Handler struct {
	Config *config.Config
}

func NewHandler(config *config.Config) Handler {
	return Handler{
		Config: config,
	}
}
