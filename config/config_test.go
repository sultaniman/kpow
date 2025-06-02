package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Check if environment variables
// picked up and set in config
func TestConfigInit(t *testing.T) {
	t.Setenv("PORT", "10001")
	t.Setenv("LOG_LEVEL", "ERROR")
	t.Setenv("MAILER_FROM", "max@mustermann.de")
	t.Setenv("MAILER_TO", "me@mustermann.de")
	t.Setenv("MAILER_DSN", "smtp://user:secret@mustermann.de:587")
	t.Setenv("KEY_KIND", "age")
	t.Setenv("KEY_PATH", "/opt/kpow/pub.pgp")
	t.Setenv("INBOX_PATH", "/opt/kpow/inbox")
	t.Setenv("INBOX_CRON", "5 * * * *")

	config, err := GetConfig("")
	assert.NoError(t, err)
	assert.Equal(t, 10001, config.Server.Port)
	assert.Equal(t, "ERROR", config.Server.LogLevel)
	assert.Equal(t, "max@mustermann.de", config.Mailer.From)
	assert.Equal(t, "me@mustermann.de", config.Mailer.To)
	assert.Equal(t, "smtp://user:secret@mustermann.de:587", config.Mailer.DSN)
	assert.Equal(t, "age", string(config.Key.Kind))
	assert.Equal(t, "/opt/kpow/pub.pgp", config.Key.Path)
	assert.Equal(t, "/opt/kpow/inbox", config.Inbox.Path)
	assert.Equal(t, "5 * * * *", config.Inbox.Cron)
}

func TestConfigValidate(t *testing.T) {

}
