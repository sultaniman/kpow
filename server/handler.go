package server

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"

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

func getMailer(mailerConfig *config.Mailer) (mailer.Mailer, error) {
	parts, err := url.Parse(mailerConfig.DSN)
	if err != nil {
		return nil, err
	}
	port, err := strconv.ParseInt(parts.Port(), 10, 0)
	if err != nil {
		return nil, err
	}
	password, _ := parts.User.Password()
	switch parts.Scheme {
	case "smtp":
		return mailer.NewSMTPMailer(&mailer.MailerConfig{
			Host:      parts.Host,
			Port:      int(port),
			Username:  parts.User.Username(),
			Password:  password,
			FromEmail: mailerConfig.From,
			ToEmail:   mailerConfig.To,
		})
	case "dummy":
		return mailer.NewDummyMailer()
	default:
		return nil, nil
	}
}

func getWebhookHandler(webhookUrl string) (mailer.Mailer, error) {
	if webhookUrl == "" {
		return nil, nil
	}

	parts, err := url.Parse(webhookUrl)
	if err != nil {
		return nil, err
	}

	if parts.Scheme != "https" {
		return nil, errors.New("webhook url should be https only")
	}

	return mailer.NewWebhookMailer(webhookUrl), nil
}
func NewHandler(config *config.Config) (*Handler, error) {
	if crypt, err := enc.LoadKey(&config.Key); err == nil {
		mailerHandler, err := getMailer(&config.Mailer)
		if err != nil {
			return nil, err
		}

		webhookHandler, err := getWebhookHandler(config.Webhook.Url)
		if err != nil {
			return nil, err
		}

		return &Handler{
			Config:             config,
			EncryptionProvider: crypt,
			Mailer:             mailerHandler,
			WebhookHandler:     webhookHandler,
		}, nil
	} else {
		return nil, err
	}
}
