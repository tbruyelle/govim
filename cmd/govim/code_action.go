package main

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/govim/govim"
	"github.com/govim/govim/cmd/govim/internal/golang_org_x_tools/lsp/protocol"
	"github.com/govim/govim/cmd/govim/internal/types"
)

func (v *vimstate) runGoTest(flags govim.CommandFlags, args ...string) error {
	b, _, err := v.bufCursorPos()
	if err != nil {
		return fmt.Errorf("failed to get cursor position: %v", err)
	}
	start, end, err := v.rangeFromFlags(b, flags)
	if err != nil {
		return err
	}

	ca, err := v.server.CodeAction(context.Background(), &protocol.CodeActionParams{
		TextDocument: b.ToTextDocumentIdentifier(),
		Range: protocol.Range{
			Start: start.ToPosition(),
			End:   end.ToPosition(),
		},
		Context: protocol.CodeActionContext{
			Only: []protocol.CodeActionKind{protocol.GoTest},
		},
		WorkDoneProgressParams: protocol.WorkDoneProgressParams{},
	})
	if err != nil {
		return err
	}
	if len(ca) > 1 {
		return fmt.Errorf("got %d CodeActions, expected no more than 1", len(ca))
	}

	c := ca[0]
	var token protocol.ProgressToken
	if c := v.config.ExperimentalProgressPopups; c != nil && *c {
		token = fmt.Sprintf("govim%d", rand.Uint64())
		if _, ok := v.progressPopups[token]; ok {
			return fmt.Errorf("failed to init progress, duplicate token")
		}
		v.progressPopups[token] = &types.ProgressPopup{}
	}

	_, err = v.server.ExecuteCommand(context.Background(), &protocol.ExecuteCommandParams{
		Command:   c.Command.Command,
		Arguments: c.Command.Arguments,
		WorkDoneProgressParams: protocol.WorkDoneProgressParams{
			WorkDoneToken: token,
		},
	})

	return err
}
