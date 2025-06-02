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

const envPrefix = "KPOW_"

var (
	port         int
	host         string = "0.0.0.0"
	configFile   string
	password     string
	pubKeyPath   string
	mailerDsn    string
	fromEmail    string
	toEmail      string
	logLevel     string
	advertiseKey bool
	startCmd     = &cobra.Command{
		Use:   "start",
		Short: "Start server",
		RunE: func(cmd *cobra.Command, args []string) error {
			config.WarnAboutPassphrase()
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

			if logLevel != "" {
				level, err := appConfig.ParseLogLevel(logLevel)
				if err != nil {
					return err
				}

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
	env.SetEnvPrefix(envPrefix)
	port := env.GetInt("PORT")

	// viper.SetEnvPrefix(envPrefix)
	startCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Path to config file")

	// Server options
	startCmd.PersistentFlags().IntVar(&port, "port", -1, "Server port")
	startCmd.PersistentFlags().StringVar(&host, "host", "", "Server host")

	// Mailer options
	startCmd.PersistentFlags().StringVarP(
		&fromEmail, "mailer-from", "f", "",
		"From email address",
	)

	startCmd.PersistentFlags().StringVarP(
		&toEmail, "mailer-to", "t", "",
		"Recipient email (usually your email)",
	)

	startCmd.PersistentFlags().StringVarP(
		&mailerDsn, "mailer", "m", "",
		"Mailer DSN, example: smtp://user:password@smtp.example.com:587",
	)

	// Encryption and key options
	startCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "Password for message encryption")
	startCmd.PersistentFlags().StringVarP(&pubKeyPath, "pubkey", "k", "", "Path to public key file")
	startCmd.PersistentFlags().BoolVarP(&advertiseKey, "advertise-key", "s", false, "Advertise public key")
	startCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "WARN", "Log level")
}
