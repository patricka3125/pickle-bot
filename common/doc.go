package common

import (
	"context"
	"fmt"
	"os"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkdocx "github.com/larksuite/oapi-sdk-go/v3/service/docx/v1"
)

type Player struct {
	Number string
	Name   string
	Paid   bool
}

type Roster struct {
	Players []Player
	Spots   int
}

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

func SignupRoster(blockID string, items []*larkdocx.Block) (*Roster, error) {
	if len(items) == 0 {
		return nil, fmt.Errorf("data is nil")
	}

	var (
		blockIter, cellSize int
		block               *larkdocx.Block
		result              *Roster
	)
	result = &Roster{
		Players: make([]Player, 0),
		Spots:   0,
	}

	// Get to table block with blockID
	for blockIter, block = range items {
		if block == nil {
			continue
		}

		// found signup table
		if *block.BlockId == blockID && *block.BlockType == 31 {
			result.Spots = *block.Table.Property.RowSize - 1
			cellSize = *block.Table.Property.ColumnSize * *block.Table.Property.RowSize
			break
		}
	}
	if cellSize == 0 {
		return nil, fmt.Errorf("table with blockID=%q not found", blockID)
	}

	var (
		headerCellCount, cellIter int
		curPlayer                 Player
	)

	for i := 1; cellIter < cellSize; i++ {
		if block = items[blockIter+i]; block == nil {
			return nil, fmt.Errorf("block is nil, blockID=%q, i=%d", blockID, i)
		}

		// Ignore table cell block
		if *block.BlockType == 32 {
			continue
		}

		// Skip parsing the table header
		if headerCellCount < 3 {
			headerCellCount++
			cellIter++
			continue
		}

		// recover from table cell parsing panic
		defer func() {
			if r := recover(); r != nil {
				fmt.Fprintf(os.Stderr, "failed to parse table cell\n")
				fmt.Fprintf(os.Stderr, "panic: %v\n", r)
			}
		}()

		switch cellIter % 3 {
		// No.
		case 0:
			curPlayer.Number = *block.Text.Elements[0].TextRun.Content
		// Name
		case 1:
			// Handle text element and mention element
			for _, element := range block.Text.Elements {
				if element.TextRun != nil {
					curPlayer.Name += *element.TextRun.Content
				}
				if element.MentionUser != nil {
					curPlayer.Name += fmt.Sprintf(`<at user_id=\"%s\"></at>`, *element.MentionUser.UserId)
				}
			}
		// Paid
		case 2:
			curPlayer.Paid = *block.Todo.Style.Done
			if curPlayer.Name != "" {
				result.Players = append(result.Players, curPlayer)
			}

			curPlayer = Player{}
		}

		cellIter++
	}

	return result, nil
}
