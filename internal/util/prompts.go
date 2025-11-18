package util

import "cliTodoist/internal/colors"

var NavigationPrompt string = colors.Gray +
	"\rUse " + colors.Yellow + "↑/↓" + colors.Gray + " to navigate, " +
	colors.Yellow + "Enter" + colors.Gray + " to select, " +
	colors.Yellow + "Esc" + colors.Gray + " to exit:\n" + colors.Reset

var NavigationPromptDeletion string = colors.Gray +
	"\rUse " + colors.Yellow + "↑/↓" + colors.Gray + " to navigate, " +
	colors.Yellow + "Space" + colors.Gray + " to select items, " +
	colors.Yellow + "Enter" + colors.Gray + " to delete selected tasks, " +
	colors.Yellow + "Esc" + colors.Gray + " to exit:\n" + colors.Reset

var PromptBackToMenu string = colors.Gray + "Print" + colors.Yellow + " 'e' " + colors.Gray + "to navigate back to Menu\n" + colors.Reset
