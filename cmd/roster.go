package cmd

import (
	"context"
	"fmt"
	"time"

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

			roster, err := common.SignupRoster(args[0], data)
			if err != nil {
				return err
			}

			courts := prompt("Please enter courts")
			dateIn := prompt("Please enter date (MM/DD/YYYY)")
			startTime := prompt("Please enter start time (HH:MM)")
			endTime := prompt("Please enter end time (HH:MM)")
			host := prompt("Please enter host uid (ou_xxx)")
			payment := prompt("Please enter payment link")

			date, err := time.Parse("01/02/2006", dateIn)
			if err != nil {
				return fmt.Errorf("invalid date: %w", err)
			}

			fmt.Println("\n------------------------------------------------------")
			fmt.Printf("\nCourts: %s\nPlayers: %d/%d\nHost: %s\nPayment link: %s\n%s, %s %s - %s\n\n",
				courts, len(roster.Players), roster.Spots, host, payment,
				date.Weekday().String(), dateIn, startTime, endTime)
			for _, player := range roster.Players {
				fmt.Printf("player: %+v\n", player)
			}

			proceed := prompt("\nProceed with sending message? (y/n)")
			if proceed != "y" {
				fmt.Println("Aborting...")
				return nil
			}

			return nil
		},
	}
)

func prompt(msg string) string {
	var input string
	fmt.Printf("%s : ", msg)
	fmt.Scanf("%s", &input)

	return input
}
