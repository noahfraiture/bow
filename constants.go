package main

import "os"

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
)

func isDevMode() bool {
	return os.Getenv("BOW_DEV") == "1"
}
