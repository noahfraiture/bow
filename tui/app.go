package tui

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"syscall"
	"time"
)

// App manages the terminal user interface, handling panel layout, input, and rendering.
// It runs the main event loop, positions panels, and redraws the screen.
type App struct {
	term        *terminal
	panels      []Panel
	layout      Layout
	activeIdx   int
	running     bool
	sigch       chan os.Signal
	handler     GlobalHandler
	noDraw      bool     // For testing: skip drawing
	previousOps []drawOp // Previous frame operations for double buffering
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
		slog.Warn("could not enable raw mode", "error", err)
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
	clearScreen()
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

func (a *App) callOnPanelSwitch() {
	panelName := a.panels[a.activeIdx].GetBase().Title
	if panelName == "" {
		panelName = strconv.Itoa(a.activeIdx)
	}
	a.handler.OnPanelSwitch(a, panelName)
}

func (a *App) handleMessage(msg InputMessage) {
	handled, redraw := a.panels[a.activeIdx].Update(msg)
	if redraw {
		// Redraw all panels so that other panels that share data can be redraw too
		a.draw()
	}
	if handled {
		return
	}

	redraw = a.handler.UpdateGlobal(a, msg)
	if redraw {
		a.draw()
	}
}
