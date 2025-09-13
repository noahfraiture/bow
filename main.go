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
	baseStr := cp.PanelBase.Draw(active)
	if baseStr == "" {
		return ""
	}
	lines := strings.Split(baseStr, "\n")
	if len(lines) < 5 {
		return baseStr
	}

	// Draw the counter value on line 2
	countStr := fmt.Sprintf("Count: %d", cp.Count)
	if len(countStr) > cp.W-2 {
		countStr = countStr[:cp.W-5] + "..."
	}
	spaces := cp.W - 2 - len(countStr)
	if spaces < 0 {
		spaces = 0
	}
	padded := countStr + strings.Repeat(" ", spaces)
	lines[1] = tui.ClrWhite + "│" + tui.Reset + tui.ClrWhite + padded + tui.Reset + tui.ClrWhite + "│" + tui.Reset

	// Draw instructions on lines 3 and 4
	instructions := []string{
		"Use + to increment",
		"Use - to decrement",
	}
	for i, line := range instructions {
		if i+2 >= len(lines)-1 {
			break
		}
		if len(line) > cp.W-2 {
			line = line[:cp.W-5] + "..."
		}
		spaces := cp.W - 2 - len(line)
		if spaces < 0 {
			spaces = 0
		}
		padded := line + strings.Repeat(" ", spaces)
		lines[i+2] = tui.ClrWhite + "│" + tui.Reset + tui.ClrWhite + padded + tui.Reset + tui.ClrWhite + "│" + tui.Reset
	}

	return strings.Join(lines, "\n")
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
