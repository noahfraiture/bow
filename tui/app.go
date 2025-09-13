package tui

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

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
