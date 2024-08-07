package common

import (
	"context"
	"fmt"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkdocx "github.com/larksuite/oapi-sdk-go/v3/service/docx/v1"
)

type Player struct {
	num  int
	name string
	paid bool
}

type Roster []Player

func GetDocumentBlocks(ctx context.Context, client *lark.Client, docID string) ([]*larkdocx.Block, error) {
	var (
		pageToken string
		result    []*larkdocx.Block = make([]*larkdocx.Block, 0)
	)

	reqBuilder := larkdocx.NewListDocumentBlockReqBuilder().
		DocumentId(docID).
		PageSize(500).
		DocumentRevisionId(-1)

	for {
		readReq := reqBuilder.PageToken(pageToken).Build()
		resp, err := client.Docx.DocumentBlock.List(ctx, readReq)
		if err != nil {
			return nil, fmt.Errorf("failed to send request, err=%w", err)
		}
		if !resp.Success() {
			fmt.Println(resp.Code, resp.Msg, resp.RequestId())
			return nil, fmt.Errorf("failed to list document blocks: %+v", resp)
		}

		result = append(result, resp.Data.Items...)
		if resp.Data.HasMore != nil && *resp.Data.HasMore {
			if resp.Data.PageToken == nil {
				return nil, fmt.Errorf("invalid page token in response")
			}
			pageToken = *resp.Data.PageToken
			continue
		}

		break
	}

	return result, nil
}

func ParseSignupTable(ctx context.Context, client *lark.Client, docID, blockID string) error {
	readReq := larkdocx.NewGetDocumentBlockReqBuilder().
		DocumentId(docID).
		BlockId(blockID).
		DocumentRevisionId(-1).
		Build()

	resp, err := client.Docx.DocumentBlock.Get(ctx, readReq)
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

func getSignupRoster(blockID string, items []*larkdocx.Block) (Roster, error) {
	if len(items) == 0 {
		return nil, fmt.Errorf("data is nil")
	}

	var (
		i, columnSize, cellSize int
		block                   *larkdocx.Block
	)

	// Get to table block with blockID
	for i, block = range items {
		if block == nil {
			continue
		}

		// found signup table
		if *block.BlockId == blockID || *block.BlockType == 32 {
			columnSize = *block.Table.Property.ColumnSize
			cellSize = *block.Table.Property.ColumnSize * *block.Table.Property.RowSize
			break
		}
	}
	if columnSize == 0 || cellSize == 0 {
		return nil, fmt.Errorf("table with blockID=%q not found")
	}

	for i < cellSize {
		block = items[i]
		if block == nil {
			return nil, fmt.Errorf("block is nil, blockID=%q, i=%d")
		}

		// Ignore table cell block
		if *block.BlockType == 32 {
			continue
		}
	}

	return nil, nil
}
