package main

import (
	"app/tui"
	"fmt"
	"strings"
)

type MethodPanel struct {
	tui.ListPanel[string]
	req *Request
}

func NewMethodPanel(req *Request) *MethodPanel {
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	return &MethodPanel{
		ListPanel: tui.ListPanel[string]{
			PanelBase: tui.PanelBase{
				Title:  "Method",
				Border: true,
			},
			Items:    methods,
			Selected: 0,
		},
		req: req,
	}
}

func (mp *MethodPanel) Update(msg tui.InputMessage) (handled bool, redraw bool) {
	oldSelected := mp.ListPanel.Selected
	handled, redraw = mp.ListPanel.Update(msg)
	if handled && redraw && mp.ListPanel.Selected != oldSelected {
		mp.req.Method = mp.ListPanel.Items[mp.ListPanel.Selected]
	}
	return handled, redraw
}

type UrlPanel struct {
	tui.TextPanel
	req              *Request
	respHeadersPanel *ResponseHeadersPanel
	respBodyPanel    *ResponseBodyPanel
}

func NewUrlPanel(req *Request, rhp *ResponseHeadersPanel, rbp *ResponseBodyPanel) *UrlPanel {
	return &UrlPanel{
		TextPanel: tui.TextPanel{
			PanelBase: tui.PanelBase{
				Title:  "URL",
				Border: true,
			},
		},
		req:              req,
		respHeadersPanel: rhp,
		respBodyPanel:    rbp,
	}
}

func (up *UrlPanel) Update(msg tui.InputMessage) (handled bool, redraw bool) {
	handled, redraw = up.TextPanel.Update(msg)
	if handled && redraw {
		up.req.URL = string(up.TextPanel.Text)
	}
	return handled, redraw
}

type HeadersPanel struct {
	tui.TextPanel
	req *Request
}

func NewHeadersPanel(req *Request) *HeadersPanel {
	return &HeadersPanel{
		TextPanel: tui.TextPanel{
			PanelBase: tui.PanelBase{
				Title:  "Headers",
				Border: true,
			},
		},
		req: req,
	}
}

func (hp *HeadersPanel) Update(msg tui.InputMessage) (handled bool, redraw bool) {
	handled, redraw = hp.TextPanel.Update(msg)
	if handled && redraw {
		hp.req.Headers = parseHeaders(string(hp.TextPanel.Text))
	}
	return handled, redraw
}

func parseHeaders(text string) []Header {
	var headers []Header
	lines := strings.Split(text, ";")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			headers = append(headers, Header{Key: strings.TrimSpace(parts[0]), Value: strings.TrimSpace(parts[1])})
		}
	}
	return headers
}

type BodyPanel struct {
	tui.TextPanel
	req *Request
}

func NewBodyPanel(req *Request) *BodyPanel {
	return &BodyPanel{
		TextPanel: tui.TextPanel{
			PanelBase: tui.PanelBase{
				Title:  "Body",
				Border: true,
			},
		},
		req: req,
	}
}

func (bp *BodyPanel) Update(msg tui.InputMessage) (handled bool, redraw bool) {
	handled, redraw = bp.TextPanel.Update(msg)
	if handled && redraw {
		bp.req.Body = string(bp.TextPanel.Text)
	}
	return handled, redraw
}

type ResponseHeadersPanel struct {
	tui.InfoPanel
}

func NewResponseHeadersPanel() *ResponseHeadersPanel {
	return &ResponseHeadersPanel{
		InfoPanel: tui.InfoPanel{
			PanelBase: tui.PanelBase{
				Title:  "Response Headers",
				Border: true,
			},
			Lines: []string{"No response yet"},
		},
	}
}

func (rhp *ResponseHeadersPanel) UpdateResponse(resp Response) {
	if resp.Error != "" {
		rhp.InfoPanel.Lines = []string{"Error:", resp.Error}
		return
	}
	var lines []string
	for k, v := range resp.Headers {
		lines = append(lines, fmt.Sprintf("%s: %s", k, strings.Join(v, ", ")))
	}
	rhp.InfoPanel.Lines = lines
}

type ResponseBodyPanel struct {
	tui.InfoPanel
}

func NewResponseBodyPanel() *ResponseBodyPanel {
	return &ResponseBodyPanel{
		InfoPanel: tui.InfoPanel{
			PanelBase: tui.PanelBase{
				Title:  "Response Body",
				Border: true,
			},
			Lines: []string{"No response yet"},
		},
	}
}

func (rbp *ResponseBodyPanel) UpdateResponse(resp Response) {
	if resp.Error != "" {
		rbp.InfoPanel.Lines = []string{"Error:", resp.Error}
		return
	}
	body := resp.Body
	if len(body) > 1000 {
		body = body[:1000] + "..."
	}
	lines := strings.Split(body, "\n")
	if len(lines) > 20 {
		lines = lines[:20]
		lines = append(lines, "...")
	}
	rbp.InfoPanel.Lines = lines
}

type StatusPanel struct {
	tui.InfoPanel
}

func NewStatusPanel() *StatusPanel {
	return &StatusPanel{
		InfoPanel: tui.InfoPanel{
			PanelBase: tui.PanelBase{
				Title:  "Status",
				Border: true,
			},
			Lines: []string{"Ready"},
		},
	}
}

func (sp *StatusPanel) SetStatus(status string) {
	sp.InfoPanel.Lines = []string{status}
}
