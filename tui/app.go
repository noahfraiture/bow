package tui

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// App manages the terminal user interface, handling panel layout, input, and rendering.
// It runs the main event loop, positions panels, and redraws the screen.
type App struct {
	term      *terminal
	panels    []Panel
	layout    Layout
	activeIdx int
	running   bool
	sigch     chan os.Signal
	handler   GlobalHandler
	noDraw    bool // For testing: skip drawing
}

// NewApp creates a new App instance with the given layout and global handler.
// If handler is nil, uses DefaultGlobalHandler.
// Initializes terminal settings and positions panels.
func NewApp(layout Layout, handler GlobalHandler) *App {
	if handler == nil {
		handler = &DefaultGlobalHandler{}
	}
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
		handler: handler,
	}
	app.layoutPanels(layout)
	return app
}

func (a *App) layoutPanels(layout Layout) {
	a.panels = layout.position(0, 0, a.term.cols, a.term.rows-1)
	if a.activeIdx > len(a.panels) {
		a.activeIdx = len(a.panels) - 1
	}
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
		disableRawMode(a.term.prevStty)
		fmt.Print(ShowCursor)
		clearScreen()
	}()

	fmt.Print(HideCursor)
	for _, panel := range a.panels {
		panel.Update(InputMessage{})
	}
	a.draw()

	for a.running {
		if slices.ContainsFunc(a.panels, func(panel Panel) bool {
			return panel.GetBase().stopping
		}) {
			break
		}
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
	if a.noDraw {
		return
	}
	a.layoutPanels(a.layout)
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
	status := a.handler.GetStatus()
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
		cursorX := min(startX+tp.Cursor, maxX)
		writeAt(cursorX, startY, "")
	} else {
		fmt.Print(HideCursor)
	}
	fmt.Print(reset)
}

// SwitchPanel switches to the next or previous panel based on direction.
// direction > 0: next, < 0: previous, 0: no change.
func (a *App) SwitchPanel(direction int) {
	if direction == 0 {
		return
	}
	a.activeIdx = (a.activeIdx + direction + len(a.panels)) % len(a.panels)
	a.callOnPanelSwitch()
}

// FocusPanel focuses the panel with the given name (Title or index string).
// Returns true if found and focused.
func (a *App) FocusPanel(name string) bool {
	for i, panel := range a.panels {
		if panel.GetBase().Title == name || strconv.Itoa(i) == name {
			a.activeIdx = i
			a.callOnPanelSwitch()
			return true
		}
	}
	return false
}

// Stop stops the application by setting running to false.
func (a *App) Stop() {
	a.running = false
}

// callOnPanelSwitch is a helper to call OnPanelSwitch with the current panel name.
func (a *App) callOnPanelSwitch() {
	panelName := a.panels[a.activeIdx].GetBase().Title
	if panelName == "" {
		panelName = strconv.Itoa(a.activeIdx)
	}
	a.handler.OnPanelSwitch(a, panelName)
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
	handled, redraw := a.handler.UpdateGlobal(a, msg)
	if handled {
		if redraw {
			a.draw()
		}
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
