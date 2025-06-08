package cmd

import (
	"errors"
	"fmt"

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
	rpm                      int
	numBurst                 int
	rateLimitCooldownSeconds int
	// pubkey
	pubKeyPath   string
	keyKind      string
	advertiseKey bool
	// mailer
	mailerDsn string
	fromEmail string
	toEmail   string
	// inbox
	inboxPath      string
	inboxCron      string
	inboxBatchSize int
	startCmd       = &cobra.Command{
		Use:   "start",
		Short: "Start server",
		RunE: func(cmd *cobra.Command, args []string) error {
			appConfig, err := config.GetConfig(configFile)
			if err != nil {
				return err
			}

			if port > 0 {
				appConfig.Server.Port = port
			}

			if host != "" {
				appConfig.Server.Host = host
			}

			if mailerDsn != "" {
				appConfig.Mailer.DSN = mailerDsn
			}

			if advertiseKey {
				appConfig.Key.Advertise = advertiseKey
			}

			if level, err := appConfig.ParseLogLevel(logLevel); err != nil {
				return err
			} else {
				zerolog.SetGlobalLevel(level)
			}

			if errorList := appConfig.Validate(); len(errorList) > 0 {
				server.LogErrors(errorList)
				return errors.New("configuration error")
			}

			if env.GetBool("DEBUG") {
				godump.Dump(appConfig)
			}

			app, err := server.CreateServer(appConfig)
			if err != nil {
				return err
			}

			scheduler := cron.NewScheduler(appConfig.Inbox.Cron)
			scheduler.AddFunc(appConfig.Inbox.Cron, cron.InboxCleaner(appConfig.Inbox.Path))
			scheduler.Start()
			defer scheduler.Stop()

			err = app.Start(fmt.Sprintf("%s:%d", appConfig.Server.Host, appConfig.Server.Port))

			log.Fatal().Err(err)

			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(startCmd)
	// Read system $PORT env value and
	// use it below if KPOW_PORT is not set
	systemPort := env.GetInt("PORT")
	env.SetEnvPrefix(envPrefix)

	port := env.GetInt("PORT")
	if port == 0 {
		port = systemPort
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
		&rpm, "limiter-rpm", 0,
		"Rate limiter, requests per minute",
	)

	startCmd.PersistentFlags().IntVar(
		&numBurst, "limiter-burst", -1,
		"Rate limiter burst requests",
	)

	startCmd.PersistentFlags().IntVar(
		&rateLimitCooldownSeconds, "limiter-cooldown", -1,
		"Rate limiter cooldown time",
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
		&mailerDsn, "mailer", "",
		"Mailer DSN, example: smtp://user:password@smtp.example.com:587",
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
		&inboxCron, "inbox-cron", defaultCronSpec,
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
		&messageSize, "size", defaultMessageSizeBytes,
		"Size of the message in bytes",
	)
}
