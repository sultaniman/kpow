package mailer

type MailerConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	UseTLS   bool
}

type Mailer interface {
	Send(message string) error
}

type SMTPMailer struct {
	config MailerConfig
}

func (m *SMTPMailer) Send(message string) error {
	// Implementation of sending email using SMTP
	return nil
}

func NewSMTPMailer(config MailerConfig) *SMTPMailer {
	return &SMTPMailer{
		config: config,
	}
}
