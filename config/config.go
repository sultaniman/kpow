package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"net/mail"
	"net/url"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sultaniman/env"
)

type KeyKind string

const (
	PGP KeyKind = "pgp"
	Age KeyKind = "age"
	RSA KeyKind = "rsa"
)

const (
	Port               = 8080
	Host               = "localhost"
	DefaultMessageSize = 120
)

// For kind=PGP and unset path
// password pgp encryption is used.
type KeyInfo struct {
	Path      string
	KeyBytes  []byte `toml:"-"`
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
	From       string
	To         string
	DSN        string
	MaxRetries int `toml:"max_retries"`
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
	RPM             int `toml:"rpm"`
	Burst           int `toml:"burst"`
	CooldownSeconds int `toml:"cooldown"`
}

type Config struct {
	Server      ServerConfig
	Key         KeyInfo
	Mailer      Mailer
	RateLimiter *RateLimiter `toml:"rate_limiter"`

	// Use webhook instead of mailer
	Webhook Webhook

	// Inbox config
	Inbox Inbox
}

func (c *Config) Validate() []error {
	var errorList = []error{}

	absPath, err := filepath.Abs(c.Key.Path)
	if err != nil {
		errorList = append(
			errorList,
			newConfigError(
				"KEY_PATH",
				fmt.Sprintf("invalid key path %s", c.Key.Path),
			),
		)
	}

	if _, err := os.Stat(absPath); errors.Is(err, os.ErrNotExist) {
		errorList = append(
			errorList,
			newConfigError(
				"KEY_PATH",
				fmt.Sprintf("public key file does not exist %s", c.Key.Path),
			),
		)
	}

	if keyBytes, err := os.ReadFile(absPath); err == nil {
		c.Key.KeyBytes = keyBytes
	} else {
		errorList = append(
			errorList,
			newConfigError("KEY_PATH", fmt.Sprintf("unable to read pubkey %s", c.Key.Path)),
		)
	}

	if c.Key.Kind != Age && c.Key.Kind != PGP && c.Key.Kind != RSA {
		errorList = append(
			errorList,
			newConfigError("KEY_KIND", fmt.Sprintf("unsupported key kind %s", c.Key.Kind)),
		)
	}

	// Validate if rsa key bits > selected message size
	if c.Key.Kind == RSA {
		keyIsValid := true
		block, _ := pem.Decode(c.Key.KeyBytes)
		parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)

		if err != nil {
			keyIsValid = false
			errorList = append(
				errorList,
				newConfigError("KEY_KIND", err.Error()),
			)
		}

		rsaKey, ok := parsedKey.(*rsa.PublicKey)
		if !ok {
			keyIsValid = false
			errorList = append(
				errorList,
				newConfigError("KEY_KIND", "invalid rsa public key"),
			)
		}

		if rsaKey != nil {
			keyBits := rsaKey.N.BitLen()
			if keyIsValid && keyBits>>3 < c.Server.MessageSize {
				errorList = append(
					errorList,
					newConfigError(
						"KEY_KIND",
						fmt.Sprintf(
							"public key bytes %d can not be less than message bytes %d",
							keyBits>>3, c.Server.MessageSize,
						),
					),
				)
			}
		}
	}

	if c.Server.MessageSize > 0 && c.Server.MessageSize < DefaultMessageSize {
		log.
			Warn().
			Int("message_size", c.Server.MessageSize).
			Msg("Message size is too small")
	} else if c.Server.MessageSize <= 0 {
		log.
			Warn().
			Int("message_size", c.Server.MessageSize).
			Msgf("Incorrect message size, using default %d bytes", DefaultMessageSize)

		c.Server.MessageSize = DefaultMessageSize
	}

	// validate mailer options
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

	// validate webhook url
	if c.Webhook.Url != "" {
		parts, err := url.Parse(c.Webhook.Url)
		if err != nil {
			errorList = append(errorList, errors.New("invalid webhook url"))
		}

		hostname := parts.Hostname()
		if parts.Scheme != "https" && !IsLocalhost(hostname) {
			errorList = append(errorList, errors.New("webhook url should use https"))
		}
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

	if config.RateLimiter == nil {
		config.RateLimiter = &RateLimiter{}
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

	if messageSize := env.GetInt("MESSAGE_SIZE"); messageSize > 0 {
		config.Server.MessageSize = messageSize
	}

	if hideLogo := env.GetBool("HIDE_LOGO"); hideLogo {
		config.Server.HideLogo = hideLogo
	}

	if customBanner := env.GetString("CUSTOM_BANNER"); customBanner != "" {
		config.Server.CustomBanner = customBanner
	}

	// rate limiter
	if rpm := env.GetInt("LIMITER_RPM"); rpm > 0 {
		config.RateLimiter.RPM = rpm
	}

	if numBurstRequests := env.GetInt("LIMITER_BURST"); numBurstRequests > 0 {
		config.RateLimiter.Burst = numBurstRequests
	}

	if rateLimitCooldownSeconds := env.GetInt("LIMITER_COOLDOWN"); rateLimitCooldownSeconds > 0 {
		config.RateLimiter.CooldownSeconds = rateLimitCooldownSeconds
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

	if webhookUrl := env.GetString("WEBHOOK_URL"); webhookUrl != "" {
		config.Webhook.Url = webhookUrl
	}

	if maxRetries := env.GetInt("MAX_RETRIES"); maxRetries > 0 {
		config.Mailer.MaxRetries = maxRetries
	}

	// key
	if keyKind := env.GetString("KEY_KIND"); keyKind != "" {
		config.Key.Kind = KeyKind(keyKind)
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

	if inboxBatchSize := env.GetInt("INBOX_BATCH_SIZE"); inboxBatchSize > 0 {
		config.Inbox.BatchSize = inboxBatchSize
	}

	return config, nil
}
