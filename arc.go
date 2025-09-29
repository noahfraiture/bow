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
	NeedsReview
	Draft
	ChangesPlanned
	Accepted
	NeedsRevision
)

var stringToStatus = map[string]status{
	"Needs Review":    NeedsReview,
	"Draft":           Draft,
	"Changes Planned": ChangesPlanned,
	"Accepted":        Accepted,
	"Needs Revision":  NeedsRevision,
}

func (s status) String() string {
	switch s {
	case NeedsReview:
		return colorYellow + "Needs Review" + colorReset
	case Draft:
		return colorGreen + "Draft" + colorReset
	case ChangesPlanned:
		return colorRed + "Changes Planned" + colorReset
	case Accepted:
		return colorCyan + "Accepted" + colorReset
	case NeedsRevision:
		return colorRed + "Needs Revision" + colorReset
	default:
		panic(fmt.Sprintf("Unknown status: %d", s))
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

func (dp *diffPanel) Draw(_ bool) string {
	var buffer bytes.Buffer
	for i, item := range dp.Items {
		selected := ""
		if dp.Selected == i {
			selected = colorRed + "*" + colorReset
		}
		buffer.WriteString(fmt.Sprintf("%s %s\n", selected, item.String()))
	}
	return buffer.String()
}

func (dp *diffPanel) Update(msg tui.InputMessage) (handled bool, redraw bool) {
	handled, redraw = dp.ListPanel.Update(msg)
	if len(dp.Items) > 0 && dp.Selected >= 0 && dp.Selected < len(dp.Items) {
		*dp.diff = dp.Items[dp.Selected]
	}
	return handled, redraw
}

var diffRe = regexp.MustCompile(`(Needs Review|Draft|Changes Planned|Accepted|Needs Revision).+(D\d{5}): (.*)`)

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

func getDiff() ([]diff, error) {
	if isDevMode() {
		return []diff{{
			status:  2,
			id:      "1",
			message: "1",
		}, {
			status:  1,
			id:      "2",
			message: "2",
		}}, nil
	}
	cmd := exec.Command("arc", "list")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run 'arc list' command: %w", err)
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
