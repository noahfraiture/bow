package tui

import (
	"bytes"
	"strings"
	"unicode/utf8"
)

func truncateToWidth(s string, w int) string {
	if displayWidth(s) <= w {
		return s
	}
	stripped := stripANSI(s)
	r := []rune(stripped)
	var truncatedStripped string
	if len(r) <= w {
		truncatedStripped = stripped
	} else if w <= 2 {
		truncatedStripped = string(r[:w])
	} else {
		truncatedStripped = string(r[:w-2]) + ".."
	}
	// Re-add colors if present
	if strings.HasPrefix(s, "\x1b[7m") && strings.HasSuffix(s, "\x1b[0m") {
		return "\x1b[7m" + truncatedStripped + "\x1b[0m"
	} else if strings.HasPrefix(s, "\x1b[33m") && strings.HasSuffix(s, "\x1b[0m") {
		return "\x1b[33m" + truncatedStripped + "\x1b[0m"
	}
	return truncatedStripped
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

func stripANSI(s string) string {
	var result strings.Builder
	inEscape := false
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == '\x1b' {
			inEscape = true
		}
		if inEscape {
			if r == 'm' {
				inEscape = false
			}
		} else {
			result.WriteRune(r)
		}
		i += size
	}
	return result.String()
}

func displayWidth(s string) int {
	count := 0
	inEscape := false
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == '\x1b' {
			inEscape = true
		}
		if inEscape {
			if r == 'm' {
				inEscape = false
			}
		} else {
			count++
		}
		i += size
	}
	return count
}
