package cmd

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigurationOverrides(t *testing.T) {
	pwd, _ := os.Getwd()
	projectRoot := filepath.Dir(pwd)
	keyPath := path.Join(projectRoot, "server/enc/testkeys/pubkey.pub")
	configPath := path.Join(projectRoot, "config.toml")
	err := startCmd.ParseFlags([]string{
		"--port=10000",
		"--host=127.0.0.1",
		"--config=" + configPath,
		"--pubkey=" + keyPath,
		"--key-kind=pgp",
		"--advertise-key=true",
		"--limiter-rpm=100001",
		"--limiter-burst=10",
		"--limiter-cooldown=33",
		"--mailer-from=myform@kpow.friends",
		"--mailer-to=ping@me.friends",
		"--mailer-dsn=smtp://user:pass@me.friends:587",
		"--webhook-url=https://kpow.friends/callback",
		"--inbox-path=" + projectRoot,
		"--inbox-cron=*/10 * * * *",
		"--batch-size=4",
		"--log-level=ERROR",
		"--banner=" + projectRoot + "/banner.html",
		"--hide-logo=true",
		"--message-size=512",
	})

	assert.NoError(t, err)

	appConfig, err := getConfig()
	assert.NoError(t, err)

	// Server
	assert.Equal(t, 10000, appConfig.Server.Port)
	assert.Equal(t, "127.0.0.1", appConfig.Server.Host)
	assert.Equal(t, "ERROR", appConfig.Server.LogLevel)
	assert.Equal(t, true, appConfig.Server.HideLogo)
	assert.Equal(t, 512, appConfig.Server.MessageSize)
	assert.Equal(t, projectRoot+"/banner.html", appConfig.Server.CustomBanner)

	// Key
	assert.Equal(t, keyPath, appConfig.Key.Path)
	assert.Equal(t, "pgp", string(appConfig.Key.Kind))
	assert.Equal(t, true, appConfig.Key.Advertise)

	// Rate limiter
	assert.Equal(t, 100001, appConfig.RateLimiter.RPM)
	assert.Equal(t, 10, appConfig.RateLimiter.Burst)
	assert.Equal(t, 33, appConfig.RateLimiter.CooldownSeconds)

	// Mailer
	assert.Equal(t, "myform@kpow.friends", appConfig.Mailer.From)
	assert.Equal(t, "ping@me.friends", appConfig.Mailer.To)
	assert.Equal(t, "smtp://user:pass@me.friends:587", appConfig.Mailer.DSN)
	assert.Equal(t, "https://kpow.friends/callback", appConfig.Webhook.Url)

	// Inbox
	assert.Equal(t, projectRoot, appConfig.Inbox.Path)
	assert.Equal(t, "*/10 * * * *", appConfig.Inbox.Cron)
	assert.Equal(t, 4, appConfig.Inbox.BatchSize)
}
