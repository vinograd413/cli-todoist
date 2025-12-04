package util

import "cliTodoist/internal/colors"

const (
	ShowCursor         = "\033[?25h"
	HideCursor         = "\033[?25l"
	CursorUpFormat     = "\033[%dA"
	ClearAndHideCursor = "\033c\033[?25l"
	ClearNotHide       = "\033c"
	ClearDown          = "\033[J"
	// Selections
	MoveUp   = byte(65)
	MoveDown = byte(66)
	Enter    = byte(13)
	Escape   = byte(27)
	Space    = byte(32)
	// Console window max height
	MinWindowHeight = 15
	// Max table height
	MaxTableHeight = 8
)

var NavigationKeys = map[byte]bool{
	MoveUp:   true,
	MoveDown: true,
}

const SymbGreenCross = "❎"
const SymbGreenCheck = "✅"

var CursorSelection = colors.Yellow + "➤  " + colors.Reset
var SymbRedCross = colors.Red + "✘" + colors.Reset
var SymbGreenSimleCheck = colors.Green + "✔" + colors.Reset
var CursorSelectionDeletion = colors.Yellow + "➤ " + SymbRedCross
