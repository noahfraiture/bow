package tui

// Input handling
func (a *App) handleByte(b byte) {
	if b == KeyEsc {
		next1, err := a.term.reader.ReadByte()
		if err != nil {
			return
		}
		if next1 == '[' {
			next2, err := a.term.reader.ReadByte()
			if err != nil {
				return
			}
			switch next2 {
			case 'A':
				a.onUp()
			case 'B':
				a.onDown()
			case 'C':
				a.onRight()
			case 'D':
				a.onLeft()
			}
			return
		}
		return
	}

	switch b {
	case KeyTab:
		a.switchPanel()
	case KeyEnter, '\n':
		a.onEnter()
	case KeyBackspace, 8:
		a.onBackspace()
	default:
		if b >= 32 && b <= 126 {
			a.onRune(rune(b))
		}
	}
}

func (a *App) switchPanel() {
	a.activeIdx = (a.activeIdx + 1) % len(a.panels)
}

func (a *App) onUp() {
	if lp, ok := a.panels[a.activeIdx].(*ListPanel); ok {
		if lp.Selected > 0 {
			lp.Selected--
		}
	}
}

func (a *App) onDown() {
	if lp, ok := a.panels[a.activeIdx].(*ListPanel); ok {
		if lp.Selected < len(lp.Items)-1 {
			lp.Selected++
		}
	}
}

func (a *App) onLeft() {
	if tp, ok := a.panels[a.activeIdx].(*TextPanel); ok {
		if tp.Cursor > 0 {
			tp.Cursor--
		}
	}
}

func (a *App) onRight() {
	if tp, ok := a.panels[a.activeIdx].(*TextPanel); ok {
		if tp.Cursor < len(tp.Text) {
			tp.Cursor++
		}
	}
}

func (a *App) onEnter() {
	switch p := a.panels[a.activeIdx].(type) {
	case *ListPanel:
		// echo to info if there's an info panel
		for _, panel := range a.panels {
			if ip, ok := panel.(*InfoPanel); ok {
				ip.Lines = append(ip.Lines, "", "Selected: "+p.Items[p.Selected])
			}
		}
	case *TextPanel:
		for _, panel := range a.panels {
			if ip, ok := panel.(*InfoPanel); ok {
				ip.Lines = append(ip.Lines, "", "Input: "+string(p.Text))
			}
		}
		p.Text = []rune{}
		p.Cursor = 0
	}
}

func (a *App) onBackspace() {
	if tp, ok := a.panels[a.activeIdx].(*TextPanel); ok {
		if tp.Cursor > 0 && len(tp.Text) > 0 {
			i := tp.Cursor
			tp.Text = append(tp.Text[:i-1], tp.Text[i:]...)
			tp.Cursor--
		}
	}
}

func (a *App) onRune(r rune) {
	if tp, ok := a.panels[a.activeIdx].(*TextPanel); ok {
		i := tp.Cursor
		before := tp.Text[:i]
		after := tp.Text[i:]
		newText := make([]rune, 0, len(before)+1+len(after))
		newText = append(newText, before...)
		newText = append(newText, r)
		newText = append(newText, after...)
		tp.Text = newText
		tp.Cursor++
	} else {
		switch r {
		case 'j':
			a.onDown()
		case 'k':
			a.onUp()
		case 'h':
			a.onLeft()
		case 'l':
			a.onRight()
		case 'q', 'Q':
			a.running = false
		}
	}
}
