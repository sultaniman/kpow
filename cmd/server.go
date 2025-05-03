package cmd

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sultaniman/kpow/server"
)

var (
	port         int
	host         string
	configFile   string
	password     string
	pubKeyPath   string
	mailerDsn    string
	fromEmail    string
	toEmail      string
	advertiseKey bool
	config       = server.NewConfig()
	startCmd     = &cobra.Command{
		Use:   "start",
		Short: "Start server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if viper.GetBool("debug") {
				spew.Dump(config)
			}

			return nil
		},
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	startCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Path to config file")

	// Server options
	startCmd.PersistentFlags().IntVar(&port, "port", server.Port, "Server port")
	startCmd.PersistentFlags().StringVar(&host, "host", server.Host, "Server host")

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
}

func initConfig() {
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()
	viper.AllowEmptyEnv(true)

	if configFile == "" {
		return
	}

	viper.SetConfigType("toml")
	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal().Err(err)
	}

	if err := viper.Unmarshal(config); err != nil {
		log.Fatal().Err(err)
	}
}
