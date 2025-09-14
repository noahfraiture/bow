package main

import (
	"app/tui"
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/object"
)

type commit struct {
	*object.Commit
}

func (c commit) String() string {
	msg := strings.TrimRight(strings.SplitN(c.Message, "\n", 2)[0], "\n")
	return fmt.Sprintf("%s%s%s: %s", colorYellow, c.Hash.String()[:6], colorReset, msg)
}

type commitPanel struct {
	*tui.ListPanel[commit]
	commit *commit
}

func (cp *commitPanel) Draw(_ bool) string {
	var buffer bytes.Buffer
	for i, item := range cp.Items {
		selected := ""
		if cp.Selected == i {
			selected = colorRed + "*" + colorReset
		}
		buffer.WriteString(fmt.Sprintf("%s %s\n", selected, item.String()))
	}
	return buffer.String()
}

func (cp *commitPanel) Update(msg tui.InputMessage) bool {
	redraw := cp.ListPanel.Update(msg)
	if len(cp.Items) > 0 && cp.Selected >= 0 && cp.Selected < len(cp.Items) {
		*cp.commit = cp.Items[cp.Selected]
	}
	return redraw
}

func getCommits() ([]commit, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current working directory: %w", err)
	}

	repo, err := git.PlainOpenWithOptions(dir, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return nil, fmt.Errorf("failed to open git repository at %s: %w", dir, err)
	}

	commitsIter, err := repo.Log(&git.LogOptions{Order: git.LogOrderCommitterTime})
	if err != nil {
		return nil, fmt.Errorf("failed to get commit log: %w", err)
	}
	commits := []commit{}
	count := 0
	err = commitsIter.ForEach(func(c *object.Commit) error {
		if count >= 20 {
			return nil // Stop after 20 commits
		}
		commits = append(commits, commit{c})
		count++
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to iterate commits: %w", err)
	}
	return commits, nil
}

func newCommitPanel(name string, commits []commit) commitPanel {
	return commitPanel{
		ListPanel: &tui.ListPanel[commit]{
			PanelBase: tui.PanelBase{
				Title:  name,
				Border: true,
			},
			Items: commits,
		},
		commit: &commit{},
	}
}
