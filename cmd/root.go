package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/patricka3125/pickle-bot/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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
			ctx := context.Background()
			var cfg common.Config
			if err := viper.Unmarshal(&cfg); err != nil {
				return err
			}

			client = lark.NewClient(cfg.OpenAPI.AppID, cfg.OpenAPI.AppKey)

			data, err := common.GetDocumentBlocks(ctx, client, cfg.PickleBall.DocumentID)
			if err != nil {
				return err
			}

			for _, block := range data {
				fmt.Println(*block.BlockId)
			}

			return nil
		},
	}
)

func Execute() error {
	homeDir, _ := os.UserHomeDir()
	rootCmd.PersistentFlags().String(cfgFlag, filepath.Join(homeDir, "/.picklebot/config.yaml"), "config file path")

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
