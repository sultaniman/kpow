package mailer

import (
	"fmt"

	"github.com/wneessen/go-mail"
)

type SMTPMailer struct {
	client    *mail.Client
	fromEmail string
	toEmail   string
}

func (m *SMTPMailer) Send(message Message) error {
	email := mail.NewMsg()
	email.Subject(message.Subject)
	email.SetBodyString(mail.TypeTextPlain, message.EncryptedMessage)
	if err := m.client.DialAndSend(email); err != nil {
		return fmt.Errorf("failed to send mail: %s", err)
	}
	return nil
}

func NewSMTPMailer(config MailerConfig) (*SMTPMailer, error) {
	client, err := mail.NewClient(
		config.Host,
		mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(config.Username),
		mail.WithPassword(config.Password),
	)
	if err != nil {
		return nil, err
	}

	return &SMTPMailer{
		client:    client,
		fromEmail: config.FromEmail,
		toEmail:   config.ToEmail,
	}, nil
}
