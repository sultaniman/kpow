package cmd

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const envPrefix = "kpow"

var rootCmd = &cobra.Command{
	Use:   "kpow",
	Short: "KPow 💥 is a secure contact form.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err)
	}
}

func init() {
	cobra.OnInitialize(prepareCmds)
	rootCmd.AddCommand(startCmd)
}

func prepareCmds() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	logLevel, err := zerolog.ParseLevel(viper.GetString("log_level"))
	if err == nil {
		zerolog.SetGlobalLevel(logLevel)
	} else {
		log.Fatal().Err(err)
	}
}
