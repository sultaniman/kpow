package cmd

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const envPrefix = "kpow"

var rootCmd = &cobra.Command{
	Use:   "kpow",
	Short: "KPow 💥 is a minimal & privacy focused contact form.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err)
	}
	spew.Dump(config)
}

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(startCmd)
}

func initConfig() {
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()

	if configFile != "" {
		viper.SetConfigType("toml")
		viper.SetConfigFile(configFile)
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal().Err(err)
	}

	if err := viper.Unmarshal(config); err != nil {
		log.Fatal().Err(err)
	}
}
