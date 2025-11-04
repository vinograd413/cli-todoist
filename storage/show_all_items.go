package storage

import (
	"cliTodoist/colors"
	"cliTodoist/internal/input"
	"cliTodoist/internal/table"
	"cliTodoist/internal/util"
	"fmt"
)

func (d *DB) ShowAllItems(i input.Input) error {
	util.ClearScreen()
	renderer := table.AquaSecTableRenderer{}

	defer func() {
		fmt.Print(util.ShowCursor)
	}()

	fmt.Println(colors.Yellow + colors.Bold + "List of all tasks:" + colors.Reset)

	tasks, err := d.GetAllTasks(i)
	if err != nil {
		return err
	}

	if len(tasks) > 0 {
		PrintTasksAsTable(i.File(), &renderer, tasks)
	} else {
		fmt.Println(colors.Red + "    There is no open tasks yet!" + colors.Reset + "\n")
	}

	util.WaitForAnyKey(i, colors.Yellow+"Hit Any button to return to Menu"+colors.Reset)
	fmt.Print(util.ShowCursor)
	return nil
}
