package storage

import (
	"cliTodoist/colors"
	"cliTodoist/internal/input"
	"cliTodoist/internal/util"
	"fmt"
)

func (d *DB) ShowAllItems(i input.Input) error {
	util.ClearScreen()

	defer func() {
		fmt.Print(util.ShowCursor)
	}()

	fmt.Println(colors.Yellow + colors.Bold + "List of all tasks:" + colors.Reset)

	table, _, err := d.ListAllTasks(i)
	if err != nil {
		return err
	}

	if table != nil {
		table.Render()
	} else {
		fmt.Println(colors.Red + "    There is no open tasks yet!" + colors.Reset + "\n")
	}

	util.WaitForAnyKey(i, colors.Yellow+"Hit Any button to return to Menu"+colors.Reset)
	fmt.Print(util.ShowCursor)
	return nil
}
