package server

import "github.com/BurntSushi/toml"

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
	Title     string
	Port      int
	Host      string
	LogLevel  string
	KeyInfo   KeyInfo
	MailerDSN string
}

func FromToml(path string) (*Config, error) {
	var config = &Config{}
	if _, err := toml.DecodeFile(path, config); err != nil {
		return nil, err
	}

	return config, nil
}

func NewConfig() *Config {
	return &Config{
		Title:    "KPow",
		Port:     Port,
		Host:     Host,
		LogLevel: "info",
		KeyInfo:  KeyInfo{AdvertiseKey: true},
	}
}
