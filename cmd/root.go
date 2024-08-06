package cmd

import "github.com/spf13/cobra"

const (
	cfgFlag = "config"
)

var (
	rootCmd = &cobra.Command{
		Use:               "picklebot",
		Short:             "Picklebot is a CLI tool for Pickleball Lark bot tasks.",
		PersistentPreRunE: initconfig,
	}
)

func Execute() error {
	rootCmd.PersistentFlags().String(cfgFlag, "./config.yaml", "config file path")

	if err := rootCmd.Execute(); err != nil {
		return err
	}

	return nil
}
