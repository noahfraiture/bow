package tui

import (
	"strings"
	"unicode/utf8"
)

// PanelBase provides common fields and methods for panels.
// It should be embedded (using an unnamed field) in custom panel structs to enable
// positioning, borders, and default behavior. Users can override Update and Draw methods.
type PanelBase struct {
	x, y   int
	w, h   int
	Title  string
	Border bool
}

// GetBase returns the PanelBase instance.
// This is used internally by the layout system; users typically don't need to call it directly.
func (pb *PanelBase) GetBase() *PanelBase {
	return pb
}

// Panel defines the interface for UI panels that can update and draw themselves.
// Implementations should embed PanelBase (or another existing panel) to inherit
// common functionality like positioning and borders. Override Update and Draw as needed, but keep GetBase for internal use.
type Panel interface {
	GetBase() *PanelBase
	Update(input byte) bool
	Draw(active bool) string
}

// Update handles input for the base panel.
// Default implementation does nothing and returns false.
// Override in custom panels to handle user input.
func (pb *PanelBase) Update(input byte) bool {
	return false
}

// Draw renders the base panel's content.
// Default implementation returns empty lines fitting the panel's dimensions.
// Override in custom panels to provide custom content.
func (pb *PanelBase) Draw(active bool) string {
	if pb.w <= 2 || pb.h <= 2 {
		return ""
	}
	lines := make([]string, pb.h-2)
	for i := range lines {
		lines[i] = ""
	}
	return strings.Join(lines, "\n")
}

func (pb *PanelBase) wrapWithBorder(content string, active bool) string {
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
	top := color + "┌" + strings.Repeat("─", pb.w-2) + "┐" + Reset
	if pb.Title != "" {
		title := " [" + pb.Title + "] "
		if utf8.RuneCountInString(title) <= pb.w-2 {
			top = color + "┌" + title + strings.Repeat("─", pb.w-2-utf8.RuneCountInString(title)) + "┐" + Reset
		}
	}
	lines = append(lines, top)

	for _, line := range contentLines {
		truncated := truncateToWidth(line, pb.w-2)
		padded := truncated + strings.Repeat(" ", pb.w-2-displayWidth(truncated))
		borderedLine := color + "│" + Reset + ClrWhite + padded + Reset + color + "│" + Reset
		lines = append(lines, borderedLine)
	}

	for len(lines) < pb.h-1 {
		emptyLine := color + "│" + Reset + ClrWhite + strings.Repeat(" ", pb.w-2) + Reset + color + "│" + Reset
		lines = append(lines, emptyLine)
	}

	bottom := color + "└" + strings.Repeat("─", pb.w-2) + "┘" + Reset
	lines = append(lines, bottom)

	return strings.Join(lines, "\n")
}

type layout interface {
	position(x, y, w, h int) []Panel
}

// PanelNode represents a single panel in the layout.
// It positions the panel to fill the entire given area.
type PanelNode struct {
	Panel Panel
}

// HorizontalSplit divides the area into left and right sections.
// Left takes half the width, right takes the rest.
type HorizontalSplit struct {
	Left  layout
	Right layout
}

// VerticalSplit divides the area into top and bottom sections.
// Top takes half the height, bottom takes the rest.
type VerticalSplit struct {
	Top    layout
	Bottom layout
}

func (pn *PanelNode) position(x, y, w, h int) []Panel {
	pb := pn.Panel.GetBase()
	pb.x, pb.y, pb.w, pb.h = x, y, w, h
	return []Panel{pn.Panel}
}

func (hs *HorizontalSplit) position(x, y, w, h int) []Panel {
	leftW := w / 2
	rightW := w - leftW
	var panels []Panel
	panels = append(panels, hs.Left.position(x, y, leftW, h)...)
	panels = append(panels, hs.Right.position(x+leftW, y, rightW, h)...)
	return panels
}

func (vs *VerticalSplit) position(x, y, w, h int) []Panel {
	topH := h / 2
	bottomH := h - topH
	var panels []Panel
	panels = append(panels, vs.Top.position(x, y, w, topH)...)
	panels = append(panels, vs.Bottom.position(x, y+topH, w, bottomH)...)
	return panels
}
