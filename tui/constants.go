package tui

const (
	Esc        = "\x1b"
	Clear      = Esc + "[2J"
	Home       = Esc + "[H"
	HideCursor = Esc + "[?25l"
	ShowCursor = Esc + "[?25h"
	reset      = Esc + "[0m"
	bold       = Esc + "[1m"
	Reverse    = Esc + "[7m"

	clrWhite  = Esc + "[37m"
	clrCyan   = Esc + "[36m"
	clrYellow = Esc + "[33m"
	clrGreen  = Esc + "[32m"
)
