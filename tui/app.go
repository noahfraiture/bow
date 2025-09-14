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

// App manages the terminal user interface, handling panel layout, input, and rendering.
// It runs the main event loop, positions panels, and redraws the screen.
type App struct {
	term      *terminal
	panels    []Panel
	layout    layout
	activeIdx int
	running   bool
	sigch     chan os.Signal
}

// NewApp creates a new App instance with the given layout.
// Initializes terminal settings and positions panels.
func NewApp(layout layout) *App {
	cols, rows, err := getTermSize()
	if err != nil {
		cols, rows = 80, 24
	}
	term := &terminal{
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

func (a *App) layoutPanels(layout layout) {
	a.panels = layout.position(0, 0, a.term.cols, a.term.rows-1)
	a.activeIdx = 0
}

// Run starts the application's main loop, handling input and rendering until quit.
// Enables raw mode, processes events, and cleans up on exit.
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
					a.layoutPanels(a.layout)
					a.draw() // redraw after resize
				}
			} else {
				a.running = false
			}
		}
	}()

	defer func() {
		DisableRawMode(a.term.prevStty)
		fmt.Print(ShowCursor)
		clearScreen()
	}()

	fmt.Print(HideCursor)
	a.draw()

	for a.running {
		msg, err := a.parseInput()
		if err != nil {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		a.handleMessage(msg)
		a.draw()
	}
}

func (a *App) draw() {
	clearScreen()
	for i, p := range a.panels {
		active := i == a.activeIdx
		content := p.Draw(active)
		if content == "" {
			continue
		}
		full := p.GetBase().wrapWithBorder(content, active)
		if full == "" {
			continue
		}
		lines := strings.Split(full, "\n")
		for j, line := range lines {
			if j >= p.GetBase().h {
				break
			}
			writeAt(p.GetBase().x, p.GetBase().y+j, line)
		}
	}
	status := " Tab: switch  •  ↑/↓: navigate  •  ←/→: move cursor  •  Enter: confirm  •  q/Ctrl-C: quit "
	writeAt(0, a.term.rows-1, padRightRuneString(status, a.term.cols))

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
		writeAt(cursorX, startY, "")
	} else {
		fmt.Print(HideCursor)
	}
	fmt.Print(Reset)
}

func (a *App) switchPanel() {
	a.activeIdx = (a.activeIdx + 1) % len(a.panels)
}

// parseInput reads and parses input into an InputMessage
func (a *App) parseInput() (InputMessage, error) {
	b, err := a.term.reader.ReadByte()
	if err != nil {
		return InputMessage{}, err
	}

	raw := []byte{b}

	// Handle escape sequences
	if b == byte(KeyEsc) {
		next1, err := a.term.reader.ReadByte()
		if err != nil {
			return InputMessage{}, err
		}
		raw = append(raw, next1)

		if next1 == '[' {
			next2, err := a.term.reader.ReadByte()
			if err != nil {
				return InputMessage{}, err
			}
			raw = append(raw, next2)

			switch next2 {
			case 'A':
				return newKeyMessage(KeyUp, raw), nil
			case 'B':
				return newKeyMessage(KeyDown, raw), nil
			case 'C':
				return newKeyMessage(KeyRight, raw), nil
			case 'D':
				return newKeyMessage(KeyLeft, raw), nil
			case 'H':
				return newKeyMessage(KeyHome, raw), nil
			case 'F':
				return newKeyMessage(KeyEnd, raw), nil
			case '5':
				tilde, err := a.term.reader.ReadByte()
				if err != nil {
					return InputMessage{}, err
				}
				raw = append(raw, tilde)
				if tilde == '~' {
					return newKeyMessage(KeyPageUp, raw), nil
				}
			case '6':
				tilde, err := a.term.reader.ReadByte()
				if err != nil {
					return InputMessage{}, err
				}
				raw = append(raw, tilde)
				if tilde == '~' {
					return newKeyMessage(KeyPageDown, raw), nil
				}
			case '2':
				tilde, err := a.term.reader.ReadByte()
				if err != nil {
					return InputMessage{}, err
				}
				raw = append(raw, tilde)
				if tilde == '~' {
					return newKeyMessage(KeyInsert, raw), nil
				}
			case '3':
				tilde, err := a.term.reader.ReadByte()
				if err != nil {
					return InputMessage{}, err
				}
				raw = append(raw, tilde)
				if tilde == '~' {
					return newKeyMessage(KeyDelete, raw), nil
				}
			}
		} else if next1 >= 'O' && next1 <= 'Z' {
			// Function keys F1-F4
			switch next1 {
			case 'P':
				return newKeyMessage(KeyF1, raw), nil
			case 'Q':
				return newKeyMessage(KeyF2, raw), nil
			case 'R':
				return newKeyMessage(KeyF3, raw), nil
			case 'S':
				return newKeyMessage(KeyF4, raw), nil
			}
		}
		// Unknown escape sequence, return as special key
		return newKeyMessage(KeyEsc, raw), nil
	}

	// Handle special keys
	switch b {
	case byte(KeyTab):
		return newKeyMessage(KeyTab, raw), nil
	case byte(KeyEnter):
		return newKeyMessage(KeyEnter, raw), nil
	case byte(KeyBackspace):
		return newKeyMessage(KeyBackspace, raw), nil
	case byte(KeySpace):
		return newCharMessage(' ', raw), nil
	}

	// Handle printable characters
	if b >= 32 && b <= 126 {
		return newCharMessage(rune(b), raw), nil
	}

	// Handle control characters
	if b < 32 {
		msg := newCharMessage(rune(b), raw)
		msg.modifiers = append(msg.modifiers, ModCtrl)
		return msg, nil
	}

	// Default to special key
	return newKeyMessage(KeyEsc, raw), nil
}

// handleMessage processes an InputMessage
func (a *App) handleMessage(msg InputMessage) {
	switch {
	case msg.IsKey(KeyTab):
		a.switchPanel()
		return
	case msg.IsChar('q'), msg.IsChar('Q'):
		a.running = false
		return
	case msg.IsChar('\x03'): // Ctrl+C
		a.running = false
		return
	}

	if a.activeIdx < len(a.panels) {
		needsRedraw := a.panels[a.activeIdx].Update(msg)
		if needsRedraw {
			a.draw()
		}
	}
}

// padRightRuneString pads s with spaces to width w, truncating if longer.
func padRightRuneString(s string, w int) string {
	r := []rune(s)
	if len(r) >= w {
		return string(r[:w])
	}
	return string(r) + strings.Repeat(" ", w-len(r))
}
