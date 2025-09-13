package main

import (
	"app/tui"
	"fmt"
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

func (cp *CounterPanel) Draw(active bool) {
	cp.PanelBase.Draw(active)

	// Draw the counter value
	countStr := fmt.Sprintf("Count: %d", cp.Count)
	tui.WriteAt(cp.X+2, cp.Y+2, tui.ClrWhite+countStr+tui.Reset)

	// Draw instructions
	instructions := []string{
		"Use + to increment",
		"Use - to decrement",
	}
	for i, line := range instructions {
		tui.WriteAt(cp.X+2, cp.Y+4+i, tui.ClrWhite+line+tui.Reset)
	}
}

func (cp *CounterPanel) GetBase() *tui.PanelBase {
	return &cp.PanelBase
}

func (cp *CounterPanel) SetPosition(x, y, w, h int) {
	cp.X, cp.Y, cp.W, cp.H = x, y, w, h
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
		PanelBase: tui.PanelBase{Title: "Input", Border: true},
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
