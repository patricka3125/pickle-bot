package common

import (
	"context"
	"fmt"
	"time"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

func SendMessage(ctx context.Context, client *lark.Client, hostID, receiveId string,
	date time.Time, startTime, endTime, payment string, note string,
	withFee bool, courtSize int, courts string, roster Roster) error {
	content := fmt.Sprintf(`{"text":"<b>%s %s %s - %s</b>\nCourt #: %s\n`, date.Weekday().String(), date.Format("01/02"), startTime, endTime, courts) +
		fmt.Sprintf(`Host: <at user_id=\"%s\"></at>\n`, hostID)

	if withFee {
		var fee float32 = 20.0 * float32(courtSize) / float32(len(roster.Players))
		content += fmt.Sprintf(`Reserve fee: <b>$%0.02f</b> to %s\n`, fee, payment)
	}
	if note != "" {
		content += fmt.Sprintf(`Note: %s\n`, note)
	}

	content += `------------------------------------ \n` +
		fmt.Sprintf(`Players (%d/%d) :\n\n`, len(roster.Players), roster.Spots)

	for _, player := range roster.Players {
		if !withFee || !player.Paid {
			content += fmt.Sprintf(`%s. %s\n`, player.Number, player.Name)
		}
	}
	content += `"}`

	var receiveIdType string
	if receiveId[:2] == "oc" {
		receiveIdType = "chat_id"
	} else if receiveId[:2] == "ou" {
		receiveIdType = "open_id"
	} else {
		return fmt.Errorf("receiveID is bad format: %s", receiveId)
	}

	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(receiveIdType).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(receiveId).
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
