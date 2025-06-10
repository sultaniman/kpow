package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/sultaniman/kpow/config"
	"github.com/sultaniman/kpow/server/enc"
	"github.com/sultaniman/kpow/server/mailer"
)

type Handler struct {
	Config             *config.Config
	EncryptionProvider enc.KeyLike
	Mailer             mailer.Mailer
	WebhookHandler     mailer.Mailer
}

func (h *Handler) internalError(message string) *echo.HTTPError {
	return echo.NewHTTPError(
		http.StatusInternalServerError,
		message,
	)
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
	case http.StatusRequestEntityTooLarge:
		serverError.Title = "Your message is too big..."
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

func NewHandler(config *config.Config) (*Handler, error) {
	if cryptoHandler, err := enc.LoadKey(&config.Key); err == nil {
		mailerHandler, err := mailer.GetMailer(&config.Mailer)
		if err != nil {
			return nil, err
		}

		webhookHandler, err := mailer.NewWebhookMailer(config.Webhook.Url)
		if err != nil {
			return nil, err
		}

		return &Handler{
			Config:             config,
			EncryptionProvider: cryptoHandler,
			Mailer:             mailerHandler,
			WebhookHandler:     webhookHandler,
		}, nil
	} else {
		return nil, err
	}
}
