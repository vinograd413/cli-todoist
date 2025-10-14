package util

import (
	"cliTodoist/internal/input"
	"fmt"
)

// Clear screen and hide cursor
func ClearScreen() {
	fmt.Print(ClearAndHideCursor)
}

// Clear screen without hiding cursor
func ClearScreenPlain() {
	fmt.Print(ClearNotHide)
}

func WaitForAnyKey(input input.Input, prompt string) {
	state, _ := input.SetRawMode()
	fmt.Print(prompt)
	input.ReadKey()
	input.RestoreMode(state)
}
