package main

import (
	"app/tui"
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

type status int

const (
	_ status = iota
	NeedReview
	Draft
	RequestedChange
)

var stringToStatus = map[string]status{
	"Needs Review":    NeedReview,
	"Draft":           Draft,
	"Request Changes": RequestedChange,
}

func (s status) String() string {
	switch s {
	case NeedReview:
		return colorYellow + "Needs Review" + colorReset
	case Draft:
		return colorGreen + "Draft" + colorReset
	case RequestedChange:
		return colorRed + "Requested Changes" + colorReset
	default:
		panic("color not implemented")
	}

}

type diff struct {
	status  status
	id      string
	message string
}

func (d diff) String() string {
	return fmt.Sprintf("%-27s %s%s%s: %s", d.status.String(), colorYellow, d.id, colorReset, d.message)
}

type diffPanel struct {
	*tui.ListPanel[diff]
	diff *diff
}

func (c *diffPanel) Draw(_ bool) string {
	var buffer bytes.Buffer
	for i, item := range c.Items {
		selected := ""
		if c.Selected == i {
			selected = colorRed + "*" + colorReset
		}
		buffer.WriteString(fmt.Sprintf("%s %s\n", selected, item.String()))
	}
	return buffer.String()
}

func (c *diffPanel) Update(msg tui.InputMessage) bool {
	redraw := c.ListPanel.Update(msg)
	*c.diff = c.Items[c.Selected]
	return redraw
}

var diffRe = regexp.MustCompile(`(Needs Review|Draft).+(D\d{5}): (.*)`)

func parseDiff(line string) (diff, bool) {
	line = strings.TrimSpace(line)
	if line == "" {
		return diff{}, false
	}
	matches := diffRe.FindStringSubmatch(line)
	if len(matches) != 4 {
		return diff{}, false
	}
	statusStr := matches[1]
	id := matches[2]
	message := matches[3]
	status, ok := stringToStatus[statusStr]
	if !ok {
		return diff{}, false
	}
	return diff{status: status, id: id, message: message}, true
}

func getDiffTest() ([]diff, error) {
	return []diff{
		{
			status:  NeedReview,
			id:      "D12345",
			message: "message taht need review",
		},
		{
			status:  Draft,
			id:      "D67890",
			message: "still a draft",
		},
		{
			status:  RequestedChange,
			id:      "D07312",
			message: "you have  to do changes",
		},
	}, nil
}

func getDiff() ([]diff, error) {
	cmd := exec.Command("arc", "list")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run arc list: %w", err)
	}
	lines := strings.Split(string(output), "\n")
	var diffs []diff
	for _, line := range lines {
		if d, ok := parseDiff(line); ok {
			diffs = append(diffs, d)
		}
	}
	return diffs, nil
}

func newDiffPanel(name string, diffs []diff) diffPanel {
	return diffPanel{
		ListPanel: &tui.ListPanel[diff]{
			PanelBase: tui.PanelBase{
				Title:  name,
				Border: true,
			},
			Items: diffs,
		},
		diff: &diff{},
	}
}
