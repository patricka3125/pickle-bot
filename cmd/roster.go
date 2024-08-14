package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/patricka3125/pickle-bot/common"
	"github.com/spf13/cobra"
)

var (
	withFeeFlag bool
	rosterCmd   = &cobra.Command{
		Use:   "roster <signup_table_block_id>",
		Short: "List the final roster for a pickleball session.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			withFee, err := cmd.Flags().GetBool(feeFlag)
			if err != nil {
				return err
			}

			data, err := common.GetDocumentBlocks(ctx, client, cfg.PickleBall.DocumentID)
			if err != nil {
				return err
			}

			roster, err := common.SignupRoster(args[0], data)
			if err != nil {
				return err
			}

			for _, player := range roster.Players {
				fmt.Printf("%+v\n", player)
			}
			fmt.Println()

			courtSizeIn := prompt("Please enter court size")
			courtSize, err := strconv.Atoi(courtSizeIn)
			if courtSize <= 0 {
				return fmt.Errorf("invalid court size: courtSize=%s", courtSizeIn)
			}

			courts := prompt("Please enter courts")
			dateIn := prompt("Please enter date (MM/DD/YYYY)")
			startTime := prompt("Please enter start time (HH:MM)")
			endTime := prompt("Please enter end time (HH:MM)")
			hostID := prompt("Please enter host uid (ou_xxx)")

			var note string
			if addNote, _ := cmd.Flags().GetBool(noteFlag); addNote {
				note = prompt("Please enter additional note")
			}

			var payment string
			if withFee {
				payment = prompt("Please enter payment link")
			}

			date, err := time.Parse("01/02/2006", dateIn)
			if err != nil {
				return fmt.Errorf("invalid date: %w", err)
			}

			fmt.Println("\n------------------------------------------------------")
			fmt.Printf("\nCourts: %s\nPlayers: %d/%d\nHost: %s\nPayment link: %s\n%s, %s %s - %s\n\n",
				courts, len(roster.Players), roster.Spots, hostID, payment,
				date.Weekday().String(), dateIn, startTime, endTime)

			proceed := prompt("\nProceed with sending message? (y/n)")
			if proceed != "y" {
				fmt.Println("Aborting...")
				return nil
			}

			if err := common.SendMessage(ctx, client,
				hostID, cfg.OpenAPI.ReceiveID,
				date, startTime, endTime, payment, note,
				withFee, courtSize, courts, *roster); err != nil {
				return err
			}

			return nil
		},
	}
)

func prompt(msg string) string {
	fmt.Printf("%s: ", msg)
	in := bufio.NewReader(os.Stdin)
	input, _ := in.ReadString('\n')

	return strings.TrimSpace(input)
}
