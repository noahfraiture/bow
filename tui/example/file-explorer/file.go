package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type fileItem struct {
	name  string
	path  string
	isDir bool
}

func (fi fileItem) String() string {
	prefix := "[FILE]"
	if fi.isDir {
		prefix = "[DIR] "
	}
	return fmt.Sprintf("%s %s", prefix, fi.name)
}

func readDir(path string) ([]fileItem, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var items []fileItem
	for _, entry := range entries {
		if entry.Name() == "." || entry.Name() == ".." {
			continue
		}
		fullPath := filepath.Join(path, entry.Name())
		items = append(items, fileItem{
			name:  entry.Name(),
			path:  fullPath,
			isDir: entry.IsDir(),
		})
	}
	return items, nil
}

func getPreview(item fileItem) []string {
	if item.isDir {
		subItems, err := readDir(item.path)
		if err != nil {
			return []string{"Error reading directory:", err.Error()}
		}
		var lines []string
		for i, sub := range subItems {
			if i >= 20 {
				break
			}
			lines = append(lines, sub.String())
		}
		if len(subItems) > 20 {
			lines = append(lines, "...")
		}
		return lines
	} else {
		content, err := os.ReadFile(item.path)
		if err != nil {
			return []string{"Error reading file:", err.Error()}
		}
		text := string(content)
		if len(text) > 1000 {
			text = text[:1000] + "..."
		}
		lines := strings.Split(text, "\n")
		if len(lines) > 20 {
			lines = lines[:20]
			lines = append(lines, "...")
		}
		return lines
	}
}
