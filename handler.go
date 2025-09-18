package main

import (
	"app/tui"
	"os/exec"
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

func (h *handler) UpdateGlobal(app *tui.App, msg tui.InputMessage) (redraw bool) {
	switch {
	case msg.IsChar('u'):
		if h.activeCommand != Update {
			h.activeCommand = Update
			*h.rightPanel = &tui.VerticalSplit{
				Top:    &tui.PanelNode{Panel: &h.panels.diffs},
				Bottom: &tui.PanelNode{Panel: &h.panels.updateMsg},
			}
			return true
		}
	case msg.IsChar('c'):
		if h.activeCommand != Create {
			h.activeCommand = Create
			*h.rightPanel = &tui.PanelNode{Panel: &h.panels.createMsg}
			return true
		}
	case msg.HasModifier(tui.ModCtrl) && msg.IsChar('s'):
		h.runUpdate()
		app.Stop()
	default:
		return h.DefaultGlobalHandler.UpdateGlobal(app, msg)
	}
	return false
}

func (h *handler) runUpdate() error {
	cmd := exec.Command(
		"arc",
		"diff", h.diffFromCommit.Hash.String(),
		"--head", h.diffOnCommit.Hash.String(),
		"--update", h.diffToUpdate.id,
		"--message", *h.updateMsg,
	)
	return cmd.Run()
}
