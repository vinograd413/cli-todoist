package util

const (
	ShowCursor         = "\033[?25h"
	HideCursor         = "\033[?25l"
	CursorUpFormat     = "\033[%dA"
	ClearAndHideCursor = "\033c\033[?25l"
	ClearNotHide       = "\033c"
	// Selections
	MoveUp   = byte(65)
	MoveDown = byte(66)
	Enter    = byte(13)
	Escape   = byte(27)
)

var NavigationKeys = map[byte]bool{
	MoveUp:   true,
	MoveDown: true,
}
