package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type Message struct {
	Subject      string `form:"subject" validate:"required"`
	SubjectError string
	Content      string `form:"content" validate:"required"`
	ContentError string
}

func (m *Message) Check() {
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

		message.Check()
		formData.Message = *message
	}

	if h.Config.Key.Advertise && len(h.Config.Key.KeyBytes) > 0 {
		formData.PubKey = string(h.Config.Key.KeyBytes)
	}

	err := ctx.Render(http.StatusOK, "form.html", formData)
	if err != nil {
		log.Err(err).Msg("Template rendering error")
	}

	return err
}
