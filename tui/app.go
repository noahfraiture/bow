package tui

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type App struct {
	term      *Terminal
	panels    []Panel
	layout    Layout
	activeIdx int
	running   bool
	sigch     chan os.Signal
}

func NewApp(layout Layout) *App {
	cols, rows, err := getTermSize()
	if err != nil {
		cols, rows = 80, 24
	}
	term := &Terminal{
		cols:   cols,
		rows:   rows,
		reader: bufio.NewReader(os.Stdin),
	}
	app := &App{
		term:    term,
		layout:  layout,
		running: true,
	}
	app.layoutPanels(layout)
	return app
}

func (a *App) layoutPanels(layout Layout) {
	a.panels = layout.Position(0, 0, a.term.cols, a.term.rows-1) // leave space for status
	a.activeIdx = 0
}

func (a *App) Run() {
	prev, err := enableRawMode()
	if err == nil {
		a.term.prevStty = prev
	} else {
		fmt.Fprintln(os.Stderr, "warning: could not enable raw mode:", err)
	}

	a.sigch = make(chan os.Signal, 1)
	signal.Notify(a.sigch, syscall.SIGWINCH, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for s := range a.sigch {
			if s == syscall.SIGWINCH {
				cols, rows, err := getTermSize()
				if err == nil {
					a.term.cols = cols
					a.term.rows = rows
					a.layoutPanels(a.layout) // need to store layout
					a.draw()                 // redraw after resize
				}
			} else {
				a.running = false
			}
		}
	}()

	defer func() {
		disableRawMode(a.term.prevStty)
		fmt.Print(ShowCursor)
		clearScreen()
	}()

	fmt.Print(HideCursor)
	a.draw()

	for a.running {
		b, err := a.term.reader.ReadByte()
		if err != nil {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		if b == 0x03 || b == 0x04 {
			a.running = false
			break
		}
		a.handleByte(b)
		a.draw()
	}
}

// Drawing
func (a *App) draw() {
	clearScreen()
	for i, p := range a.panels {
		active := i == a.activeIdx
		content := p.Draw(active)
		if content == "" {
			continue
		}
		full := p.GetBase().WrapWithBorder(content, active)
		if full == "" {
			continue
		}
		lines := strings.Split(full, "\n")
		for j, line := range lines {
			if j >= p.GetBase().h {
				break
			}
			WriteAt(p.GetBase().x, p.GetBase().y+j, line)
		}
	}
	status := " Tab: switch  •  ↑/↓: navigate  •  ←/→: move cursor  •  Enter: confirm  •  q/Ctrl-C: quit "
	WriteAt(0, a.term.rows-1, padRightRuneString(status, a.term.cols))

	// Handle cursor visibility and position for text panels
	if tp, ok := a.panels[a.activeIdx].(*TextPanel); ok && a.activeIdx < len(a.panels) {
		fmt.Print(ShowCursor)
		var startX, startY, maxX int
		if tp.Border {
			startX = tp.x + 1
			startY = tp.y + 1
			maxX = tp.x + tp.w - 2
		} else {
			startX = tp.x
			if tp.Title != "" {
				startY = tp.y + 1
			} else {
				startY = tp.y
			}
			maxX = tp.x + tp.w - 1
		}
		cursorX := startX + tp.Cursor
		if cursorX > maxX {
			cursorX = maxX
		}
		WriteAt(cursorX, startY, "")
	} else {
		fmt.Print(HideCursor)
	}
	fmt.Print(Reset)
}

// Input handling
func (a *App) handleByte(b byte) {
	// Handle escape sequences for arrow keys
	if b == KeyEsc {
		next1, err := a.term.reader.ReadByte()
		if err != nil {
			return
		}
		if next1 == '[' {
			next2, err := a.term.reader.ReadByte()
			if err != nil {
				return
			}
			// Convert escape sequences to simple bytes for panel.Update()
			switch next2 {
			case 'A':
				b = 65 // up arrow -> 'A'
			case 'B':
				b = 66 // down arrow -> 'B'
			case 'C':
				b = 67 // right arrow -> 'C'
			case 'D':
				b = 68 // left arrow -> 'D'
			}
		} else {
			return
		}
	}

	// Handle app-level keys first
	switch b {
	case KeyTab:
		a.switchPanel()
		return
	case 'q', 'Q', 3, 4: // quit keys (q, Q, Ctrl-C, Ctrl-D)
		a.running = false
		return
	}

	// Let the active panel handle the input
	if a.activeIdx < len(a.panels) {
		needsRedraw := a.panels[a.activeIdx].Update(b)
		if needsRedraw {
			a.draw()
		}
	}
}

func (a *App) switchPanel() {
	a.activeIdx = (a.activeIdx + 1) % len(a.panels)
}
