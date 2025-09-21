package main

import (
	"app/tui"
	"fmt"
	"strings"
)

const (
	Reverse   = "\x1b[7m"
	reset     = "\x1b[0m"
	clrYellow = "\x1b[33m"
)

type ProcessListPanel struct {
	tui.ListPanel[Process]
	detailPanel *ProcessDetailPanel
}

func NewProcessListPanel(processes []Process, detail *ProcessDetailPanel) *ProcessListPanel {
	return &ProcessListPanel{
		ListPanel: tui.ListPanel[Process]{
			PanelBase: tui.PanelBase{
				Title:  "Processes",
				Border: true,
			},
			Items:    processes,
			Selected: 0,
		},
		detailPanel: detail,
	}
}

func (plp *ProcessListPanel) Update(msg tui.InputMessage) (handled bool, redraw bool) {
	oldSelected := plp.ListPanel.Selected
	handled, redraw = plp.ListPanel.Update(msg)
	if handled && redraw && plp.detailPanel != nil && plp.ListPanel.Selected != oldSelected {
		if plp.ListPanel.Selected < len(plp.ListPanel.Items) {
			plp.detailPanel.UpdateProcess(plp.ListPanel.Items[plp.ListPanel.Selected])
		}
	}
	return handled, redraw
}

func (plp *ProcessListPanel) Draw(active bool) string {
	var lines []string
	for i, item := range plp.ListPanel.Items {
		color := ""
		if i == plp.ListPanel.Selected {
			if active {
				color = Reverse
			} else {
				color = clrYellow
			}
		}
		lines = append(lines, fmt.Sprintf("%s%s%s", color, item.String(), reset))
	}
	return strings.Join(lines, "\n")
}

type ProcessDetailPanel struct {
	tui.InfoPanel
}

func NewProcessDetailPanel() *ProcessDetailPanel {
	return &ProcessDetailPanel{
		InfoPanel: tui.InfoPanel{
			PanelBase: tui.PanelBase{
				Title:  "Details",
				Border: true,
			},
			Lines: []string{"Select a process"},
		},
	}
}

func (pdp *ProcessDetailPanel) UpdateProcess(p Process) {
	pdp.InfoPanel.Lines = []string{
		fmt.Sprintf("PID: %d", p.PID),
		fmt.Sprintf("Name: %s", p.Name),
		fmt.Sprintf("CPU: %.1f%%", p.CPU),
		fmt.Sprintf("Memory: %.1f%%", p.Memory),
		fmt.Sprintf("Status: %s", p.Status),
	}
}
