package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/sultaniman/env"
	"github.com/sultaniman/kpow/config"
)

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify configuration file and exit",
	RunE: func(cmd *cobra.Command, args []string) error {
		env.SetEnvPrefix(envPrefix)
		cfg, err := config.GetConfig(configFile)
		if err != nil {
			return err
		}
		if errs := cfg.Validate(); len(errs) > 0 {
			// Use helper to print all validation errors
			config.LogErrors(errs)
			return errors.New("configuration error")
		}
		fmt.Println("all good")
		return nil
	},
}

func init() {
	verifyCmd.PersistentFlags().StringVarP(
		&configFile, "config", "c", "",
		"Path to config file",
	)
}
