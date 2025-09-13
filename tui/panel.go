package tui

import (
	"strings"
)

type PanelBase struct {
	X, Y   int // made public for custom panels
	W, H   int // made public for custom panels
	Title  string
	Border bool
}

type Panel interface {
	GetBase() *PanelBase
	SetPosition(x, y, w, h int)
	Update(input byte) bool // returns true if redraw needed
	Draw(active bool)       // handles drawing the panel (active indicates if this panel has focus)
}

// Default implementations for PanelBase
func (pb *PanelBase) Update(input byte) bool {
	return false // base implementation does nothing
}

func (pb *PanelBase) Draw(active bool) {
	// Base implementation draws border and title
	if pb.W < 2 || pb.H < 2 {
		return
	}

	color := ClrWhite
	if active {
		color = ClrCyan
	}
	WriteAt(pb.X, pb.Y, color+"┌"+strings.Repeat("─", pb.W-2)+"┐"+Reset)

	if pb.Title != "" {
		title := " [" + pb.Title + "] "
		if len(title) <= pb.W-2 {
			WriteAt(pb.X+1, pb.Y, title)
		}
	}

	for i := 1; i < pb.H-1; i++ {
		WriteAt(pb.X, pb.Y+i, color+"│"+Reset)
		WriteAt(pb.X+pb.W-1, pb.Y+i, color+"│"+Reset)
	}

	WriteAt(pb.X, pb.Y+pb.H-1, color+"└"+strings.Repeat("─", pb.W-2)+"┘"+Reset)
}

// Layout types for tree-like structure
type Layout interface {
	Position(x, y, w, h int) []Panel
}

type PanelNode struct {
	Panel Panel
}

type HorizontalSplit struct {
	Left  Layout
	Right Layout
}

type VerticalSplit struct {
	Top    Layout
	Bottom Layout
}

// Position methods
func (pn *PanelNode) Position(x, y, w, h int) []Panel {
	pn.Panel.SetPosition(x, y, w, h)
	return []Panel{pn.Panel}
}

func (hs *HorizontalSplit) Position(x, y, w, h int) []Panel {
	leftW := w / 2
	rightW := w - leftW
	var panels []Panel
	panels = append(panels, hs.Left.Position(x, y, leftW, h)...)
	panels = append(panels, hs.Right.Position(x+leftW, y, rightW, h)...)
	return panels
}

func (vs *VerticalSplit) Position(x, y, w, h int) []Panel {
	topH := h / 2
	bottomH := h - topH
	var panels []Panel
	panels = append(panels, vs.Top.Position(x, y, w, topH)...)
	panels = append(panels, vs.Bottom.Position(x, y+topH, w, bottomH)...)
	return panels
}
