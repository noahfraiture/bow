package main

import (
	"app/tui"
	"fmt"
)

type command string

const (
	Update command = "Update"
	Create command = "Create"
)

type commandPanel struct {
	*tui.ListPanel[command]
	from *commit
	to   *commit
	diff *diff
}

func (c *commandPanel) Update(msg tui.InputMessage) bool {
	if msg.IsKey(tui.KeyEnter) {
		panic(fmt.Sprintln(c.from, c.to, c.diff))
	}
	return c.ListPanel.Update(msg)
}

func NewCmdPanel(name string, from, to *commit, diff *diff) commandPanel {
	return commandPanel{
		ListPanel: &tui.ListPanel[command]{
			PanelBase: tui.PanelBase{
				Title:  name,
				Border: true,
			},
			Items: []command{Update, Create},
		},
		from: from,
		to:   to,
		diff: diff,
	}
}
