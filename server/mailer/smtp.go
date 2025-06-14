package mailer

import (
	"fmt"

	"github.com/sultaniman/kpow/config"
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
	email.From(m.fromEmail)
	email.To(m.toEmail)
	email.SetBodyString(mail.TypeTextPlain, message.EncryptedMessage)
	if err := m.client.DialAndSend(email); err != nil {
		return fmt.Errorf("failed to send mail: %s", err)
	}
	return nil
}

func NewSMTPMailer(mailerConfig *MailerConfig) (Mailer, error) {
	tlsPolicy := mail.TLSMandatory
	if config.IsLocalhost(mailerConfig.Host) {
		tlsPolicy = mail.TLSOpportunistic
	}

	client, err := mail.NewClient(
		mailerConfig.Host,
		mail.WithPort(mailerConfig.Port),
		mail.WithTLSPortPolicy(tlsPolicy),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(mailerConfig.Username),
		mail.WithPassword(mailerConfig.Password),
	)

	if err != nil {
		return nil, err
	}

	return &SMTPMailer{
		client:    client,
		fromEmail: mailerConfig.FromEmail,
		toEmail:   mailerConfig.ToEmail,
	}, nil
}
