package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/goforj/godump"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/sultaniman/env"
	"github.com/sultaniman/kpow/config"
	"github.com/sultaniman/kpow/server"
	"github.com/sultaniman/kpow/server/cron"
)

const (
	envPrefix               = "KPOW_"
	defaultMessageSizeBytes = 240
	defaultCronSpec         = "*/5 * * * *"
	defaultBatchSize        = 5
)

var (
	// server config
	port         int
	host         string = "0.0.0.0"
	configFile   string
	logLevel     string
	customBanner string
	hideLogo     bool
	messageSize  int
	// rate limiter
	limiterRPM             int
	limiterBurst           int
	limiterCooldownSeconds int
	// pubkey
	pubKeyPath   string
	keyKind      string
	advertiseKey bool
	// mailer
	mailerDsn string
	fromEmail string
	toEmail   string
	// webhook
	webhookUrl string
	// inbox
	inboxPath      string
	inboxCron      string
	inboxBatchSize int
	startCmd       = &cobra.Command{
		Use:   "start",
		Short: "Start server",
		RunE: func(cmd *cobra.Command, args []string) error {
			appConfig, err := getConfig()
			if err != nil {
				return err
			}

			if env.GetBool("DEBUG") {
				godump.Dump(appConfig)
			}

			app, err := server.CreateServer(appConfig)
			if err != nil {
				return err
			}

			scheduler := cron.NewScheduler(appConfig.Inbox.Cron)
			_, err = scheduler.AddFunc(appConfig.Inbox.Cron, cron.InboxCleaner(appConfig.Inbox.Path))
			if err != nil {
				return err
			}

			scheduler.Start()
			defer scheduler.Stop()

			err = app.Start(fmt.Sprintf("%s:%d", appConfig.Server.Host, appConfig.Server.Port))

			log.Fatal().Err(err)

			return nil
		},
	}
)

func getConfig() (*config.Config, error) {
	appConfig, err := config.GetConfig(configFile)
	if err != nil {
		return nil, err
	}

	// server
	if port > 0 {
		appConfig.Server.Port = port
	}

	if host != "" {
		appConfig.Server.Host = host
	}

	if messageSize > 0 {
		appConfig.Server.MessageSize = messageSize
	}

	if hideLogo {
		appConfig.Server.HideLogo = hideLogo
	}

	if customBanner != "" {
		bannerBytes, err := os.ReadFile(customBanner)
		if err != nil {
			return nil, err
		}

		appConfig.Server.CustomBanner = string(bannerBytes)
	}

	if level, err := appConfig.ParseLogLevel(logLevel); err != nil {
		return nil, err
	} else {
		zerolog.SetGlobalLevel(level)
	}

	// key
	if keyKind != "" {
		appConfig.Key.Kind = config.KeyKind(keyKind)
	}

	if pubKeyPath != "" {
		appConfig.Key.Path = pubKeyPath
	}

	if advertiseKey {
		appConfig.Key.Advertise = advertiseKey
	}

	// mailer
	if mailerDsn != "" {
		appConfig.Mailer.DSN = mailerDsn
	}

	if toEmail != "" {
		appConfig.Mailer.To = toEmail
	}

	if fromEmail != "" {
		appConfig.Mailer.From = fromEmail
	}

	// webhook
	if webhookUrl != "" {
		appConfig.Webhook.Url = webhookUrl
	}

	// inbox
	if inboxPath != "" {
		appConfig.Inbox.Path = inboxPath
	}

	if inboxCron != "" {
		appConfig.Inbox.Cron = inboxCron
	}

	if inboxBatchSize > 1 {
		appConfig.Inbox.BatchSize = inboxBatchSize
	}

	if appConfig.RateLimiter == nil {
		appConfig.RateLimiter = &config.RateLimiter{}
	}

	// rate limiter
	if limiterRPM > 0 {
		appConfig.RateLimiter.RPM = limiterRPM
	}

	if limiterBurst > 0 {
		appConfig.RateLimiter.Burst = limiterBurst
	}

	if limiterCooldownSeconds > 0 {
		appConfig.RateLimiter.CooldownSeconds = limiterCooldownSeconds
	}

	if errorList := appConfig.Validate(); len(errorList) > 0 {
		server.LogErrors(errorList)
		return nil, errors.New("configuration error")
	}

	return appConfig, nil
}

func init() {
	rootCmd.AddCommand(startCmd)
	// Read system $PORT env value and
	// use it below if KPOW_PORT is not set
	port = env.GetInt("PORT")
	env.SetEnvPrefix(envPrefix)
	if serverPort := env.GetInt("PORT"); serverPort > 0 {
		port = serverPort
	}

	// viper.SetEnvPrefix(envPrefix)
	startCmd.PersistentFlags().StringVarP(
		&configFile, "config", "c", "",
		"Path to config file",
	)

	// Server options
	startCmd.PersistentFlags().IntVar(
		&port, "port", -1,
		"Server port",
	)

	startCmd.PersistentFlags().StringVar(
		&host, "host", "",
		"Server host",
	)

	startCmd.PersistentFlags().IntVar(
		&limiterRPM, "limiter-rpm", 0,
		"Rate limiter, requests per minute",
	)

	startCmd.PersistentFlags().IntVar(
		&limiterBurst, "limiter-burst", -1,
		"Rate limiter burst requests",
	)

	startCmd.PersistentFlags().IntVar(
		&limiterCooldownSeconds, "limiter-cooldown", -1,
		"Rate limiter cooldown seconds",
	)

	// Mailer options
	startCmd.PersistentFlags().StringVar(
		&fromEmail, "mailer-from", "",
		"From email address",
	)

	startCmd.PersistentFlags().StringVar(
		&toEmail, "mailer-to", "",
		"Recipient email (usually your email)",
	)

	startCmd.PersistentFlags().StringVar(
		&mailerDsn, "mailer-dsn", "",
		"Mailer DSN, example: smtp://user:password@smtp.example.com:587",
	)

	// Webhook URL
	startCmd.PersistentFlags().StringVar(
		&webhookUrl, "webhook-url", "",
		"Webhook URL (must be a https url)",
	)

	// Key options
	startCmd.PersistentFlags().StringVar(
		&pubKeyPath, "pubkey", "",
		"Path to public key file",
	)

	startCmd.PersistentFlags().StringVar(
		&keyKind, "key-kind", "",
		"Type of public key one of RSA, AGE, PGP",
	)

	startCmd.PersistentFlags().BoolVar(
		&advertiseKey, "advertise-key", false,
		"Advertise public key",
	)

	// Inbox options
	startCmd.PersistentFlags().StringVar(
		&inboxPath, "inbox-path", "",
		"Path to message inbox",
	)

	startCmd.PersistentFlags().StringVar(
		&inboxCron, "inbox-cron", "",
		"Schedule of inbox cleaner",
	)

	startCmd.PersistentFlags().IntVar(
		&inboxBatchSize, "batch-size", defaultBatchSize,
		"Schedule of inbox cleaner",
	)

	// Server options
	startCmd.PersistentFlags().StringVar(
		&logLevel, "log-level", "INFO",
		"Log level",
	)

	startCmd.PersistentFlags().StringVar(
		&customBanner, "banner", "",
		"Custom banner above the form (path to html file)",
	)

	startCmd.PersistentFlags().BoolVar(
		&hideLogo, "hide-logo", false,
		"Hide logo above the form",
	)

	startCmd.PersistentFlags().IntVar(
		&messageSize, "message-size", defaultMessageSizeBytes,
		"Size of the message in bytes",
	)
}
