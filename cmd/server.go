package cmd

import (
	"github.com/kortschak/utter"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/sultaniman/env"
	"github.com/sultaniman/kpow/server"
)

const envPrefix = "kpow"

var (
	port         int
	host         string
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
			config, err := server.GetConfig(configFile)
			if err != nil {
				return err
			}

			if port > 0 {
				config.Server.Port = port
			}

			if host != "" {
				config.Server.Host = host
			}

			if mailerDsn != "" {
				config.Mailer.DSN = mailerDsn
			}

			if advertiseKey {
				config.Key.Advertise = advertiseKey
			}

			if logLevel != "" {
				level, err := config.ParseLogLevel(logLevel)
				if err != nil {
					return err
				}

				zerolog.SetGlobalLevel(level)
			}

			if err := config.Validate(); err != nil {
				return err
			}

			if env.GetBool("debug") {
				utter.Dump(config)
			}

			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(startCmd)
	port := env.GetInt("KPOW_PORT")

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
