package cmd

import (
	"context"
	"fmt"

	"github.com/patricka3125/pickle-bot/common"
	"github.com/spf13/cobra"
)

var (
	rosterCmd = &cobra.Command{
		Use:   "roster <signup_table_block_id>",
		Short: "List the final roster for a pickleball session.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			data, err := common.GetDocumentBlocks(ctx, client, cfg.PickleBall.DocumentID)
			if err != nil {
				return err
			}

			fmt.Println(args[0])
			roster, err := common.SignupRoster(args[0], data)
			if err != nil {
				return err
			}

			for _, player := range roster {
				fmt.Printf("player: %+v\n", player)
			}
			return nil
		},
	}
)
