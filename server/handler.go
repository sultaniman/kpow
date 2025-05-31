package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/sultaniman/kpow/config"
)

type Handler struct {
	Config *config.Config
}

type ServerError struct {
	Code   int
	Title  string
	Reason string
}

func errorHandler(err error, ctx echo.Context) {
	if ctx.Response().Committed {
		return
	}

	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	serverError := ServerError{
		Code:   code,
		Title:  "Server Error",
		Reason: "",
	}

	switch code {
	case http.StatusForbidden:
		serverError.Title = "CSRF Error"
		serverError.Reason = "Invalid CSRF token"
	case http.StatusNotFound:
		serverError.Title = "Not Found"
	default:
		serverError.Title = "Unknown Error"
		serverError.Reason = "Oopsie! Unknown Error"
	}

	log.
		Err(err).
		Int("code", code).
		Str("URL", ctx.Request().RequestURI).
		Msg("")

	ctx.Render(code, "error.html", serverError)
}

func NewHandler(config *config.Config) Handler {
	return Handler{
		Config: config,
	}
}
