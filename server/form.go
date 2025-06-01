package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/sultaniman/kpow/config"
)

type Message struct {
	Subject      string `form:"subject" validate:"required"`
	SubjectError string
	Content      string `form:"content" validate:"required"`
	ContentError string
}

func (m *Message) CheckForm() {
	if m.Subject == "" {
		m.SubjectError = "Subject is required"
	}

	if m.Content == "" {
		m.ContentError = "Message is required"
	}
}

type FormData struct {
	CSRFToken string
	Title     string

	Message Message
	PubKey  string

	Note    string
	IsError bool
}

const CSRFTokenKey = "csrfToken"

func (h *Handler) RenderForm(ctx echo.Context) error {
	csrfToken := ctx.Get(CSRFTokenKey).(string)
	formData := FormData{
		CSRFToken: csrfToken,
		Title:     h.Config.Server.Title,
		PubKey:    "",
		Message:   Message{},
	}

	if ctx.Request().Method == "POST" {
		message := new(Message)
		if err := ctx.Bind(message); err != nil {
			log.Warn().Err(err).Msg("Failed to bind form data")

			formData.Note = fmt.Sprintf("Invalid form data: %v", err)
			formData.IsError = true
		}

		message.CheckForm()
		formData.Message = *message
	}

	if h.Config.Key.Advertise && h.Config.Key.Kind == config.PGP {
		pubkey, err := os.ReadFile(h.Config.Key.Path)
		if err != nil {
			log.Fatal().Err(err).Msgf("Unable to read public key: %s", h.Config.Key.Path)
		}

		formData.PubKey = string(pubkey)
	}

	err := ctx.Render(http.StatusOK, "form.html", formData)
	if err != nil {
		log.Err(err).Msg("Template rendering error")
	}

	return err
}
