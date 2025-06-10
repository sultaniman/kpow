package form

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/sultaniman/kpow/config"
	"github.com/sultaniman/kpow/server/enc"
	"github.com/sultaniman/kpow/server/mailer"
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

func (f *FormData) EncryptAndSend(sender mailer.Mailer, wehbhooHandler mailer.Mailer, encryptionProvider enc.KeyLike, inboxPath string) error {
	encrypted, err := encryptionProvider.Encrypt(f.Message.Content)
	if err != nil {
		log.Err(err).Msg("Encryption failed")
		return errors.New("unable encrypt the message")
	}

	go (func() {
		message := mailer.NewMessage(f.Message.Subject, encrypted, f.Message.Hash())
		failed := false
		if err = sender.Send(message); err != nil {
			log.Err(err).Str("method", "mailer").Msg("Unable to send the message")
			message.Method = "mailer"
			failed = true
			err = message.Save(inboxPath)
			if err != nil {
				log.Err(err).Str("message", message.EncryptedMessage).Msg("Unable to save message")
			}
		}

		if !failed && wehbhooHandler != nil {
			if err = wehbhooHandler.Send(message); err != nil {
				log.Err(err).Str("method", "webhook").Msg("Unable to send the message")
				message.Method = "webhook"
				err = message.Save(inboxPath)
				if err != nil {
					log.Err(err).Str("message", message.EncryptedMessage).Msg("Unable to save message")

				}
			}
		}
	})()

	// when done reset the form
	f.Message = MessageForm{}
	return nil
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

// BindFormMessage
// Binds and validates subject and message in the submitted form
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
