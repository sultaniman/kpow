package form

import (
	"crypto/sha256"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/sultaniman/kpow/config"
)

type MessageForm struct {
	Subject      string `form:"subject" validate:"required"`
	SubjectError string
	Content      string `form:"content" validate:"required"`
	ContentError string
	IsValid      bool
}

func (m *MessageForm) Check() {
	if m.Subject == "" {
		m.SubjectError = "Subject is required"
	}

	if m.Content == "" {
		m.ContentError = "Message is required"
	}

	m.IsValid = m.SubjectError == "" && m.ContentError == ""
}

func (m *MessageForm) Hash() string {
	hash := sha256.New()
	hash.Write([]byte(m.Subject))
	hash.Write([]byte(m.Content))
	hashBytes := hash.Sum(nil)
	return fmt.Sprintf("%x", hashBytes)
}

type NoteKind string

const (
	ErrorPlain  NoteKind = "plain"
	ErrorNote   NoteKind = "error"
	SuccessNote NoteKind = "success"
)

type FormData struct {
	CSRFToken string
	Title     string

	Message MessageForm
	PubKey  string

	Note     string
	NoteKind NoteKind
}

func GetFormData(csrfToken string, config *config.Config) *FormData {
	form := &FormData{
		CSRFToken: csrfToken,
		Title:     config.Server.Title,
		PubKey:    "",
		Message:   MessageForm{},
	}

	if config.Key.Advertise && len(config.Key.KeyBytes) > 0 {
		form.PubKey = string(config.Key.KeyBytes)
	}

	return form
}

func BindFormMessage(ctx echo.Context) (*MessageForm, error) {
	if ctx.Request().Method != "POST" {
		return nil, nil
	}

	message := new(MessageForm)
	if err := ctx.Bind(message); err != nil {
		log.Warn().Err(err).Msg("Failed to bind form data")
		return nil, fmt.Errorf("invalid form data: %v", err)
	}

	message.Check()
	return message, nil
}
