package mailer

import (
	"errors"
	"net/url"
	"strconv"
	"strings"

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
	if strings.HasPrefix(mailerConfig.DSN, "dummy://") {
		return NewDummyMailer()
	}

	if !strings.HasPrefix(mailerConfig.DSN, "smtp://") {
		return nil, errors.New("unsupported mailer option")
	}

	parts, err := url.Parse(mailerConfig.DSN)
	if err != nil {
		return nil, err
	}
	port, err := strconv.ParseInt(parts.Port(), 10, 0)
	if err != nil {
		return nil, err
	}

	password, _ := parts.User.Password()
	return NewSMTPMailer(&MailerConfig{
		Host:      parts.Host,
		Port:      int(port),
		Username:  parts.User.Username(),
		Password:  password,
		FromEmail: mailerConfig.From,
		ToEmail:   mailerConfig.To,
	})
}
