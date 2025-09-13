package main

import (
	"app/tui"
	"fmt"
	"strings"
)

// CustomPanel demonstrates how to create custom panels
type CounterPanel struct {
	tui.PanelBase
	Count int
}

func (cp *CounterPanel) Update(input byte) bool {
	switch input {
	case '+':
		cp.Count++
		return true
	case '-':
		cp.Count--
		return true
	}
	return false
}

func (cp *CounterPanel) Draw(active bool) string {
	countStr := fmt.Sprintf("Count: %d", cp.Count)
	instructions := []string{
		"Use + to increment",
		"Use - to decrement",
	}
	lines := []string{countStr}
	lines = append(lines, instructions...)
	return strings.Join(lines, "\n")
}

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
		PanelBase: tui.PanelBase{Title: "Input", Border: false},
		Text:      []rune{},
		Cursor:    0,
	}
	counter := &CounterPanel{
		PanelBase: tui.PanelBase{Title: "Counter", Border: true},
		Count:     0,
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
			"  + / -     - change counter",
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
			Right: &tui.HorizontalSplit{
				Left:  &tui.PanelNode{Panel: counter},
				Right: &tui.PanelNode{Panel: info},
			},
		},
	}

	app := tui.NewApp(layout)
	app.Run()
}
