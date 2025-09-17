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
	clearScreen()
	a.drawPanels()
	a.drawStatusBar()
	a.drawCursor()
	fmt.Print(reset)
}

func (a *App) drawPanels() {
	for i, p := range a.panels {
		active := i == a.activeIdx
		a.drawPanel(p, active)
	}
}

func (a *App) drawPanel(p Panel, active bool) {
	content := p.Draw(active)
	if content == "" {
		return
	}
	full := p.GetBase().wrapWithBorder(content, active)
	if full == "" {
		return
	}
	lines := strings.Split(full, "\n")
	for j, line := range lines {
		if j >= p.GetBase().h {
			break
		}
		writeAt(p.GetBase().x, p.GetBase().y+j, line)
	}
}

func (a *App) drawStatusBar() {
	status := a.handler.GetStatus()
	writeAt(0, a.term.rows-1, padRightRuneString(status, a.term.cols))
}

func (a *App) drawCursor() {
	if a.activeIdx >= len(a.panels) {
		fmt.Print(HideCursor)
		return
	}
	p := a.panels[a.activeIdx]
	active := true // since it's the active panel
	x, y, show := p.CursorPosition(active)
	if show {
		fmt.Print(ShowCursor)
		writeAt(x, y, "")
	} else {
		fmt.Print(HideCursor)
	}
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
			var buf []byte
			var final byte
			for {
				b, err := a.term.reader.ReadByte()
				if err != nil {
					return InputMessage{}, err
				}
				raw = append(raw, b)
				if (b >= 'A' && b <= 'Z') || b == '~' {
					final = b
					break
				}
				buf = append(buf, b)
			}
			params := strings.Split(string(buf), ";")
			var key Key
			var modifiers []Modifier
			if len(params) > 1 {
				modStr := params[len(params)-1]
				modCode, _ := strconv.Atoi(modStr)
				if modCode&1 != 0 {
					modifiers = append(modifiers, ModShift)
				}
				if modCode&2 != 0 {
					modifiers = append(modifiers, ModAlt)
				}
				if modCode&4 != 0 {
					modifiers = append(modifiers, ModCtrl)
				}
				params = params[:len(params)-1]
			}
			if final == '~' {
				if len(params) > 0 {
					keyCode, _ := strconv.Atoi(params[0])
					switch keyCode {
					case 1:
						key = KeyHome
					case 2:
						key = KeyInsert
					case 3:
						key = KeyDelete
					case 4:
						key = KeyEnd
					case 5:
						key = KeyPageUp
					case 6:
						key = KeyPageDown
					case 9:
						key = KeyTab
					}
				}
			} else {
				switch final {
				case 'A':
					key = KeyUp
				case 'B':
					key = KeyDown
				case 'C':
					key = KeyRight
				case 'D':
					key = KeyLeft
				case 'H':
					key = KeyHome
				case 'F':
					key = KeyEnd
				case 'Z':
					key = KeyTab
					modifiers = append(modifiers, ModShift)
				}
			}
			if key != 0 {
				return newKeyMessageWithModifiers(key, raw, modifiers), nil
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
