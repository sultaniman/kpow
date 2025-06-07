package mailer

import (
	"github.com/rs/zerolog/log"
)

type DummyMailer struct{}

func (m *DummyMailer) Send(message Message) error {
	log.
		Debug().
		Str("message", message.EncryptedMessage).
		Str("subject", message.Subject).
		Msg("Send message")

	return nil
}

func NewDummyMailer() (*DummyMailer, error) {
	return &DummyMailer{}, nil
}
