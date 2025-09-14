package tui

const (
	Esc        = "\x1b"
	Clear      = Esc + "[2J"
	Home       = Esc + "[H"
	HideCursor = Esc + "[?25l"
	ShowCursor = Esc + "[?25h"
	Reset      = Esc + "[0m"
	Bold       = Esc + "[1m"
	Reverse    = Esc + "[7m"

	ClrWhite  = Esc + "[37m"
	ClrCyan   = Esc + "[36m"
	ClrYellow = Esc + "[33m"
	ClrGreen  = Esc + "[32m"
)
