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
	msg := strings.TrimRight(c.Message, "\n")
	return fmt.Sprintf("%s%s%s: %s", colorYellow, c.Hash.String()[:6], colorReset, msg)
}

type commitPanel struct {
	*tui.ListPanel[commit]
	commit *commit
}

func (c *commitPanel) Draw(_ bool) string {
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

func (c *commitPanel) Update(msg tui.InputMessage) bool {
	redraw := c.ListPanel.Update(msg)
	*c.commit = c.Items[c.Selected]
	return redraw
}

func getCommits() ([]commit, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	repo, err := git.PlainOpen(dir)
	if err != nil {
		return nil, err
	}

	commitsIter, err := repo.CommitObjects()
	if err != nil {
		return nil, err
	}
	commits := []commit{}
	err = commitsIter.ForEach(func(c *object.Commit) error {
		commits = append(commits, commit{c})
		return nil
	})
	if err != nil {
		return nil, err
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
