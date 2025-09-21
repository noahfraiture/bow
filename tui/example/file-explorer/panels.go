package main

import (
	"app/tui"
)

type fileListPanel struct {
	tui.ListPanel[fileItem]
	currentPath  string
	previewPanel *filePreviewPanel
}

func newFileListPanel(path string, preview *filePreviewPanel) *fileListPanel {
	items, _ := readDir(path)
	return &fileListPanel{
		ListPanel: tui.ListPanel[fileItem]{
			PanelBase: tui.PanelBase{
				Title:  "Files",
				Border: true,
			},
			Items:    items,
			Selected: 0,
		},
		currentPath:  path,
		previewPanel: preview,
	}
}

func (flp *fileListPanel) Update(msg tui.InputMessage) (handled bool, redraw bool) {
	oldSelected := flp.Selected
	handled, redraw = flp.ListPanel.Update(msg)
	if redraw && flp.Selected != oldSelected {
		if flp.Selected < len(flp.Items) {
			flp.previewPanel.Lines = getPreview(flp.Items[flp.Selected])
		}
	}

	if msg.IsKey(tui.KeyEnter) && flp.Selected < len(flp.Items) {
		selected := flp.Items[flp.Selected]
		if selected.isDir {
			flp.currentPath = selected.path
			newItems, err := readDir(flp.currentPath)
			if err == nil {
				flp.Items = newItems
				flp.Selected = 0
				if flp.previewPanel != nil && len(newItems) > 0 {
					flp.previewPanel.Lines = getPreview(newItems[0])
				}
			}
			return true, true
		}
	}
	return handled, redraw
}

type filePreviewPanel struct {
	tui.InfoPanel
}

func newFilePreviewPanel() *filePreviewPanel {
	return &filePreviewPanel{
		InfoPanel: tui.InfoPanel{
			PanelBase: tui.PanelBase{
				Title:  "Preview",
				Border: true,
			},
			Lines: []string{"Select a file or directory"},
		},
	}
}
