package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"code.byted.org/patrick.liao/pickle-bot/common"

	lark "github.com/larksuite/oapi-sdk-go/v3"
)

const (
	cfgFlag = "config"
)

var (
	client  *lark.Client
	rootCmd = &cobra.Command{
		Use:               "picklebot",
		Short:             "Picklebot is a CLI tool for Pickleball Lark bot tasks.",
		PersistentPreRunE: initConfig,
		RunE: func(cmd *cobra.Command, args []string) error {
			var cfg common.OpenAPIConfig
			if err := viper.Unmarshal(&cfg); err != nil {
				return err
			}

			client = lark.NewClient(cfg.AppID, cfg.AppKey)

			return nil
		},
	}
)

func Execute() error {
	homeDir, _ := os.UserHomeDir()

	rootCmd.PersistentFlags().String(cfgFlag, filepath.Join(homeDir+"/.picklebot/config.yaml"), "config file path")
	if err := rootCmd.Execute(); err != nil {
		return err
	}

	return nil
}

func initConfig(cmd *cobra.Command, args []string) error {
	path, err := cmd.Flags().GetString(cfgFlag)
	if err != nil {
		return fmt.Errorf("flag %q does not exist", cfgFlag)
	}

	viper.SetConfigType("yaml")
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("read config file %q failed: %w", path, err)
	}

	return nil
}
