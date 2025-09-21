package tui

import (
	"fmt"
	"strings"
)

type drawBuffer struct {
	operations  []drawOp
	previousOps []drawOp
}

type drawOp struct {
	x, y    int
	content string
}

func newDrawBuffer(app *App) *drawBuffer {
	return &drawBuffer{
		operations:  make([]drawOp, 0, 100),
		previousOps: app.previousOps,
	}
}

func (db *drawBuffer) writeAt(x, y int, content string) {
	db.operations = append(db.operations, drawOp{x: x, y: y, content: content})
}

func (db *drawBuffer) flush() {
	changedOps := db.findChangedOperations()
	for _, op := range changedOps {
		writeAt(op.x, op.y, op.content)
	}
	db.updatePreviousState()
}

func (db *drawBuffer) findChangedOperations() []drawOp {
	changedOps := make([]drawOp, 0, len(db.operations))

	prevMap := make(map[string]drawOp)
	for _, op := range db.previousOps {
		key := db.makeKey(op.x, op.y)
		prevMap[key] = op
	}

	for _, currOp := range db.operations {
		key := db.makeKey(currOp.x, currOp.y)
		if prevOp, exists := prevMap[key]; !exists || prevOp.content != currOp.content {
			changedOps = append(changedOps, currOp)
		}
		delete(prevMap, key)
	}

	for _, oldOp := range prevMap {
		changedOps = append(changedOps, drawOp{
			x:       oldOp.x,
			y:       oldOp.y,
			content: "",
		})
	}
	return changedOps
}

func (db *drawBuffer) updatePreviousState() {
	db.previousOps = db.previousOps[:0]
	db.previousOps = append(db.previousOps, db.operations...)
	db.operations = db.operations[:0]
}

func (db *drawBuffer) makeKey(x, y int) string {
	return fmt.Sprintf("%d,%d", x, y)
}

func (a *App) drawPanelsBuffered(buffer *drawBuffer) {
	for i, p := range a.panels {
		active := i == a.activeIdx
		a.drawPanelBuffered(p, active, buffer)
	}
}

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

func (a *App) drawStatusBarBuffered(buffer *drawBuffer) {
	status := a.handler.GetStatus()
	buffer.writeAt(0, a.term.rows-1, padRightRuneString(status, a.term.cols))
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
	active := true
	x, y, show := p.CursorPosition(active)
	if show {
		fmt.Printf("\x1b[%d;%dH", y+1, x+1)
		fmt.Print(ShowCursor)
	} else {
		fmt.Print(HideCursor)
	}
}

func (a *App) draw() {
	if a.noDraw {
		return
	}

	// Reposition layout in case it changed
	a.layoutPanels(a.layout)

	if a.disableDoubleBuffer {
		clearScreen()
		a.drawPanels()
		a.drawStatusBar()
		a.drawCursor()
		fmt.Print(reset)
	} else {
		buffer := newDrawBuffer(a)
		a.drawPanelsBuffered(buffer)
		a.drawStatusBarBuffered(buffer)

		buffer.flush()
		a.previousOps = buffer.previousOps
		a.drawCursor()
		fmt.Print(reset)
	}
}

func padRightRuneString(s string, w int) string {
	if displayWidth(s) >= w {
		return truncateToWidth(s, w)
	}
	return s + strings.Repeat(" ", w-displayWidth(s))
}
