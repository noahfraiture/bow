package main

import (
	"app/tui"
)

type PostmanHandler struct {
	req              *Request
	respHeadersPanel *ResponseHeadersPanel
	respBodyPanel    *ResponseBodyPanel
	statusPanel      *StatusPanel
}

func NewPostmanHandler(req *Request, rhp *ResponseHeadersPanel, rbp *ResponseBodyPanel, sp *StatusPanel) *PostmanHandler {
	return &PostmanHandler{
		req:              req,
		respHeadersPanel: rhp,
		respBodyPanel:    rbp,
		statusPanel:      sp,
	}
}

func (ph *PostmanHandler) UpdateGlobal(app *tui.App, msg tui.InputMessage) bool {
	if msg.IsChar('S') {
		ph.statusPanel.SetStatus("Sending...")
		resp := SendRequest(*ph.req)
		ph.respHeadersPanel.UpdateResponse(resp)
		ph.respBodyPanel.UpdateResponse(resp)
		if resp.Error != "" {
			ph.statusPanel.SetStatus("Error: " + resp.Error)
		} else {
			ph.statusPanel.SetStatus("Response received")
		}
		return true
	}
	return false
}

func (ph *PostmanHandler) OnPanelSwitch(app *tui.App, panelName string) {}

func (ph *PostmanHandler) GetStatus() string {
	return "Ready"
}

func main() {
	req := &Request{
		Method: "GET",
		URL:    "https://httpbin.org/get",
	}

	respHeadersPanel := NewResponseHeadersPanel()
	respBodyPanel := NewResponseBodyPanel()
	statusPanel := NewStatusPanel()

	methodPanel := NewMethodPanel(req)
	urlPanel := NewUrlPanel(req, respHeadersPanel, respBodyPanel)
	headersPanel := NewHeadersPanel(req)
	bodyPanel := NewBodyPanel(req)

	// Layout: Vertical split
	// Top: Horizontal for method, url, headers
	// Middle: Body
	// Bottom: Horizontal for resp headers, body

	topSplit := &tui.HorizontalSplit{
		Panels: []tui.Layout{
			&tui.PanelNode{Panel: methodPanel, Weight: 1},
			&tui.PanelNode{Panel: urlPanel, Weight: 3},
			&tui.PanelNode{Panel: headersPanel, Weight: 2},
		},
	}

	bottomSplit := &tui.HorizontalSplit{
		Panels: []tui.Layout{
			&tui.PanelNode{Panel: respHeadersPanel, Weight: 1},
			&tui.PanelNode{Panel: respBodyPanel, Weight: 2},
		},
	}

	layout := &tui.VerticalSplit{
		Panels: []tui.Layout{
			topSplit,
			&tui.PanelNode{Panel: bodyPanel, Weight: 2},
			bottomSplit,
			&tui.PanelNode{Panel: statusPanel, Weight: 1},
		},
	}

	handler := NewPostmanHandler(req, respHeadersPanel, respBodyPanel, statusPanel)
	app := tui.NewApp(layout, handler)
	app.Run()
}
