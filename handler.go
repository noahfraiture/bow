package main

import (
	"app/tui"
)

type handler struct {
	*tui.DefaultGlobalHandler
	panels         panels
	activeCommand  command
	rightPanel     *tui.Layout
	diffFromCommit *commit
	diffOnCommit   *commit
	diffToUpdate   *diff
	updateMsg      *string
	createMsg      *string
}

func (h *handler) GetStatus() string {
	return "HEY"
}

func (h *handler) OnPanelSwitch(app *tui.App, panelName string) {}

func (h *handler) UpdateGlobal(app *tui.App, msg tui.InputMessage) (handled bool, redraw bool) {
	switch {
	case msg.IsChar('u'):
		if h.activeCommand != Update {
			h.activeCommand = Update
			*h.rightPanel = &tui.VerticalSplit{
				Top:    &tui.PanelNode{Panel: &h.panels.diffs},
				Bottom: &tui.PanelNode{Panel: &h.panels.updateMsg},
			}
			redraw = true
		}
		handled = true
	case msg.IsChar('c'):
		if h.activeCommand != Create {
			h.activeCommand = Create
			*h.rightPanel = &tui.PanelNode{Panel: &h.panels.createMsg}
			redraw = true
		}
		handled = true
	default:
		return h.DefaultGlobalHandler.UpdateGlobal(app, msg)
	}
	return
}
