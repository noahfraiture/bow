package tui

import "strings"

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

// Implement Panel for each
func (lp *ListPanel) GetBase() *PanelBase { return &lp.PanelBase }
func (tp *TextPanel) GetBase() *PanelBase { return &tp.PanelBase }
func (ip *InfoPanel) GetBase() *PanelBase { return &ip.PanelBase }

func (lp *ListPanel) SetPosition(x, y, w, h int) {
	lp.X, lp.Y, lp.W, lp.H = x, y, w, h
}
func (tp *TextPanel) SetPosition(x, y, w, h int) {
	tp.X, tp.Y, tp.W, tp.H = x, y, w, h
}
func (ip *InfoPanel) SetPosition(x, y, w, h int) {
	ip.X, ip.Y, ip.W, ip.H = x, y, w, h
}

// ListPanel implementations
func (lp *ListPanel) Update(input byte) bool {
	switch input {
	case 'k', 65: // up arrow
		if lp.Selected > 0 {
			lp.Selected--
			return true
		}
	case 'j', 66: // down arrow
		if lp.Selected < len(lp.Items)-1 {
			lp.Selected++
			return true
		}
	}
	return false
}

func (lp *ListPanel) Draw(active bool) string {
	baseStr := lp.PanelBase.Draw(active)
	if baseStr == "" {
		return ""
	}
	lines := strings.Split(baseStr, "\n")
	if len(lines) < 3 {
		return baseStr
	}

	// Fill content lines
	contentHeight := len(lines) - 2
	for i, item := range lp.Items {
		if i >= contentHeight {
			if i == contentHeight {
				// Last line: "..."
				ellipsis := "..."
				if len(ellipsis) > lp.W-2 {
					ellipsis = ellipsis[:lp.W-2]
				}
				spaces := lp.W - 2 - len(ellipsis)
				if spaces < 0 {
					spaces = 0
				}
				lines[len(lines)-2] = ClrWhite + "│" + Reset + ClrWhite + ellipsis + strings.Repeat(" ", spaces) + Reset + ClrWhite + "│" + Reset
			}
			break
		}
		color := ClrWhite
		if i == lp.Selected {
			if active {
				color = Reverse
			} else {
				color = ClrYellow
			}
		}
		truncated := truncateToWidth(item, lp.W-2)
		spaces := lp.W - 2 - len(truncated)
		if spaces < 0 {
			spaces = 0
		}
		padded := truncated + strings.Repeat(" ", spaces)
		lines[i+1] = ClrWhite + "│" + Reset + color + padded + Reset + ClrWhite + "│" + Reset
	}

	return strings.Join(lines, "\n")
}

// TextPanel implementations
func (tp *TextPanel) Update(input byte) bool {
	switch input {
	case 68: // left arrow
		if tp.Cursor > 0 {
			tp.Cursor--
			return true
		}
	case 67: // right arrow
		if tp.Cursor < len(tp.Text) {
			tp.Cursor++
			return true
		}
	case 127, 8: // backspace
		if tp.Cursor > 0 && len(tp.Text) > 0 {
			i := tp.Cursor
			tp.Text = append(tp.Text[:i-1], tp.Text[i:]...)
			tp.Cursor--
			return true
		}
	case 13, 10: // enter
		// Clear text on enter
		tp.Text = []rune{}
		tp.Cursor = 0
		return true
	default:
		if input >= 32 && input <= 126 { // printable characters
			i := tp.Cursor
			before := tp.Text[:i]
			after := tp.Text[i:]
			newText := make([]rune, 0, len(before)+1+len(after))
			newText = append(newText, before...)
			newText = append(newText, rune(input))
			newText = append(newText, after...)
			tp.Text = newText
			tp.Cursor++
			return true
		}
	}
	return false
}

func (tp *TextPanel) Draw(active bool) string {
	baseStr := tp.PanelBase.Draw(active)
	if baseStr == "" {
		return ""
	}
	lines := strings.Split(baseStr, "\n")
	if len(lines) < 3 {
		return baseStr
	}

	// Draw text content on first content line
	textStr := string(tp.Text)
	if len(textStr) > tp.W-2 {
		textStr = truncateToWidth(textStr, tp.W-2)
	}
	padded := textStr + strings.Repeat(" ", tp.W-2-len(textStr))
	lines[1] = ClrWhite + "│" + Reset + ClrWhite + padded + Reset + ClrWhite + "│" + Reset

	// Note: Cursor drawing is handled separately in app.draw() if needed
	return strings.Join(lines, "\n")
}

// InfoPanel implementations
func (ip *InfoPanel) Update(input byte) bool {
	// InfoPanel doesn't handle input directly
	return false
}

func (ip *InfoPanel) Draw(active bool) string {
	baseStr := ip.PanelBase.Draw(active)
	if baseStr == "" {
		return ""
	}
	lines := strings.Split(baseStr, "\n")
	if len(lines) < 3 {
		return baseStr
	}

	contentHeight := len(lines) - 2
	for i, line := range ip.Lines {
		if i >= contentHeight {
			if i == contentHeight {
				// Last line: "..."
				ellipsis := "..."
				if len(ellipsis) > ip.W-2 {
					ellipsis = ellipsis[:ip.W-2]
				}
				spaces := ip.W - 2 - len(ellipsis)
				if spaces < 0 {
					spaces = 0
				}
				lines[len(lines)-2] = ClrWhite + "│" + Reset + ClrWhite + ellipsis + strings.Repeat(" ", spaces) + Reset + ClrWhite + "│" + Reset
			}
			break
		}
		truncated := truncateToWidth(line, ip.W-2)
		spaces := ip.W - 2 - len(truncated)
		if spaces < 0 {
			spaces = 0
		}
		padded := truncated + strings.Repeat(" ", spaces)
		lines[i+1] = ClrWhite + "│" + Reset + ClrWhite + padded + Reset + ClrWhite + "│" + Reset
	}

	return strings.Join(lines, "\n")
}
