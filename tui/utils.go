package tui

import (
	"bytes"
	"strings"
	"unicode/utf8"
)

// Utilities
func truncateToWidth(s string, w int) string {
	r := []rune(s)
	if len(r) <= w {
		return s
	}
	if w <= 2 {
		return string(r[:w])
	}
	return string(r[:w-2]) + ".."
}

func padRightRuneString(s string, w int) string {
	r := []rune(s)
	if len(r) >= w {
		return string(r[:w])
	}
	return string(r) + strings.Repeat(" ", w-len(r))
}

func wrapTextToWidth(s string, width int) []string {
	if width <= 0 {
		return []string{""}
	}
	words := strings.Fields(s)
	if len(words) == 0 {
		return []string{""}
	}
	var lines []string
	var cur bytes.Buffer
	for _, w := range words {
		if cur.Len()+len(w)+1 > width {
			lines = append(lines, padRightRuneString(cur.String(), width))
			cur.Reset()
		}
		if cur.Len() > 0 {
			cur.WriteByte(' ')
		}
		cur.WriteString(w)
	}
	if cur.Len() > 0 {
		lines = append(lines, padRightRuneString(cur.String(), width))
	}
	return lines
}

func runeWidth(s string) int {
	return utf8.RuneCountInString(s)
}
