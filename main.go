package main

import (
	"app/tui"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

type panels struct {
	diffFrom  commitPanel
	diffOn    commitPanel
	diffs     diffPanel
	updateMsg messagePanel
	createMsg messagePanel
}

func createApp() (*tui.App, *handler, error) {

	commits, err := getCommits()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create app: %w", err)
	}
	diffs, err := getDiff()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create app: %w", err)
	}

	panels := panels{
		diffFrom:  newCommitPanel("Diff from", commits),
		diffOn:    newCommitPanel("Diff from", commits),
		diffs:     newDiffPanel("Diff to update", diffs),
		updateMsg: newMessagePanel("Message"),
		createMsg: newMessagePanel("Message"),
	}

	defaultLayout := &tui.HorizontalSplit{
		Panels: []tui.Layout{
			&tui.VerticalSplit{
				Panels: []tui.Layout{
					&tui.PanelNode{Panel: &panels.diffFrom, Weight: 1},
					&tui.PanelNode{Panel: &panels.diffOn, Weight: 1},
				},
				Weight: 2,
			},
			&tui.VerticalSplit{
				Panels: []tui.Layout{
					&tui.PanelNode{Panel: &panels.diffs, Weight: 1},
					&tui.PanelNode{Panel: &panels.updateMsg, Weight: 0},
				},
				Weight: 1,
			},
		},
		Weight: 1,
	}

	handler := &handler{
		diffFromCommit: panels.diffFrom.commit,
		diffOnCommit:   panels.diffOn.commit,
		diffToUpdate:   panels.diffs.diff,
		updateMsg:      panels.updateMsg.msg,
		createMsg:      panels.createMsg.msg,
		panels:         panels,
		activeCommand:  Update,
		rightPanel:     &defaultLayout.Panels[1],
	}

	app := tui.NewApp(defaultLayout, handler)

	return app, handler, nil
}

func main() {
	// Setup logging
	cacheDir := filepath.Join(os.Getenv("HOME"), ".cache", "bow")
	_ = os.MkdirAll(cacheDir, 0755)
	logFile := filepath.Join(cacheDir, "bow.log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open log file: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = file.Close() }()

	handler := slog.NewTextHandler(file, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	app, h, err := createApp()
	if err != nil {
		slog.Error("failed to start application", "error", err)
		os.Exit(1)
	}

	app.Run()

	if h.lastOutput != "" {
		fmt.Print(h.lastOutput)
	}
}
