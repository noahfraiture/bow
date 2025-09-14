package tui

// KeyType represents the type of key input
type KeyType int

const (
	KeyTypeChar KeyType = iota
	KeyTypeArrow
	KeyTypeNavigation
	KeyTypeFunction
	KeyTypeSpecial
)

// Key constants for common keys
const (
	KeyUp    = 65
	KeyDown  = 66
	KeyRight = 67
	KeyLeft  = 68

	KeyHome     = 72
	KeyEnd      = 70
	KeyPageUp   = 53
	KeyPageDown = 54

	KeyF1  = 80
	KeyF2  = 81
	KeyF3  = 82
	KeyF4  = 83
	KeyF5  = 84
	KeyF6  = 85
	KeyF7  = 86
	KeyF8  = 87
	KeyF9  = 88
	KeyF10 = 89
	KeyF11 = 90
	KeyF12 = 91

	KeyInsert = 50
	KeyDelete = 51
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
	Key       KeyType
	Code      int  // Key code for special keys
	Char      rune // Character for printable keys
	Modifiers []Modifier
	Raw       []byte // Original raw bytes
}

// NewCharMessage creates a new character input message
func NewCharMessage(char rune, raw []byte) InputMessage {
	return InputMessage{
		Key:  KeyTypeChar,
		Char: char,
		Raw:  raw,
	}
}

// NewArrowMessage creates a new arrow key input message
func NewArrowMessage(code int, raw []byte) InputMessage {
	return InputMessage{
		Key:  KeyTypeArrow,
		Code: code,
		Raw:  raw,
	}
}

// NewNavigationMessage creates a new navigation key input message
func NewNavigationMessage(code int, raw []byte) InputMessage {
	return InputMessage{
		Key:  KeyTypeNavigation,
		Code: code,
		Raw:  raw,
	}
}

// NewSpecialMessage creates a new special key input message
func NewSpecialMessage(code int, raw []byte) InputMessage {
	return InputMessage{
		Key:  KeyTypeSpecial,
		Code: code,
		Raw:  raw,
	}
}

// IsChar checks if the message is a character key
func (msg InputMessage) IsChar(char rune) bool {
	return msg.Key == KeyTypeChar && msg.Char == char
}

// IsArrow checks if the message is an arrow key
func (msg InputMessage) IsArrow(code int) bool {
	return msg.Key == KeyTypeArrow && msg.Code == code
}

// IsNavigation checks if the message is a navigation key
func (msg InputMessage) IsNavigation(code int) bool {
	return msg.Key == KeyTypeNavigation && msg.Code == code
}

// IsSpecial checks if the message is a special key
func (msg InputMessage) IsSpecial(code int) bool {
	return msg.Key == KeyTypeSpecial && msg.Code == code
}

// HasModifier checks if the message has a specific modifier
func (msg InputMessage) HasModifier(mod Modifier) bool {
	for _, m := range msg.Modifiers {
		if m == mod {
			return true
		}
	}
	return false
}
