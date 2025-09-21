package tui

import (
	"strings"
	"unicode/utf8"
)

// runeWidth returns the display width of a rune.
func runeWidth(r rune) int {
	switch r {
	case '\t':
		return len(expandTabs("\t"))
	case '\n', '\r', '\v', '\f':
		return 0
	case 0:
		return 0
	}
	if r < 32 || (r >= 0x7f && r < 0xa0) {
		return 0
	}
	if (r >= 0x1100 && r <= 0x115f) || r == 0x2329 || r == 0x232a ||
		(r >= 0x2e80 && r <= 0x303e) || (r >= 0x3040 && r <= 0xa4cf) ||
		(r >= 0xac00 && r <= 0xd7a3) || (r >= 0xf900 && r <= 0xfaff) ||
		(r >= 0xfe10 && r <= 0xfe19) || (r >= 0xfe30 && r <= 0xfe6f) ||
		(r >= 0xff00 && r <= 0xff60) || (r >= 0xffe0 && r <= 0xffe6) ||
		(r >= 0x20000 && r <= 0x2fffd) || (r >= 0x30000 && r <= 0x3fffd) {
		return 2
	}
	return 1
}

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

// CursorPosition returns the cursor position for PanelBase.
// Default implementation shows no cursor.
func (pb *PanelBase) CursorPosition(active bool) (x, y int, show bool) {
	return 0, 0, false
}

// Panel defines the interface for UI panels that can update and draw themselves.
// Implementations should embed PanelBase (or another existing panel) to inherit
// common functionality like positioning and borders. Override Update and Draw as needed, but keep GetBase for internal use.
type Panel interface {
	// GetBase returns the PanelBase instance.
	// This is used internally by the layout system; users typically don't need to call it directly if the embed a PanelBase
	GetBase() *PanelBase
	// Update handles input messages for the panel.
	// Returns handled (true if input was processed) and redraw (true if panel needs redrawing).
	Update(msg InputMessage) (handled bool, redraw bool)
	// Draw renders the panel's content as a string.
	// The active parameter indicates if the panel is currently active.
	Draw(active bool) string
	// CursorPosition returns the cursor position within the panel.
	// Returns x, y coordinates and whether to show the cursor.
	CursorPosition(active bool) (x, y int, show bool)
}

// Update handles input for the base panel.
// Default implementation does nothing and returns false, false.
// Override in custom panels to handle user input.
func (pb *PanelBase) Update(msg InputMessage) (handled bool, redraw bool) {
	return false, false
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
		return pb.renderWithoutBorder(content)
	}

	return pb.renderWithBorder(content, active)
}

func (pb *PanelBase) renderWithoutBorder(content string) string {
	var builder strings.Builder
	lineCount := 0
	contentLines := strings.Split(content, "\n")
	if pb.Title != "" {
		titleLine := truncateToWidth(pb.Title, pb.w)
		paddedTitle := titleLine + strings.Repeat(" ", pb.w-displayWidth(titleLine))
		builder.WriteString(paddedTitle)
		builder.WriteString("\n")
		lineCount++
	}
	maxContentLines := pb.h - lineCount
	if len(contentLines) > maxContentLines {
		contentLines = contentLines[:maxContentLines]
		if len(contentLines) > 0 {
			contentLines[len(contentLines)-1] = "..."
		}
	}
	for _, line := range contentLines {
		truncated := truncateToWidth(line, pb.w)
		padded := truncated + strings.Repeat(" ", pb.w-displayWidth(truncated))
		builder.WriteString(padded)
		builder.WriteString("\n")
		lineCount++
	}
	for lineCount < pb.h {
		builder.WriteString(strings.Repeat(" ", pb.w))
		builder.WriteString("\n")
		lineCount++
	}
	return strings.TrimSuffix(builder.String(), "\n")
}

func (pb *PanelBase) renderWithBorder(content string, active bool) string {
	var builder strings.Builder
	lineCount := 0
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
	builder.WriteString(pb.buildTopBorder(color))
	builder.WriteString("\n")
	lineCount++
	for _, line := range contentLines {
		builder.WriteString(pb.buildContentLine(line, color))
		builder.WriteString("\n")
		lineCount++
	}
	for lineCount < pb.h-1 {
		builder.WriteString(pb.buildEmptyLine(color))
		builder.WriteString("\n")
		lineCount++
	}
	builder.WriteString(pb.buildBottomBorder(color))
	return builder.String()
}

func (pb *PanelBase) buildTopBorder(color string) string {
	top := color + "┌" + strings.Repeat("─", pb.w-2) + "┐" + reset
	if pb.Title != "" {
		title := " [" + pb.Title + "] "
		if displayWidth(title) <= pb.w-2 {
			top = color + "┌" + title + strings.Repeat("─", pb.w-2-displayWidth(title)) + "┐" + reset
		}
	}
	return top
}

func (pb *PanelBase) buildContentLine(line, color string) string {
	truncated := truncateToWidth(line, pb.w-2)
	padded := truncated + strings.Repeat(" ", pb.w-2-displayWidth(truncated))
	return color + "│" + reset + clrWhite + padded + reset + color + "│" + reset
}

func (pb *PanelBase) buildEmptyLine(color string) string {
	return color + "│" + reset + clrWhite + strings.Repeat(" ", pb.w-2) + reset + color + "│" + reset
}

func (pb *PanelBase) buildBottomBorder(color string) string {
	return color + "└" + strings.Repeat("─", pb.w-2) + "┘" + reset
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
				width := runeWidth(r)
				if visible+width <= w {
					result.WriteRune(r)
					visible += width
					if visible == w-2 && w > 2 {
						result.WriteString("..")
						truncated = true
						visible += 2
					}
				} else if w-visible >= 2 && !truncated {
					result.WriteString("..")
					truncated = true
					visible += 2
				} else {
					truncated = true
				}
			}
		}
	}
	return result.String()
}

func expandTabs(s string) string {
	return strings.ReplaceAll(s, "\t", "    ")
}

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
			count += runeWidth(r)
		}
		i += size
	}
	return count
}

type Layout interface {
	position(x, y, w, h int) []Panel
	GetWeight() int
}

// PanelNode represents a single panel in the layout.
// Take space with proportion Weight / TotalWeight of layout. Minimum Weight of 1.
// It positions the panel to fill the entire given area.
type PanelNode struct {
	Panel  Panel
	Weight int
}

// HorizontalSplit divides the area into multiple horizontal sections.
// Take space with proportion Weight / TotalWeight of layout. Minimum Weight of 1.
// Each panel gets space proportional to its weight.
type HorizontalSplit struct {
	Panels []Layout
	Weight int
}

// VerticalSplit divides the area into multiple vertical sections.
// Take space with proportion Weight / TotalWeight of layout. Minimum Weight of 1.
// Each panel gets space proportional to its weight.
type VerticalSplit struct {
	Panels []Layout
	Weight int
}

func (pn *PanelNode) position(x, y, w, h int) []Panel {
	pb := pn.Panel.GetBase()
	pb.x, pb.y, pb.w, pb.h = x, y, w, h
	return []Panel{pn.Panel}
}

func (hs *HorizontalSplit) position(x, y, w, h int) []Panel {
	return hs.positionPanels(x, y, w, h, true) // true for horizontal
}

func (vs *VerticalSplit) position(x, y, w, h int) []Panel {
	return vs.positionPanels(x, y, w, h, false) // false for vertical
}

// GetWeight returns the weight of this layout for proportional sizing
func (pn *PanelNode) GetWeight() int {
	if pn.Weight <= 0 {
		return 1
	}
	return pn.Weight
}

// GetWeight returns the weight of this layout for proportional sizing
func (hs *HorizontalSplit) GetWeight() int {
	if hs.Weight <= 0 {
		return 1
	}
	return hs.Weight
}

// GetWeight returns the weight of this layout for proportional sizing
func (vs *VerticalSplit) GetWeight() int {
	if vs.Weight <= 0 {
		return 1
	}
	return vs.Weight
}

// positionPanelsWithWeights handles the common logic for positioning panels with weights
func positionPanelsWithWeights(panels []Layout, x, y, w, h int, horizontal bool) []Panel {
	if len(panels) == 0 {
		return []Panel{}
	}

	totalWeight := 0
	for _, panel := range panels {
		totalWeight += panel.GetWeight()
	}

	if totalWeight == 0 {
		totalWeight = len(panels)
	}

	var result []Panel
	currentPos := x
	if !horizontal {
		currentPos = y
	}

	// Distribute space proportionally
	usedSpace := 0
	for i, panel := range panels {
		var panelW, panelH int
		if horizontal {
			// fill gap due to rounding division
			if i == len(panels)-1 {
				panelW = w - usedSpace
			} else {
				panelW = (panel.GetWeight() * w) / totalWeight
				usedSpace += panelW
			}
			panelH = h
		} else {
			// fill gap due to rounding division
			if i == len(panels)-1 {
				panelH = h - usedSpace
			} else {
				panelH = (panel.GetWeight() * h) / totalWeight
				usedSpace += panelH
			}
			panelW = w
		}

		var panelX, panelY int
		if horizontal {
			panelX = currentPos
			panelY = y
			currentPos += panelW
		} else {
			panelX = x
			panelY = currentPos
			currentPos += panelH
		}

		childPanels := panel.position(panelX, panelY, panelW, panelH)
		result = append(result, childPanels...)
	}

	return result
}

func (hs *HorizontalSplit) positionPanels(x, y, w, h int, horizontal bool) []Panel {
	return positionPanelsWithWeights(hs.Panels, x, y, w, h, horizontal)
}

func (vs *VerticalSplit) positionPanels(x, y, w, h int, horizontal bool) []Panel {
	return positionPanelsWithWeights(vs.Panels, x, y, w, h, horizontal)
}
