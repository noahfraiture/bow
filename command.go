package main

import (
	"app/tui"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	// FIX : patch might not be the best as phabricator now ignores the base and has no commit. Maybe checkout is better, I have "context not available" when look at the diff and the test do not run correctly
	patch, err := cp.on.Patch(cp.from.Commit)
	if err != nil {
		return fmt.Errorf("failed to get patch from commits")
	}
	patcherReader := strings.NewReader(patch.String())
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	logDir := filepath.Join(home, ".cache", "bow")
	err = os.MkdirAll(logDir, 0755)
	if err != nil {
		return err
	}
	logPath := filepath.Join(logDir, "bow.log")
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = fmt.Fprintf(file, "Running command: arc diff --raw --update %s --message %s\n", cp.diff.id, *cp.msg)
	if err != nil {
		return fmt.Errorf("failed to create command: %w", err)
	}
	switch cp.Items[cp.Selected] {
	case Create:
		// TODO : message !
		// template, could change based on the selected command and we would confirm in this panel ? but key enter new line or confirm ?
	case Update:
		// Log the command
		cmd := exec.Command("arc", "diff", "--raw", "--update", cp.diff.id, "--message", *cp.msg)
		cmd.Stdin = patcherReader
		// TODO : send to actual stdout after close or in panel ? could be the best
		cmd.Stdout = file
		cmd.Stderr = file
		return cmd.Run()
	}
	return errors.New("invalid command")
}
