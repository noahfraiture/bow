package main

import (
	"app/tui"
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/object"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
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
}

func (c commitPanel) Draw(_ bool) string {
	var buffer bytes.Buffer
	for i, item := range c.Items {
		if c.Selected == i {
			buffer.WriteString(fmt.Sprintf("%s*%s %s\n", colorRed, colorReset, item.String()))
		} else {
			buffer.WriteString(fmt.Sprintf("%s\n", item.String()))
		}
	}
	return buffer.String()
}

func createApp() (*tui.App, error) {
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

	diffFrom := commitPanel{
		ListPanel: &tui.ListPanel[commit]{
			PanelBase: tui.PanelBase{
				Title:  "Diff From",
				Border: true,
			},
			Items: commits,
		}}

	diffOn := commitPanel{
		ListPanel: &tui.ListPanel[commit]{
			PanelBase: tui.PanelBase{
				Title:  "Diff on",
				Border: true,
			},
			Items: commits,
		}}

	app := tui.NewApp(&tui.HorizontalSplit{
		Left:  &tui.PanelNode{Panel: diffFrom},
		Right: &tui.PanelNode{Panel: diffOn},
	})

	return app, nil
}

func main() {
	app, err := createApp()
	if err != nil {
		log.Fatal(err)
	}

	app.Run()
}
