package tui

import "slices"

// KeyType represents the type of key input
type KeyType int

const (
	_ KeyType = iota
	KeyTypeChar
	KeyTypeKey
)

// Key represents special keys
type Key int

const (
	KeyUp        Key = 65 + iota // 65
	KeyDown                      // 66
	KeyRight                     // 67
	KeyLeft                      // 68
	KeyHome                      // 69
	KeyEnd                       // 70
	KeyPageUp                    // 71
	KeyPageDown                  // 72
	KeyF1                        // 73
	KeyF2                        // 74
	KeyF3                        // 75
	KeyF4                        // 76
	KeyF5                        // 77
	KeyF6                        // 78
	KeyF7                        // 79
	KeyF8                        // 80
	KeyF9                        // 81
	KeyF10                       // 82
	KeyF11                       // 83
	KeyF12                       // 84
	KeyInsert                    // 85
	KeyDelete                    // 86
	KeyTab       Key = 9
	KeyEnter     Key = 13
	KeyEsc       Key = 27
	KeyBackspace Key = 127
	KeySpace     Key = 32
)

// Modifier represents keyboard modifiers
type Modifier int

const (
	ModNone Modifier = iota
	ModCtrl
	ModAlt
	ModShift
)

// InputMessage represents a structured input event
type InputMessage struct {
	keyType   KeyType
	key       Key
	char      rune
	modifiers []Modifier
	raw       []byte
}

func newCharMessage(char rune, raw []byte) InputMessage {
	return InputMessage{
		keyType: KeyTypeChar,
		char:    char,
		raw:     raw,
	}
}

func newKeyMessage(key Key, raw []byte) InputMessage {
	return InputMessage{
		keyType: KeyTypeKey,
		key:     key,
		raw:     raw,
	}
}

// IsChar checks if the message is a character key
func (msg InputMessage) IsChar(char rune) bool {
	return msg.keyType == KeyTypeChar && msg.char == char
}

// IsKey checks if the message is a specific special key
func (msg InputMessage) IsKey(key Key) bool {
	return msg.keyType == KeyTypeKey && msg.key == key
}

// HasModifier checks if the message has a specific modifier
func (msg InputMessage) HasModifier(mod Modifier) bool {
	return slices.Contains(msg.modifiers, mod)
}
