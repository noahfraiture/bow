package tui

import (
	"strings"
	"unicode/utf8"
)

type PanelBase struct {
	x, y   int
	w, h   int
	Title  string
	Border bool
}

// Default implementations for PanelBase
func (pb *PanelBase) GetBase() *PanelBase {
	return pb
}

type Panel interface {
	GetBase() *PanelBase
	Update(input byte) bool  // returns true if redraw needed
	Draw(active bool) string // returns the panel content as a string with \n for lines
}

// Default implementations for PanelBase
func (pb *PanelBase) Update(input byte) bool {
	return false // base implementation does nothing
}

func (pb *PanelBase) Draw(active bool) string {
	// Base implementation returns empty content lines
	if pb.w <= 2 || pb.h <= 2 {
		return ""
	}
	lines := make([]string, pb.h-2)
	for i := range lines {
		lines[i] = ""
	}
	return strings.Join(lines, "\n")
}

func (pb *PanelBase) WrapWithBorder(content string, active bool) string {
	if pb.w < 2 || pb.h < 2 {
		return ""
	}

	if !pb.Border {
		contentLines := strings.Split(content, "\n")
		var lines []string
		if pb.Title != "" {
			titleLine := truncateToWidth(pb.Title, pb.w)
			paddedTitle := titleLine + strings.Repeat(" ", pb.w-displayWidth(titleLine))
			lines = append(lines, paddedTitle)
		}
		maxContentLines := pb.h - len(lines)
		if len(contentLines) > maxContentLines {
			contentLines = contentLines[:maxContentLines]
			if len(contentLines) > 0 {
				contentLines[len(contentLines)-1] = "..."
			}
		}
		for _, line := range contentLines {
			truncated := truncateToWidth(line, pb.w)
			padded := truncated + strings.Repeat(" ", pb.w-displayWidth(truncated))
			lines = append(lines, padded)
		}
		for len(lines) < pb.h {
			lines = append(lines, strings.Repeat(" ", pb.w))
		}
		return strings.Join(lines, "\n")
	}

	color := ClrWhite
	if active {
		color = ClrCyan
	}

	contentLines := strings.Split(content, "\n")
	maxContentLines := pb.h - 2
	if len(contentLines) > maxContentLines {
		contentLines = contentLines[:maxContentLines]
		if len(contentLines) > 0 {
			contentLines[len(contentLines)-1] = "..."
		}
	}

	var lines []string
	// Top border with title
	top := color + "┌" + strings.Repeat("─", pb.w-2) + "┐" + Reset
	if pb.Title != "" {
		title := " [" + pb.Title + "] "
		if utf8.RuneCountInString(title) <= pb.w-2 {
			top = color + "┌" + title + strings.Repeat("─", pb.w-2-utf8.RuneCountInString(title)) + "┐" + Reset
		}
	}
	lines = append(lines, top)

	// Content lines
	for _, line := range contentLines {
		truncated := truncateToWidth(line, pb.w-2)
		padded := truncated + strings.Repeat(" ", pb.w-2-displayWidth(truncated))
		borderedLine := color + "│" + Reset + ClrWhite + padded + Reset + color + "│" + Reset
		lines = append(lines, borderedLine)
	}

	// Pad with empty lines if less than max
	for len(lines) < pb.h-1 {
		emptyLine := color + "│" + Reset + ClrWhite + strings.Repeat(" ", pb.w-2) + Reset + color + "│" + Reset
		lines = append(lines, emptyLine)
	}

	// Bottom border
	bottom := color + "└" + strings.Repeat("─", pb.w-2) + "┘" + Reset
	lines = append(lines, bottom)

	return strings.Join(lines, "\n")
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
	pb := pn.Panel.GetBase()
	pb.x, pb.y, pb.w, pb.h = x, y, w, h
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
