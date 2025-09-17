package tui

import (
	"fmt"
	"strings"
)

type drawBuffer struct {
	operations []drawOp
}

type drawOp struct {
	x, y    int
	content string
}

// newDrawBuffer creates a new drawing buffer
func newDrawBuffer() *drawBuffer {
	return &drawBuffer{
		operations: make([]drawOp, 0, 100), // Pre-allocate reasonable capacity
	}
}

func (db *drawBuffer) writeAt(x, y int, content string) {
	db.operations = append(db.operations, drawOp{x: x, y: y, content: content})
}

// flush executes all buffered drawing operations
func (db *drawBuffer) flush() {
	for _, op := range db.operations {
		writeAt(op.x, op.y, op.content)
	}
	db.operations = db.operations[:0] // Clear buffer efficiently
}

// drawPanelsBuffered collects all panel drawing operations
func (a *App) drawPanelsBuffered(buffer *drawBuffer) {
	for i, p := range a.panels {
		active := i == a.activeIdx
		a.drawPanelBuffered(p, active, buffer)
	}
}

// drawPanelBuffered collects drawing operations for a single panel
func (a *App) drawPanelBuffered(p Panel, active bool, buffer *drawBuffer) {
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
		buffer.writeAt(p.GetBase().x, p.GetBase().y+j, line)
	}
}

// drawStatusBarBuffered collects status bar drawing operations
func (a *App) drawStatusBarBuffered(buffer *drawBuffer) {
	status := a.handler.GetStatus()
	buffer.writeAt(0, a.term.rows-1, padRightRuneString(status, a.term.cols))
}

// drawCursorBuffered collects cursor positioning operations
func (a *App) drawCursorBuffered(buffer *drawBuffer) {
	if a.activeIdx >= len(a.panels) {
		fmt.Print(HideCursor)
		return
	}
	p := a.panels[a.activeIdx]
	active := true // since it's the active panel
	x, y, show := p.CursorPosition(active)
	if show {
		fmt.Print(ShowCursor)
		buffer.writeAt(x, y, "")
	} else {
		fmt.Print(HideCursor)
	}
}

// draw performs atomic screen updates using buffered rendering
func (a *App) draw() {
	if a.noDraw {
		return
	}

	// Create buffer and collect all drawing operations
	buffer := newDrawBuffer()
	a.drawPanelsBuffered(buffer)
	a.drawStatusBarBuffered(buffer)
	a.drawCursorBuffered(buffer)

	// Atomic update: clear screen and flush all operations at once
	clearScreen()
	buffer.flush()
	fmt.Print(reset)
}
