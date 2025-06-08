package config

import (
	"errors"
	"fmt"
	"net/mail"
	"net/url"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/rs/zerolog"
	"github.com/sultaniman/env"
)

type KeyKind string

const (
	PGP KeyKind = "pgp"
	Age KeyKind = "age"
)

const (
	Port = 8080
	Host = "localhost"
)

// For kind=PGP and unset path
// password pgp encryption is used.
type KeyInfo struct {
	Path      string
	KeyBytes  []byte
	Kind      KeyKind
	Advertise bool
}

type ServerConfig struct {
	Title        string
	Port         int
	Host         string
	MessageSize  int    `toml:"message_size"`
	HideLogo     bool   `toml:"hide_logo"`
	CustomBanner string `toml:"custom_banner"`
	LogLevel     string `toml:"log_level"`
}

type Mailer struct {
	From string
	To   string
	DSN  string
}

type Webhook struct {
	Url string
}

type Inbox struct {
	Path string
	Cron string
	// We want to send messages in batches
	// because otherwise we might ddos the
	// receiving side/server.
	BatchSize int `toml:"batch_size"`
}

type RateLimiter struct {
	RequestsPerMinute int
	NumBurstRequests  int
	CooldownSeconds   int
}

type Config struct {
	Server      ServerConfig
	Key         KeyInfo
	Mailer      Mailer
	RateLimiter *RateLimiter

	// Use webhook instead of mailer
	Webhook Webhook

	// Inbox config
	Inbox Inbox
}

func (c *Config) Validate() []error {
	var errorList = []error{}
	if c.Key.Kind != Age && c.Key.Kind != PGP {
		errorList = append(
			errorList,
			newConfigError("KEY_KIND", fmt.Sprintf("unsupported key kind %s", c.Key.Kind)),
		)
	}

	if c.Key.Path == "" && c.Key.Password == "" {
		errorList = append(
			errorList,
			newConfigError("KEY_PATH", "key path or password is required"),
		)
	}

	if _, err := os.Stat(c.Key.Path); errors.Is(err, os.ErrNotExist) {
		errorList = append(errorList, errors.New("public key file does not exist"))
	}

	if c.Mailer.From == "" {
		errorList = append(errorList, errors.New("mailer from is required"))
	}

	if _, err := mail.ParseAddress(c.Mailer.From); err != nil {
		errorList = append(errorList, errors.New("invalid sender address"))
	}

	if c.Mailer.To == "" {
		errorList = append(errorList, errors.New("recipient email is required"))
	}

	if _, err := mail.ParseAddress(c.Mailer.To); err != nil {
		errorList = append(errorList, errors.New("invalid recipient address"))
	}

	if c.Mailer.DSN == "" {
		errorList = append(errorList, errors.New("mailer dsn is required"))
	}

	parts, err := url.Parse(c.Mailer.DSN)
	if err != nil {
		errorList = append(errorList, errors.New("invalid mailer dsn"))
	}

	if parts.Scheme != "smtp" {
		errorList = append(errorList, errors.New("only smpt servers supported"))
	}

	return errorList
}

func (c *Config) ParseLogLevel(level string) (zerolog.Level, error) {
	logLevel, err := zerolog.ParseLevel(level)

	if err != nil {
		return 0, err
	}

	c.Server.LogLevel = level
	return logLevel, nil
}

// GetConfig loads configuration from toml file
// then substitutes values with the ones from environment.
func GetConfig(path string) (*Config, error) {
	var config = &Config{}

	if path != "" {
		if _, err := toml.DecodeFile(path, config); err != nil {
			return nil, err
		}
	}

	// server
	if title := env.GetString("TITLE"); title != "" {
		config.Server.Title = title
	}

	if serverPort, err := env.GetIntE("PORT"); err == nil {
		config.Server.Port = serverPort
	}

	if serverHost := env.GetString("HOST"); serverHost != "" {
		config.Server.Host = serverHost
	}

	if logLevel := env.GetString("LOG_LEVEL"); logLevel != "" {
		config.Server.LogLevel = logLevel
	}

	// mailer
	if fromEmail := env.GetString("MAILER_FROM"); fromEmail != "" {
		config.Mailer.From = fromEmail
	}

	if toEmail := env.GetString("MAILER_TO"); toEmail != "" {
		config.Mailer.To = toEmail
	}

	if mailerDSN := env.GetString("MAILER_DSN"); mailerDSN != "" {
		config.Mailer.DSN = mailerDSN
	}

	// key
	if keyKind := env.GetString("KEY_KIND"); keyKind != "" {
		config.Key.Kind = KeyKind(env.GetString("KEY_KIND"))
	}

	config.Key.Advertise = config.Key.Advertise || env.GetBool("ADVERTISE")

	if keyPath := env.GetString("KEY_PATH"); keyPath != "" {
		config.Key.Path = keyPath
	}

	if inboxPath := env.GetString("INBOX_PATH"); inboxPath != "" {
		config.Inbox.Path = inboxPath
	}

	if inboxCron := env.GetString("INBOX_CRON"); inboxCron != "" {
		config.Inbox.Cron = inboxCron
	}

	if keyBytes, err := os.ReadFile(config.Key.Path); err == nil {
		config.Key.KeyBytes = keyBytes
	} else {
		return nil, err
	}

	return config, nil
}
