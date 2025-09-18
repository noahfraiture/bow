package tui

import (
	"strconv"
	"strings"
)

// parseInput reads and parses input into an InputMessage
func (a *App) parseInput() (InputMessage, error) {
	b, err := a.term.reader.ReadByte()
	if err != nil {
		return InputMessage{}, err
	}

	raw := []byte{b}

	// Handle escape sequences
	if b == byte(KeyEsc) {
		return a.parseEscapeSequence(raw)
	}

	// Handle special keys
	switch b {
	case byte(KeyTab), byte(KeyEnter), byte(KeyBackspace), byte(KeySpace):
		return newKeyMessage(Key(b), raw), nil
	}

	// Handle printable characters
	if b >= 32 && b <= 126 {
		return newCharMessage(rune(b), raw), nil
	}

	// Handle control characters
	if b < 32 {
		msg := newCharMessage(rune(b+96), raw)
		msg.modifiers = append(msg.modifiers, ModCtrl)
		return msg, nil
	}

	// Default to special key
	return newKeyMessage(KeyEsc, raw), nil
}

// parseEscapeSequence handles escape sequences starting with ESC
func (a *App) parseEscapeSequence(raw []byte) (InputMessage, error) {
	next1, err := a.term.reader.ReadByte()
	if err != nil {
		return InputMessage{}, err
	}
	raw = append(raw, next1)

	if next1 == '[' {
		return a.parseCSISequence(raw)
	} else if next1 >= 'O' && next1 <= 'Z' {
		// Function keys F1-F4
		switch next1 {
		case 'P':
			return newKeyMessage(KeyF1, raw), nil
		case 'Q':
			return newKeyMessage(KeyF2, raw), nil
		case 'R':
			return newKeyMessage(KeyF3, raw), nil
		case 'S':
			return newKeyMessage(KeyF4, raw), nil
		}
	}
	// Unknown escape sequence, return as special key
	return newKeyMessage(KeyEsc, raw), nil
}

// parseCSISequence handles CSI (Control Sequence Introducer) sequences
func (a *App) parseCSISequence(raw []byte) (InputMessage, error) {
	var buf []byte
	var final byte
	for {
		b, err := a.term.reader.ReadByte()
		if err != nil {
			return InputMessage{}, err
		}
		raw = append(raw, b)
		if (b >= 'A' && b <= 'Z') || b == '~' {
			final = b
			break
		}
		buf = append(buf, b)
	}

	params := strings.Split(string(buf), ";")
	key, modifiers := a.parseCSIParams(params, final)
	if key != 0 {
		return newKeyMessageWithModifiers(key, raw, modifiers), nil
	}
	// Unknown sequence
	return newKeyMessage(KeyEsc, raw), nil
}

// parseCSIParams parses CSI parameters and returns the key and modifiers
func (a *App) parseCSIParams(params []string, final byte) (Key, []Modifier) {
	var key Key
	var modifiers []Modifier

	if len(params) > 1 {
		modifiers = parseModifiers(params[len(params)-1])
		params = params[:len(params)-1]
	}

	if final == '~' {
		if len(params) > 0 {
			keyCode, _ := strconv.Atoi(params[0])
			key = mapCSIKeyCode(keyCode)
		}
	} else {
		key, modifiers = mapCSIFinal(final, modifiers)
	}

	return key, modifiers
}

// parseModifiers parses modifier string and returns modifier slice
func parseModifiers(modStr string) []Modifier {
	var modifiers []Modifier
	modCode, _ := strconv.Atoi(modStr)
	if modCode&1 != 0 {
		modifiers = append(modifiers, ModShift)
	}
	if modCode&2 != 0 {
		modifiers = append(modifiers, ModAlt)
	}
	if modCode&4 != 0 {
		modifiers = append(modifiers, ModCtrl)
	}
	return modifiers
}

// mapCSIKeyCode maps CSI key codes to Key constants
func mapCSIKeyCode(keyCode int) Key {
	switch keyCode {
	case 1:
		return KeyHome
	case 2:
		return KeyInsert
	case 3:
		return KeyDelete
	case 4:
		return KeyEnd
	case 5:
		return KeyPageUp
	case 6:
		return KeyPageDown
	case 9:
		return KeyTab
	}
	return 0
}

// mapCSIFinal maps CSI final characters to Key constants and handles modifiers
func mapCSIFinal(final byte, modifiers []Modifier) (Key, []Modifier) {
	switch final {
	case 'A':
		return KeyUp, modifiers
	case 'B':
		return KeyDown, modifiers
	case 'C':
		return KeyRight, modifiers
	case 'D':
		return KeyLeft, modifiers
	case 'H':
		return KeyHome, modifiers
	case 'F':
		return KeyEnd, modifiers
	case 'Z':
		return KeyTab, append(modifiers, ModShift)
	}
	return 0, modifiers
}
