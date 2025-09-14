package main

import (
	"app/tui"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type command string

const (
	Update command = "Update"
	Create command = "Create"
)

type commandPanel struct {
	*tui.ListPanel[command]
	from *commit
	on   *commit
	diff *diff
	msg  *string
}

func (cp *commandPanel) Update(msg tui.InputMessage) bool {
	if msg.IsKey(tui.KeyEnter) {
		cp.Stop()
		err := cp.runCmd()
		if err != nil {
			panic(err)
		}
	}
	return cp.ListPanel.Update(msg)
}

func NewCmdPanel(name string, from, on *commit, diff *diff, msg *string) commandPanel {
	return commandPanel{
		ListPanel: &tui.ListPanel[command]{
			PanelBase: tui.PanelBase{
				Title:  name,
				Border: true,
			},
			Items: []command{Update, Create},
		},
		from: from,
		on:   on,
		diff: diff,
		msg:  msg,
	}
}

func (cp *commandPanel) runCmd() error {
	patch, err := cp.from.Patch(cp.on.Commit)
	if err != nil {
		return fmt.Errorf("failed to get patch from commits")
	}
	patcherReader := strings.NewReader(patch.String())
	switch cp.Items[cp.Selected] {
	case Create:
	case Update:
		// FIX : does nothing and no log
		// Could write logs stdout and stderr in a file
		cmd := exec.Command("arc", "diff", "--raw", "--update", cp.diff.id, "--message", *cp.msg)
		cmd.Stdin = patcherReader
		return cmd.Run()
	}
	return errors.New("invalid command")
}
