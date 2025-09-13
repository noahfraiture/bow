package tui

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
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
		p.Draw(active)
	}
	status := " Tab: switch  •  ↑/↓: navigate  •  ←/→: move cursor  •  Enter: confirm  •  q/Ctrl-C: quit "
	WriteAt(0, a.term.rows-1, padRightRuneString(status, a.term.cols))

	// Handle cursor visibility for text panels
	if _, ok := a.panels[a.activeIdx].(*TextPanel); ok && a.activeIdx < len(a.panels) {
		fmt.Print(ShowCursor)
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
	case KeyEnter, '\n':
		a.onEnter()
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

// These methods are now handled by panel.Update()
// Keeping them for backward compatibility or future use
func (a *App) onUp()    { /* handled by panel.Update() */ }
func (a *App) onDown()  { /* handled by panel.Update() */ }
func (a *App) onLeft()  { /* handled by panel.Update() */ }
func (a *App) onRight() { /* handled by panel.Update() */ }

func (a *App) onEnter() {
	switch p := a.panels[a.activeIdx].(type) {
	case *ListPanel:
		// echo to info if there's an info panel
		for _, panel := range a.panels {
			if ip, ok := panel.(*InfoPanel); ok {
				ip.Lines = append(ip.Lines, "", "Selected: "+p.Items[p.Selected])
			}
		}
	case *TextPanel:
		for _, panel := range a.panels {
			if ip, ok := panel.(*InfoPanel); ok {
				ip.Lines = append(ip.Lines, "", "Input: "+string(p.Text))
			}
		}
		// Clear text - this logic is now also in TextPanel.Update()
		p.Text = []rune{}
		p.Cursor = 0
	}
}

// These methods are now handled by panel.Update()
// Keeping them for backward compatibility
func (a *App) onBackspace() { /* handled by panel.Update() */ }

func (a *App) onRune(r rune) {
	// Vim-style navigation for non-text panels
	if _, ok := a.panels[a.activeIdx].(*TextPanel); !ok {
		switch r {
		case 'j':
			if a.activeIdx < len(a.panels) {
				a.panels[a.activeIdx].Update(66) // down arrow
				a.draw()
			}
		case 'k':
			if a.activeIdx < len(a.panels) {
				a.panels[a.activeIdx].Update(65) // up arrow
				a.draw()
			}
		case 'h':
			if a.activeIdx < len(a.panels) {
				a.panels[a.activeIdx].Update(68) // left arrow
				a.draw()
			}
		case 'l':
			if a.activeIdx < len(a.panels) {
				a.panels[a.activeIdx].Update(67) // right arrow
				a.draw()
			}
		case 'q', 'Q':
			a.running = false
		}
	}
}
