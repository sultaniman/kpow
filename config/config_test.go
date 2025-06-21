package config

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Check if environment variables
// picked up and set in config
func TestConfigInit(t *testing.T) {
	keyDir := t.TempDir()
	keyPath := path.Join(keyDir, "pubkey.pub")
	if err := os.WriteFile(keyPath, []byte("pubkey"), 0o644); err != nil {
		t.Fatalf("failed to create pubkey: %v", err)
	}
	t.Setenv("TEST_KEYS_DIR", keyDir)
	t.Setenv("PORT", "10001")
	t.Setenv("LOG_LEVEL", "ERROR")
	t.Setenv("MAILER_FROM", "max@mustermann.de")
	t.Setenv("MAILER_TO", "me@mustermann.de")
	t.Setenv("MAILER_DSN", "smtp://user:secret@mustermann.de:587")
	t.Setenv("KEY_KIND", "age")
	t.Setenv("KEY_PATH", keyPath)
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
	assert.Equal(t, keyPath, config.Key.Path)
	assert.Equal(t, "/opt/kpow/inbox", config.Inbox.Path)
	assert.Equal(t, "5 * * * *", config.Inbox.Cron)
}

func TestConfigValidate(t *testing.T) {
	// Use a non-existent key path in a temporary directory
	dir := t.TempDir()
	cfg := &Config{
		Key: KeyInfo{
			Path: path.Join(dir, "missing.pub"),
			Kind: Age,
		},
		Mailer: Mailer{
			From: "from@example.com",
			To:   "to@example.com",
		},
	}

	errs := cfg.Validate()
	assert.Greater(t, len(errs), 0)

	var messages []string
	for _, err := range errs {
		messages = append(messages, err.Error())
	}
	combined := strings.Join(messages, " ")

	assert.Contains(t, combined, "public key file does not exist")
	assert.Contains(t, combined, "unable to read pubkey")
	assert.Contains(t, combined, "mailer dsn is required")
	assert.Contains(t, combined, "only smpt servers supported")
}
