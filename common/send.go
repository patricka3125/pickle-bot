package common

import (
	"context"
	"fmt"
	"time"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

// Pickleball group chat ID: oc_31637de07de089a77d95c62e5e815bc7
// ou_d68e3cb2a2ccae151f8bc59d85ddb0c8
func SendMessage(ctx context.Context, client *lark.Client, hostID, receiveID string,
	date time.Time, startTime, endTime, payment string, courtSize int,
	roster Roster) error {
	var fee float32 = 26.0 * float32(courtSize) / float32(len(roster.Players))

	content := fmt.Sprintf(`{"text":"<b>%s %s %s - %s</b>\nCourt #: 5-8\n`, date.Weekday().String(), date.Format("01/02"), startTime, endTime) +
		fmt.Sprintf(`Host: <at user_id=\"%s\"></at>\n`, hostID) +
		fmt.Sprintf(`Reserve fee: <b>$%0.02f</b> to %s\n\n`, fee, payment) +
		`------------------------------------ \n` +
		fmt.Sprintf(`Players (%d/%d) :\n\n`, len(roster.Players), roster.Spots)

	for _, player := range roster.Players {
		content += fmt.Sprintf(`%s. %s<at user_id=\"%s\"></at>\n`, player.Number, player.InviteName, player.Ouid)
	}
	content += `"}`

	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType("open_id").
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(receiveID).
			MsgType("text").
			Content(content).
			Build()).
		Build()

	resp, err := client.Im.Message.Create(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return fmt.Errorf("send message returned error response: %+v", resp)
	}

	return nil
}
