package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type Message struct {
	Subject string
	Content string
}

type FormData struct {
	Title string

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
	err := ctx.Render(http.StatusOK, "form.html", FormData{
		Title: h.Config.Server.Title,
		Message: Message{
			Subject: "Subject",
			Content: "Content...",
		},
		PubKey: PubKeySample,
	})

	if err != nil {
		log.Err(err).Msg("Template rendering error")
	}

	return err
}
