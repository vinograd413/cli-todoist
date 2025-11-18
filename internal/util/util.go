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

func RemoveFirst[T comparable](s []T, val T) []T {
	for i, v := range s {
		if v == val {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
