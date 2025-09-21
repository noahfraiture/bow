package main

import (
	"app/tui"
	"log"
)

func main() {
	processes, err := FetchProcesses()
	if err != nil {
		log.Fatal(err)
	}

	detailPanel := NewProcessDetailPanel()
	listPanel := NewProcessListPanel(processes, detailPanel)

	if len(processes) > 0 {
		detailPanel.UpdateProcess(processes[0])
	}

	layout := &tui.HorizontalSplit{
		Panels: []tui.Layout{
			&tui.PanelNode{Panel: listPanel, Weight: 2},
			&tui.PanelNode{Panel: detailPanel, Weight: 1},
		},
	}

	app := tui.NewApp(layout, nil)
	app.Run()
}
