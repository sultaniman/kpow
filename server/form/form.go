package form

import (
	"crypto/sha256"
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
	Website      string `form:"website"` // honeypot field
	IsValid      bool
}

func (m *MessageForm) Check() {
	// honeypot: if the hidden website field is filled, it's a bot
	if m.Website != "" {
		m.IsValid = false
		return
	}

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
	Banner    string
	HideLogo  bool
	Message   MessageForm
	PubKey    string

	Note     string
	NoteKind NoteKind
}

func (f *FormData) EncryptAndSend(
	sender mailer.Mailer,
	webhookHandler mailer.Mailer,
	encryptionProvider enc.KeyLike,
	inboxPath string,
) error {
	subject := f.Message.Subject
	content := f.Message.Content
	hash := f.Message.Hash()

	encrypted, err := encryptionProvider.Encrypt(content)
	if err != nil {
		log.Err(err).Msg("encryption failed")
		return err
	}

	message := mailer.NewMessage(subject, encrypted, hash)

	// delivery is async since SMTP/webhook can be slow
	go func() {
		if err := mailer.SendMessage(message, sender, webhookHandler, inboxPath); err != nil {
			log.Err(err).Str("hash", hash).Msg("delivery failed, saved to inbox")
		}
	}()

	// reset the form
	f.Message = MessageForm{}
	return nil
}

func GetFormData(csrfToken string, config *config.Config) *FormData {
	form := &FormData{
		CSRFToken: csrfToken,
		Title:     config.Server.Title,
		Banner:    config.Server.CustomBanner,
		HideLogo:  config.Server.HideLogo,
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
