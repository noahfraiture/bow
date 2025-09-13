package tui

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

func (lp *ListPanel) Draw(active bool) {
	// Draw base first
	lp.PanelBase.Draw(active)

	// Draw list items
	for i, item := range lp.Items {
		if i >= lp.H-2 { // leave space for borders
			break
		}
		y := lp.Y + 1 + i
		color := ClrWhite
		if i == lp.Selected {
			if active {
				color = Reverse
			} else {
				color = ClrYellow
			}
		}
		truncated := truncateToWidth(item, lp.W-2)
		WriteAt(lp.X+1, y, color+truncated+Reset)
	}
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

func (tp *TextPanel) Draw(active bool) {
	// Draw base first
	tp.PanelBase.Draw(active)

	// Draw text content
	textStr := string(tp.Text)
	if len(textStr) > tp.W-2 {
		textStr = textStr[:tp.W-2]
	}
	WriteAt(tp.X+1, tp.Y+1, ClrWhite+textStr+Reset)

	// Draw cursor if active
	if active {
		cursorX := tp.X + 1 + tp.Cursor
		if cursorX >= tp.X+tp.W-1 {
			cursorX = tp.X + tp.W - 2
		}
		WriteAt(cursorX, tp.Y+1, "")
	}
}

// InfoPanel implementations
func (ip *InfoPanel) Update(input byte) bool {
	// InfoPanel doesn't handle input directly
	return false
}

func (ip *InfoPanel) Draw(active bool) {
	// Draw base first
	ip.PanelBase.Draw(active)

	// Draw info lines
	for i, line := range ip.Lines {
		if i >= ip.H-2 { // leave space for borders
			break
		}
		y := ip.Y + 1 + i
		truncated := truncateToWidth(line, ip.W-2)
		WriteAt(ip.X+1, y, ClrWhite+truncated+Reset)
	}
}
