package mailer

import (
	"net/url"
	"strconv"

	"github.com/sultaniman/kpow/config"
)

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

func GetMailer(mailerConfig *config.Mailer) (Mailer, error) {
	parts, err := url.Parse(mailerConfig.DSN)
	if err != nil {
		return nil, err
	}
	port, err := strconv.ParseInt(parts.Port(), 10, 0)
	if err != nil {
		return nil, err
	}

	password, _ := parts.User.Password()
	switch parts.Scheme {
	case "smtp":
		return NewSMTPMailer(&MailerConfig{
			Host:      parts.Host,
			Port:      int(port),
			Username:  parts.User.Username(),
			Password:  password,
			FromEmail: mailerConfig.From,
			ToEmail:   mailerConfig.To,
		})
	case "dummy":
		return NewDummyMailer()
	default:
		return nil, nil
	}
}
