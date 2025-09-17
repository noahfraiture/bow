package main

import (
	"app/tui"
	"fmt"
	"log"
)

type panels struct {
	diffFrom  commitPanel
	diffOn    commitPanel
	diffs     diffPanel
	updateMsg messagePanel
	createMsg messagePanel
}

func createApp() (*tui.App, error) {

	commits, err := getCommits()
	if err != nil {
		return nil, fmt.Errorf("failed to create app: %w", err)
	}
	diffs, err := getDiff()
	if err != nil {
		return nil, fmt.Errorf("failed to create app: %w", err)
	}

	panels := panels{
		diffFrom:  newCommitPanel("Diff from", commits),
		diffOn:    newCommitPanel("Diff from", commits),
		diffs:     newDiffPanel("Diff to update", diffs),
		updateMsg: newMessagePanel("Message"),
		createMsg: newMessagePanel("Message"),
	}

	defaultLayout := &tui.HorizontalSplit{
		Left: &tui.VerticalSplit{
			Top:    &tui.PanelNode{Panel: &panels.diffFrom},
			Bottom: &tui.PanelNode{Panel: &panels.diffOn},
		},
		Right: &tui.VerticalSplit{
			Top:    &tui.PanelNode{Panel: &panels.diffs},
			Bottom: &tui.PanelNode{Panel: &panels.updateMsg},
		},
	}

	handler := &handler{
		diffFromCommit: panels.diffFrom.commit,
		diffOnCommit:   panels.diffOn.commit,
		diffToUpdate:   panels.diffs.diff,
		updateMsg:      panels.updateMsg.msg,
		createMsg:      panels.createMsg.msg,
		panels:         panels,
		activeCommand:  Update,
		rightPanel:     &defaultLayout.Right,
	}

	app := tui.NewApp(defaultLayout, handler)

	return app, nil
}

func main() {
	app, err := createApp()
	if err != nil {
		log.Fatalf("failed to start application: %v", err)
	}

	app.Run()
}
