package server

import (
	"errors"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/spf13/viper"
)

type Kind string

const (
	PGP Kind = "pgp"
	Age Kind = "age"
)

const (
	Port = 8080
	Host = "localhost"
)

// For kind=PGP and unset path
// password pgp encryption is used.
type KeyInfo struct {
	Path         string
	KeyKind      Kind
	Password     string
	AdvertiseKey bool
}

type Config struct {
	Title      string
	Port       int
	Host       string
	LogLevel   string
	KeyInfo    KeyInfo
	MailerFrom string
	MailerTo   string
	MailerDSN  string
}

func (c *Config) Validate() error {
	if c.KeyInfo.KeyKind != Age && c.KeyInfo.KeyKind != PGP {
		return errors.New(fmt.Sprintf("unknown key kind %s", c.KeyInfo.KeyKind))
	}

	if c.KeyInfo.Path == "" && c.KeyInfo.Password == "" {
		return errors.New("key path or password is required")
	}

	if _, err := os.Stat(c.KeyInfo.Path); errors.Is(err, os.ErrNotExist) {
		return errors.New("public key file does not exist")
	}

	if c.MailerFrom == "" {
		return errors.New("mailer from is required")
	}

	if c.MailerTo == "" {
		return errors.New("recipient email is required")
	}

	if c.MailerDSN == "" {
		return errors.New("mailer dsn is required")
	}

	return nil
}

func GetConfig(path string) (*Config, error) {
	// If config path is empty then we load the default config
	// Else we load the config from the given path
	// Then assign values from viper later they will be substituted
	// with relevant CLI arguments
	var config *Config
	if path == "" {
		config = &Config{
			Title:    "KPow",
			Port:     Port,
			Host:     Host,
			LogLevel: "info",
			KeyInfo:  KeyInfo{AdvertiseKey: true},
		}
	} else {
		config = &Config{}
		if _, err := toml.DecodeFile(path, config); err != nil {
			return nil, err
		}
	}

	// server
	config.Title = viper.GetString("title")
	config.Port = viper.GetInt("port")
	config.Host = viper.GetString("host")
	config.LogLevel = viper.GetString("log_level")

	// mailer
	config.MailerFrom = viper.GetString("mailer_from")
	config.MailerTo = viper.GetString("mailer_to")
	config.MailerDSN = viper.GetString("mailer_dsn")

	// key
	config.KeyInfo.KeyKind = Kind(viper.GetString("key_kind"))
	config.KeyInfo.Password = viper.GetString("password")
	config.KeyInfo.AdvertiseKey = viper.GetBool("advertise")
	config.KeyInfo.Path = viper.GetString("pubkey_path")

	return config, nil
}
