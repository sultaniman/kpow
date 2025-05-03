package cmd

import (
	"github.com/spf13/cobra"
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
	startCmd = &cobra.Command{
		Use: "start",
		Short: "Start server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
)

func init() {

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Path to config file")

	// Server options
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", server.Port, "Server port")
	rootCmd.PersistentFlags().StringVarP(&host, "host", "h", server.Host, "Server host")

	// Mailer options
	rootCmd.PersistentFlags().StringVarP(
		&fromEmail, "mailer-from", "f", "",
		"From email address",
	)
	rootCmd.PersistentFlags().StringVarP(
		&toEmail, "mailer-to", "t", "",
		"Recipient email (usually your email)",
	)
	rootCmd.PersistentFlags().StringVarP(
		&mailerDsn, "mailer", "m", "",
		"Mailer DSN, example: smtp://user:password@smtp.example.com:587",
	)

	// Encryption and key options
	rootCmd.PersistentFlags().StringVar(&password, "password", "", "Password for message encryption")
	rootCmd.PersistentFlags().StringVarP(&pubKeyPath, "pubkey", "k", "", "Path to public key file")
	rootCmd.PersistentFlags().BoolVarP(&advertiseKey, "advertise-key", "ka", false, "Advertise public key")
}
