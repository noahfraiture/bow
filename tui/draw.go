package tui

import (
	"fmt"
	"strings"
)

// Drawing
func (a *App) draw() {
	clearScreen()
	for i, p := range a.panels {
		active := i == a.activeIdx
		switch panel := p.(type) {
		case *ListPanel:
			a.drawPanelBase(panel.PanelBase, active)
			a.drawList(panel, active)
		case *TextPanel:
			a.drawPanelBase(panel.PanelBase, active)
			a.drawTextInput(panel, active)
		case *InfoPanel:
			a.drawPanelBase(panel.PanelBase, active)
			a.drawInfo(panel, active)
		}
	}
	status := " Tab: switch  •  ↑/↓: navigate  •  ←/→: move cursor  •  Enter: confirm  •  q/Ctrl-C: quit "
	writeAt(0, a.term.rows-1, padRightRuneString(status, a.term.cols))
	if tp, ok := a.panels[a.activeIdx].(*TextPanel); ok {
		x := tp.x + 1 + runeWidth(string(tp.Text[:tp.Cursor]))
		y := tp.y + 1
		if x >= tp.x+tp.w-1 {
			x = tp.x + tp.w - 2
		}
		writeAt(x, y, "")
		fmt.Print(ShowCursor)
	} else {
		fmt.Print(HideCursor)
	}
	fmt.Print(Reset)
}

func (a *App) drawPanelBase(b PanelBase, active bool) {
	color := ClrWhite
	if active {
		color = ClrCyan
	}
	if b.w < 2 || b.h < 2 {
		return
	}
	writeAt(b.x, b.y, color+"┌"+strings.Repeat("─", b.w-2)+"┐"+Reset)
	if b.Title != "" {
		title := " [" + b.Title + "] "
		tr := truncateToWidth(title, b.w-4)
		writeAt(b.x+2, b.y, color+Bold+tr+Reset)
	}
	for yy := 1; yy < b.h-1; yy++ {
		writeAt(b.x, b.y+yy, color+"│"+Reset)
		writeAt(b.x+b.w-1, b.y+yy, color+"│"+Reset)
		writeAt(b.x+1, b.y+yy, padRightRuneString("", b.w-2))
	}
	writeAt(b.x, b.y+b.h-1, color+"└"+strings.Repeat("─", b.w-2)+"┘"+Reset)
}

func (a *App) drawList(lp *ListPanel, active bool) {
	maxLines := lp.h - 2
	for i := 0; i < maxLines && i < len(lp.Items); i++ {
		item := lp.Items[i]
		prefix := "  "
		if i == lp.Selected {
			if active {
				writeAt(lp.x+1, lp.y+1+i, Reverse+"> "+truncateToWidth(item, lp.w-4)+Reset)
				continue
			} else {
				prefix = "• "
			}
		}
		writeAt(lp.x+1, lp.y+1+i, padRightRuneString(prefix+truncateToWidth(item, lp.w-4), lp.w-2))
	}
}

func (a *App) drawTextInput(tp *TextPanel, active bool) {
	content := string(tp.Text)
	visW := tp.w - 2
	display := content
	if runeWidth(content) > visW {
		runes := tp.Text
		if len(runes) > visW {
			display = string(runes[len(runes)-visW:])
		}
	}
	writeAt(tp.x+1, tp.y+1, padRightRuneString(display, visW))
	if active {
		cursorAbs := tp.Cursor
		total := len(tp.Text)
		start := 0
		if runeWidth(content) > visW {
			start = total - visW
		}
		cursorInDisplay := cursorAbs - start
		if cursorInDisplay < 0 {
			cursorInDisplay = 0
		}
		if cursorInDisplay > visW-1 {
			cursorInDisplay = visW - 1
		}
		var under string
		if start+cursorInDisplay < len(tp.Text) {
			under = string(tp.Text[start+cursorInDisplay])
		} else {
			under = " "
		}
		writeAt(tp.x+1+cursorInDisplay, tp.y+1, Reverse+under+Reset)
	}
}

func (a *App) drawInfo(ip *InfoPanel, active bool) {
	maxLines := ip.h - 2
	for i := 0; i < maxLines && i < len(ip.Lines); i++ {
		writeAt(ip.x+1, ip.y+1+i, padRightRuneString(truncateToWidth(ip.Lines[i], ip.w-2), ip.w-2))
	}
}
