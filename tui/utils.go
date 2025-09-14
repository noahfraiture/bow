package tui

import (
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
