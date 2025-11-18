package util

import (
	"cliTodoist/internal/colors"
	"cliTodoist/internal/input"
	"strings"
)

func AskRepeatOperation(i input.Input, action string) (bool, error) {
	for {
		response, err := i.ReadLine(colors.Yellow + "Do you want to delete another task (Y/N): ")
		if err != nil {
			return false, err
		}
		if len(response) > 0 && strings.ToLower(response[:1]) == "y" {
			return true, nil
		} else if len(response) > 0 && strings.ToLower(response[:1]) == "n" {
			return false, nil
		}
	}
}
