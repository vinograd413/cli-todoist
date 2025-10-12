package util

import (
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

func WaitForAnyKey(input Input, prompt string) {
	state, _ := input.SetRawMode()
	fmt.Print(prompt)
	input.ReadKey()
	input.RestoreMode(state)
}
