package tui

import (
	"bufio"
	"os"
)

type Terminal struct {
	cols     int
	rows     int
	reader   *bufio.Reader
	prevStty string
}

type PanelBase struct {
	x, y   int
	w, h   int
	Title  string
	Border bool
}

type ListPanel struct {
	PanelBase
	Items    []string
	Selected int
}

type TextPanel struct {
	PanelBase
	Text   []rune
	Cursor int
}

type InfoPanel struct {
	PanelBase
	Lines []string
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

type Panel interface {
	GetBase() *PanelBase
	SetPosition(x, y, w, h int)
}

// Implement Panel for each
func (lp *ListPanel) GetBase() *PanelBase { return &lp.PanelBase }
func (tp *TextPanel) GetBase() *PanelBase { return &tp.PanelBase }
func (ip *InfoPanel) GetBase() *PanelBase { return &ip.PanelBase }

func (lp *ListPanel) SetPosition(x, y, w, h int) {
	lp.x, lp.y, lp.w, lp.h = x, y, w, h
}
func (tp *TextPanel) SetPosition(x, y, w, h int) {
	tp.x, tp.y, tp.w, tp.h = x, y, w, h
}
func (ip *InfoPanel) SetPosition(x, y, w, h int) {
	ip.x, ip.y, ip.w, ip.h = x, y, w, h
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

type App struct {
	term      *Terminal
	panels    []Panel
	layout    Layout
	activeIdx int
	running   bool
	sigch     chan os.Signal
}
