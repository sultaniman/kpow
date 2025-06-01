package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type Message struct {
	Subject      string
	SubjectError string
	Content      string
	ContentError string
}

type FormData struct {
	CSRFToken string
	Title     string

	Message Message
	PubKey  string

	Note    string
	IsError bool
}

const PubKeySample = `-----BEGIN AGE ENCRYPTED FILE-----
YWdlLWVuY3J5cHRpb24ub3JnL3YxCi0+IFgyNTUxOSA3WnpnL0h6OHp6bDEySHFq
azlWWW9EUmxNUXk4a29NOElUWERFSXJaNXhVClV4d3pocXNGL2ZNZnRwV0RvY3NT
STUzRHBmRzFJcEhhdmtFTGZabkcyZmcKLS0tIFhvNnJReVBnOXc4Zm1FOFdZdHhB
SzlEblMralZqcmloSXRucHNSV2Fqc1EK3KFfOX1Ln968kq1tX1iaI+9RoSqekVOF
na03n83y9DttvF2XOw==
-----END AGE ENCRYPTED FILE-----`

func (h *Handler) RenderForm(ctx echo.Context) error {
	formData := FormData{
		CSRFToken: ctx.Get("csrfToken").(string),
		Title:     h.Config.Server.Title,
		PubKey:    PubKeySample,
		Message:   Message{},
	}

	if ctx.Request().Method == "POST" {
		message := new(Message)
		if err := ctx.Bind(message); err != nil {
			formData.Note = fmt.Sprintf("invalid form data: %v, ", err)
			formData.IsError = true
		}

		if message.Subject == "" {
			message.SubjectError = "Subject is required"
		}

		if message.Content == "" {
			message.ContentError = "Message is required"
		}

		formData.Message = *message
	}

	err := ctx.Render(http.StatusOK, "form.html", formData)
	if err != nil {
		log.Err(err).Msg("Template rendering error")
	}

	return err
}
