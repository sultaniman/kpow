package mailer

type MailerConfig struct {
	Host      string
	Port      int
	Username  string
	Password  string
	FromEmail string
	ToEmail   string
}

type Mailer interface {
	Send(message Message) error
}
