package main

import (
	"app/tui"
)

type messagePanel struct {
	*tui.TextPanel
	msg *string
}

func (mp *messagePanel) Update(msg tui.InputMessage) (handled bool, redraw bool) {
	handled, redraw = mp.TextPanel.Update(msg)
	*mp.msg = string(mp.Text)
	return handled, redraw
}

func newMessagePanelUpdate(name string) messagePanel {
	return messagePanel{
		TextPanel: &tui.TextPanel{
			PanelBase: tui.PanelBase{
				Title:  name,
				Border: true,
			},
			Text: make([]rune, 0),
		},
		msg: new(string),
	}
}

func newMessagePanelCreate(name string) messagePanel {
	msg := "\n\nSummary: \n\nTest case: \n\nReviewers: \n\n"
	return messagePanel{
		TextPanel: &tui.TextPanel{
			PanelBase: tui.PanelBase{
				Title:  name,
				Border: true,
			},
			Text: []rune(msg),
		},
		msg: &msg,
	}
}
