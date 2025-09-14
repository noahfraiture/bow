package tui

import (
	"strings"
	"unicode/utf8"
)

func truncateToWidth(s string, w int) string {
	if displayWidth(s) <= w {
		return s
	}
	var result strings.Builder
	visible := 0
	inEscape := false
	truncated := false
	for _, r := range s {
		if r == '\x1b' {
			inEscape = true
			result.WriteRune(r)
		} else if inEscape {
			result.WriteRune(r)
			if r == 'm' {
				inEscape = false
			}
		} else {
			if !truncated {
				if visible < w {
					result.WriteRune(r)
					visible++
					if visible == w-2 && w > 2 {
						result.WriteString("..")
						truncated = true
						visible += 2
					}
				} else {
					truncated = true
				}
			}
		}
	}
	return result.String()
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
