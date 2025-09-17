package main

import "app/tui"

type messagePanel struct {
	*tui.TextPanel
	msg *string
}

func (mp *messagePanel) Update(msg tui.InputMessage) (handled bool, redraw bool) {
	handled, redraw = mp.TextPanel.Update(msg)
	*mp.msg = string(mp.Text)
	return handled, redraw
}

func newMessagePanel(name string) messagePanel {
	return messagePanel{
		TextPanel: &tui.TextPanel{
			PanelBase: tui.PanelBase{
				Title:  name,
				Border: true,
			},
			Text: []rune{},
		},
		msg: new(string),
	}
}
