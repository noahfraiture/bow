package main

import (
	"app/tui"
)

func main() {
	currentPath := "."

	previewPanel := newFilePreviewPanel()
	listPanel := newFileListPanel(currentPath, previewPanel)

	if len(listPanel.Items) > 0 {
		previewPanel.Lines = getPreview(listPanel.Items[0])
	}

	layout := &tui.HorizontalSplit{
		Panels: []tui.Layout{
			&tui.PanelNode{Panel: listPanel, Weight: 2},
			&tui.PanelNode{Panel: previewPanel, Weight: 1},
		},
	}

	app := tui.NewApp(layout, nil)
	app.Run()
}
