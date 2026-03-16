package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kpow",
	Short: "KPow 💥 – is a secure 🔒 contact form.",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err)
	}
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}
