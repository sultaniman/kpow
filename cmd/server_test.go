package cmd

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func resetFlags() {
	port = -1
	host = "0.0.0.0"
	configFile = ""
	logLevel = ""
	customBanner = ""
	hideLogo = false
	messageSize = 0
	limiterRPM = 0
	limiterBurst = 0
	limiterCooldownSeconds = 0
	pubKeyPath = ""
	keyKind = ""
	advertiseKey = false
	mailerDsn = ""
	fromEmail = ""
	toEmail = ""
	maxRetries = 0
	webhookUrl = ""
	inboxPath = ""
	inboxCron = ""
}

func TestConfigurationOverrides(t *testing.T) {
	resetFlags()
	pwd, _ := os.Getwd()
	projectRoot := filepath.Dir(pwd)
	keyPath := path.Join(projectRoot, "server/enc/testkeys/pubkey.gpg")
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
	assert.Contains(t, appConfig.Server.CustomBanner, "This a banner")
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
}

func TestBannerFlag(t *testing.T) {
	resetFlags()
	pwd, _ := os.Getwd()
	projectRoot := filepath.Dir(pwd)
	configPath := path.Join(projectRoot, "config.toml")
	bannerPath := path.Join(projectRoot, "banner.html")
	err := startCmd.ParseFlags([]string{"--config=" + configPath, "--banner=" + bannerPath})
	assert.NoError(t, err)
	appConfig, err := getConfig()
	assert.NoError(t, err)
	assert.Contains(t, appConfig.Server.CustomBanner, "This a banner")
}

func TestHideLogoFlag(t *testing.T) {
	resetFlags()
	pwd, _ := os.Getwd()
	projectRoot := filepath.Dir(pwd)
	configPath := path.Join(projectRoot, "config.toml")
	err := startCmd.ParseFlags([]string{"--config=" + configPath, "--hide-logo=true"})
	assert.NoError(t, err)
	appConfig, err := getConfig()
	assert.NoError(t, err)
	assert.Equal(t, true, appConfig.Server.HideLogo)
}
