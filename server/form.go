package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/sultaniman/kpow/server/form"
)

const CSRFTokenKey = "csrfToken"

func (h *Handler) RenderForm(ctx echo.Context) error {
	csrfToken := ctx.Get(CSRFTokenKey).(string)
	formData := form.GetFormData(csrfToken, h.Config)
	if messageForm, err := form.BindFormMessage(ctx); err == nil {
		formData.Message = *messageForm
	} else {
		formData.NoteKind = form.ErrorNote
		formData.Note = err.Error()
	}

	if formData.Message.IsValid {
		err := formData.EncryptAndSend(h.EncryptionProvider, h.Config.Inbox.Path)
		if err != nil {
			return h.internalError(err.Error())
		}

		formData.Note = "Wonderful! Your message is scheduled for delivery."
		formData.NoteKind = form.SuccessNote
	}

	err := ctx.Render(http.StatusOK, "form.html", formData)
	if err != nil {
		log.Err(err).Msg("Template rendering error")
	}

	return err
}
