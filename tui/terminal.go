package tui

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type terminal struct {
	cols     int
	rows     int
	reader   *bufio.Reader
	prevStty string
}

func enableRawMode() (prev string, err error) {
	out, err := exec.Command("sh", "-c", "stty -g < /dev/tty").Output()
	if err != nil {
		return "", err
	}
	prev = strings.TrimSpace(string(out))
	if err := exec.Command("sh", "-c", "stty raw -echo < /dev/tty").Run(); err != nil {
		return prev, err
	}
	return prev, nil
}

func disableRawMode(prev string) {
	if prev != "" {
		_ = exec.Command("sh", "-c", "stty "+prev+" < /dev/tty").Run()
	} else {
		_ = exec.Command("sh", "-c", "stty sane < /dev/tty").Run()
	}
}

func getTermSize() (cols, rows int, err error) {
	out, err := exec.Command("sh", "-c", "stty size < /dev/tty").Output()
	if err != nil {
		return 80, 24, err
	}
	parts := strings.Fields(string(out))
	if len(parts) >= 2 {
		r, _ := strconv.Atoi(parts[0])
		c, _ := strconv.Atoi(parts[1])
		return c, r, nil
	}
	return 80, 24, fmt.Errorf("unexpected stty size output")
}

func writeAt(x, y int, s string) {
	fmt.Printf("\x1b[%d;%dH%s", y+1, x+1, s)
}

func clearScreen() {
	fmt.Print(Clear)
	fmt.Print(Home)
}
