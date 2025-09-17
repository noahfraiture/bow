package tui

import (
	"strings"
	"unicode/utf8"
)

// PanelBase provides common fields and methods for panels.
// It should be embedded (using an unnamed field) in custom panel structs to enable
// positioning, borders, and default behavior. Users can override Update and Draw methods.
type PanelBase struct {
	x, y     int
	w, h     int
	Title    string
	Border   bool
	stopping bool
}

// Stop give the possibility to properly stop the app.
// Sending stop to any panel will stop the application at the next app update
func (pb *PanelBase) Stop() {
	pb.stopping = true
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
	Update(msg InputMessage) bool
	Draw(active bool) string
}

// Update handles input for the base panel.
// Default implementation does nothing and returns false.
// Override in custom panels to handle user input.
func (pb *PanelBase) Update(msg InputMessage) bool {
	return false
}

// Draw renders the base panel's content.
// Default implementation returns empty lines fitting the panel's dimensions.
// Override in custom panels to provide custom content.
func (pb *PanelBase) Draw(active bool) string {
	return ""
}

func (pb *PanelBase) wrapWithBorder(content string, active bool) string {
	if pb.w <= 2 || pb.h <= 2 {
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

	color := clrWhite
	if active {
		color = clrCyan
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
	top := color + "┌" + strings.Repeat("─", pb.w-2) + "┐" + reset
	if pb.Title != "" {
		title := " [" + pb.Title + "] "
		if utf8.RuneCountInString(title) <= pb.w-2 {
			top = color + "┌" + title + strings.Repeat("─", pb.w-2-utf8.RuneCountInString(title)) + "┐" + reset
		}
	}
	lines = append(lines, top)

	for _, line := range contentLines {
		truncated := truncateToWidth(line, pb.w-2)
		padded := truncated + strings.Repeat(" ", pb.w-2-displayWidth(truncated))
		borderedLine := color + "│" + reset + clrWhite + padded + reset + color + "│" + reset
		lines = append(lines, borderedLine)
	}

	for len(lines) < pb.h-1 {
		emptyLine := color + "│" + reset + clrWhite + strings.Repeat(" ", pb.w-2) + reset + color + "│" + reset
		lines = append(lines, emptyLine)
	}

	bottom := color + "└" + strings.Repeat("─", pb.w-2) + "┘" + reset
	lines = append(lines, bottom)

	return strings.Join(lines, "\n")
}

// truncateToWidth truncates s to width w, preserving ANSI codes and adding ".." if needed.
func truncateToWidth(s string, w int) string {
	if displayWidth(s) <= w {
		return s
	}
	var result strings.Builder
	visible := 0
	inEscape := false
	truncated := false
	for _, r := range s {
		if r == '\x1b' {
			inEscape = true
			result.WriteRune(r)
		} else if inEscape {
			result.WriteRune(r)
			if r == 'm' {
				inEscape = false
			}
		} else {
			if !truncated {
				if visible < w {
					result.WriteRune(r)
					visible++
					if visible == w-2 && w > 2 {
						result.WriteString("..")
						truncated = true
						visible += 2
					}
				} else {
					truncated = true
				}
			}
		}
	}
	return result.String()
}

// displayWidth calculates the visible width of s, ignoring ANSI escape sequences.
func displayWidth(s string) int {
	count := 0
	inEscape := false
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == '\x1b' {
			inEscape = true
		}
		if inEscape {
			if r == 'm' {
				inEscape = false
			}
		} else {
			count++
		}
		i += size
	}
	return count
}

type Layout interface {
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
	Left  Layout
	Right Layout
}

// VerticalSplit divides the area into top and bottom sections.
// Top takes half the height, bottom takes the rest.
type VerticalSplit struct {
	Top    Layout
	Bottom Layout
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
