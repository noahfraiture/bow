package tui

import "strings"

// ListPanel is a panel that displays a list of items with selection.
// Embed PanelBase for positioning and borders. Override Draw to customize appearance.
type ListPanel struct {
	PanelBase
	Items    []string
	Selected int
}

// TextPanel is a panel for text input and display.
// Embed PanelBase for positioning and borders. Handles cursor movement and text editing.
type TextPanel struct {
	PanelBase
	Text   []rune
	Cursor int
}

// InfoPanel is a panel for displaying informational lines.
// Embed PanelBase for positioning and borders. Simple display-only panel.
type InfoPanel struct {
	PanelBase
	Lines []string
}

// Update handles input for the ListPanel, updating selection based on keys.
// Returns true if the panel needs to be redrawn.
func (lp *ListPanel) Update(input byte) bool {
	switch input {
	case 'k', 65:
		if lp.Selected > 0 {
			lp.Selected--
			return true
		}
	case 'j', 66:
		if lp.Selected < len(lp.Items)-1 {
			lp.Selected++
			return true
		}
	}
	return false
}

// Draw renders the ListPanel's content as a string.
// Highlights the selected item based on active state.
func (lp *ListPanel) Draw(active bool) string {
	var lines []string
	for i, item := range lp.Items {
		color := ""
		if i == lp.Selected {
			if active {
				color = Reverse
			} else {
				color = ClrYellow
			}
		}
		lines = append(lines, color+item+Reset)
	}
	return strings.Join(lines, "\n")
}

// Update handles input for the TextPanel, managing cursor and text editing.
// Supports arrow keys, backspace, enter, and printable characters.
// Returns true if the panel needs to be redrawn.
func (tp *TextPanel) Update(input byte) bool {
	switch input {
	case 68:
		if tp.Cursor > 0 {
			tp.Cursor--
			return true
		}
	case 67:
		if tp.Cursor < len(tp.Text) {
			tp.Cursor++
			return true
		}
	case 127, 8:
		if tp.Cursor > 0 && len(tp.Text) > 0 {
			i := tp.Cursor
			tp.Text = append(tp.Text[:i-1], tp.Text[i:]...)
			tp.Cursor--
			return true
		}
	case 13, 10:
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

// Draw renders the TextPanel's text content as a string.
// Returns the current text or a space if empty.
func (tp *TextPanel) Draw(active bool) string {
	if len(tp.Text) == 0 {
		return " "
	}
	return string(tp.Text)
}

// Update handles input for the InfoPanel.
// No-op implementation; InfoPanel does not respond to input.
// Returns false.
func (ip *InfoPanel) Update(input byte) bool {
	return false
}

// Draw renders the InfoPanel's lines as a string.
// Joins all lines with newlines.
func (ip *InfoPanel) Draw(active bool) string {
	return strings.Join(ip.Lines, "\n")
}
