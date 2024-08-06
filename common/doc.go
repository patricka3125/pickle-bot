package common

import (
	"context"
	"fmt"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkdocx "github.com/larksuite/oapi-sdk-go/v3/service/docx/v1"
)

func ReadSignupDoc(ctx context.Context, client *lark.Client) error {
	readReq := larkdocx.NewListDocumentBlockReqBuilder().
		DocumentId(`FStaddwyloEuu9xgR1BuBFnTsLb`).
		PageSize(500).
		DocumentRevisionId(-1).
		Build()

	resp, err := client.Docx.DocumentBlock.List(ctx, readReq)
	if err != nil {
		return fmt.Errorf("failed to send request, err=%w", err)
	}
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return fmt.Errorf("failed to list document blocks: %+v", resp)
	}

	fmt.Println(larkcore.Prettify(resp))

	return nil
}
