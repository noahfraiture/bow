package main

import (
	"app/tui"
)

func main() {
	files := &tui.ListPanel{
		PanelBase: tui.PanelBase{Title: "Files", Border: true},
		Items:     []string{"main.go", "config.json", "README.md", "Dockerfile", "go.mod", "utils.go"},
		Selected:  0,
	}
	commands := &tui.ListPanel{
		PanelBase: tui.PanelBase{Title: "Commands", Border: true},
		Items:     []string{"git status", "git add .", "git commit", "git push", "go run .", "go test"},
		Selected:  0,
	}
	input := &tui.TextPanel{
		PanelBase: tui.PanelBase{Title: "Input", Border: true},
		Text:      []rune{},
		Cursor:    0,
	}
	input2 := &tui.TextPanel{
		PanelBase: tui.PanelBase{Title: "Input", Border: true},
		Text:      []rune{},
		Cursor:    0,
	}
	info := &tui.InfoPanel{
		PanelBase: tui.PanelBase{Title: "Info", Border: true},
		Lines: []string{
			"Multi-panel TUI demo",
			"",
			"Controls:",
			"  Tab       - switch panels",
			"  ↑/↓       - navigate lists",
			"  ←/→       - move cursor in input",
			"  Enter     - select/confirm",
			"  Backspace - delete in input",
			"  q / Ctrl-C / Ctrl-D - quit",
			"",
			"Current panel: Files",
		},
	}

	layout := &tui.VerticalSplit{
		Top: &tui.HorizontalSplit{
			Left:  &tui.PanelNode{Panel: files},
			Right: &tui.PanelNode{Panel: commands},
		},
		Bottom: &tui.HorizontalSplit{
			Left: &tui.HorizontalSplit{
				Left:  &tui.PanelNode{Panel: input},
				Right: &tui.PanelNode{Panel: input2},
			},
			Right: &tui.PanelNode{Panel: info},
		},
	}

	app := tui.NewApp(layout)
	app.Run()
}
