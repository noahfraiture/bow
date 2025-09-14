package main

import (
	"app/tui"
	"fmt"
	"log"
)

func createApp() (*tui.App, error) {
	commits, err := getCommits()
	if err != nil {
		return nil, fmt.Errorf("failed to create app: %w", err)
	}
	diffFrom := newCommitPanel("Diff from", commits)
	diffOn := newCommitPanel("Diff on", commits)
	diffOn.Selected = 1

	diffs, err := getDiff()
	if err != nil {
		return nil, fmt.Errorf("failed to create app: %w", err)
	}
	diffToUpdate := newDiffPanel("Diff to update", diffs)

	var msg string
	msgToUpdate := newMessagePanel("Message", &msg)

	cmd := NewCmdPanel(
		"Command",
		diffFrom.commit,
		diffOn.commit,
		diffToUpdate.diff,
		msgToUpdate.msg,
	)

	layout := &tui.HorizontalSplit{
		Left: &tui.VerticalSplit{
			Top:    &tui.PanelNode{Panel: &diffFrom},
			Bottom: &tui.PanelNode{Panel: &diffOn},
		},
		Right: &tui.VerticalSplit{
			Top: &tui.PanelNode{Panel: &diffToUpdate},
			Bottom: &tui.VerticalSplit{
				Top:    &tui.PanelNode{Panel: &msgToUpdate},
				Bottom: &tui.PanelNode{Panel: &cmd},
			},
		},
	}
	app := tui.NewApp(layout)

	return app, nil
}

func main() {
	app, err := createApp()
	if err != nil {
		log.Fatalf("failed to start application: %v", err)
	}

	app.Run()
}
