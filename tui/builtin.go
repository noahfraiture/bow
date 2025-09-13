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
	if len(tp.Text) == 0 {
		return " "
	}
	return string(tp.Text)
}

// InfoPanel implementations
func (ip *InfoPanel) Update(input byte) bool {
	// InfoPanel doesn't handle input directly
	return false
}

func (ip *InfoPanel) Draw(active bool) string {
	return strings.Join(ip.Lines, "\n")
}
