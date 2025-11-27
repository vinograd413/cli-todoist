package storage

import (
	"cliTodoist/internal/colors"
	"cliTodoist/internal/input"
	"cliTodoist/internal/table"
	"cliTodoist/internal/tasks"
	"cliTodoist/internal/util"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"go.etcd.io/bbolt"
)

func (d *DB) UpdateItem(i input.Input) (bool, error) {
	var input string
	var empty bool = false
	var numToDelete int
	var taskToUpdate *tasks.Task

	util.ClearScreenPlain()

	listOfTasks, err := d.GetAllTasks(i)
	if err != nil {
		return false, err
	}

	table := table.Table{Renderer: d.Renderer, TaskList: listOfTasks}

	if len(listOfTasks) == 0 {
		fmt.Print(util.HideCursor)
		fmt.Println(colors.Red + "    There is no open tasks, nothing to Update yet!" + colors.Reset + "\n")
		util.WaitForAnyKey(i, colors.Yellow+"Hit Any button to return to Menu"+colors.Reset)
		fmt.Print(util.ShowCursor)
		return false, nil
	} else {
		fmt.Println(colors.Yellow + colors.Bold + "Which task you want to update?" + colors.Reset)
		fmt.Println(colors.Gray + "Print" + colors.Yellow + " 'e' " + colors.Gray + "to navigate back to Menu\n" + colors.Reset)
		fmt.Println(colors.Yellow + colors.Bold + "Open tasks:" + colors.Reset)
		table.PrintTasksAsTable(i.File())
	}

	for {

		input, err := i.ReadLine(colors.Yellow + "Enter number: " + colors.Reset)
		if err != nil {
			return false, err
		}
		if inputLen := len(input); inputLen == 0 {
			if empty {
				for i := 0; i < 4; i++ {
					fmt.Printf("\033[2K\033[1A")
				}
			} else {
				for i := 0; i < 3; i++ {
					fmt.Printf("\033[2K\033[1A")
				}
			}
			fmt.Printf("\033[2K\r")
			fmt.Println(colors.Red +
				"Your should enter correct number!" +
				colors.Reset)
			empty = true
		} else if input == "e" {
			return false, nil
		} else {
			var inputErr error
			numToDelete, inputErr = strconv.Atoi(input)

			if inputErr != nil {
				if empty {
					for i := 0; i < 4; i++ {
						fmt.Printf("\033[2K\033[1A")
					}
				} else {
					for i := 0; i < 3; i++ {
						fmt.Printf("\033[2K\033[1A")
					}
				}
				fmt.Printf("\033[2K\r")
				fmt.Println(colors.Red + "You should enter a valid number" + colors.Reset)
				empty = true
			} else if len(listOfTasks) < numToDelete {
				if empty {
					for i := 0; i < 4; i++ {
						fmt.Printf("\033[2K\033[1A")
					}
				} else {
					for i := 0; i < 3; i++ {
						fmt.Printf("\033[2K\033[1A")
					}
				}
				fmt.Printf("\033[2K\r")
				fmt.Println(colors.Red + "You should enter only existing numbers of tasks" + colors.Reset)
				empty = true
			} else {
				break
			}
		}
	}

	taskToUpdate = listOfTasks[numToDelete-1]

	util.ClearScreenPlain()
	fmt.Println(colors.Yellow + colors.Bold + "Update task" + colors.Reset)
	fmt.Println(colors.Yellow + "Task header: " + colors.Reset + colors.Gray + taskToUpdate.Text + colors.Reset)

	for {
		fmt.Println(colors.Yellow + colors.Bold + "Print 'e' to navigate back to Menu" + colors.Reset)
		fmt.Print()

		input, err = i.ReadLine(colors.Yellow + "Enter new header: " + colors.Reset)
		if err != nil {
			return false, err
		}
		if input == "e" {
			return false, nil
		} else if inputLen := len(input); inputLen <= minHeaderLen {
			fmt.Println(colors.Red +
				fmt.Sprintf("Your header for task should be longer %d chars!", minHeaderLen) +
				colors.Reset)
		} else {
			break
		}
	}

	taskToUpdate.Text = input

	err = d.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(tasksBucket))
		if b == nil {
			return errors.New("Bucket does not exists")
		}

		buf, err := json.Marshal(taskToUpdate)
		if err != nil {
			return err
		}

		err = b.Put(taskToUpdate.ID, buf)

		return err
	})

	if err != nil {
		return false, err
	}

	fmt.Println(colors.Green + "Task updated successfully!")

	return util.AskRepeatOperation(i, "update")
}

func (d *DB) UpdateItemInteractive(i input.Input) (bool, error) {
	var input string
	var err error
	var taskToUpdate *tasks.Task

	util.ClearScreen()

	defer func() {
		fmt.Print(util.ShowCursor)
	}()

	listOfTasks, err := d.GetAllTasks(i)
	if err != nil {
		return false, err
	}

	if len(listOfTasks) == 0 {
		fmt.Print(util.HideCursor)
		fmt.Println(colors.Red + "    There is no open tasks, nothing to Update yet!" + colors.Reset + "\n")
		util.WaitForAnyKey(i, colors.Yellow+"Hit Any button to return to Menu"+colors.Reset)
		fmt.Print(util.ShowCursor)
		return false, nil
	} else {
		fmt.Println(colors.Yellow + colors.Bold + "Which task you want to update?" + colors.Reset)
		fmt.Println(util.NavigationPrompt)
	}

	table := table.Table{Renderer: d.Renderer, TaskList: listOfTasks}

	err = table.DisplayTableSingleSelection(i)
	if err != nil {
		return false, err
	}

	if len(table.SelectedTasks) == 0 {
		return false, nil
	}

	taskToUpdate = table.SelectedTasks[0]

	util.ClearScreenPlain()
	fmt.Println(colors.Yellow + colors.Bold + "Update task" + colors.Reset)
	fmt.Println(colors.Yellow + "Task header: " + colors.Reset + colors.Gray + taskToUpdate.Text + colors.Reset)

	for {
		fmt.Println(util.PromptBackToMenu)
		fmt.Print()

		input, err = i.ReadLine(colors.Yellow + "Enter new header: " + colors.Reset)
		if err != nil {
			return false, err
		}
		if input == "e" {
			return false, nil
		} else if inputLen := len(input); inputLen <= minHeaderLen {
			fmt.Println(colors.Red +
				fmt.Sprintf("Your header for task should be longer %d chars!", minHeaderLen) +
				colors.Reset)
		} else {
			break
		}
	}

	taskToUpdate.Text = input

	err = d.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(tasksBucket))
		if b == nil {
			return errors.New("Bucket does not exists")
		}

		buf, err := json.Marshal(taskToUpdate)
		if err != nil {
			return err
		}

		err = b.Put(taskToUpdate.ID, buf)

		return err
	})

	if err != nil {
		return false, err
	}

	fmt.Println(colors.Green + "Task updated successfully!")

	return util.AskRepeatOperation(i, "update")
}

func (d *DB) UpdateStatus(i input.Input) (bool, error) {
	var err error

	util.ClearScreen()

	defer func() {
		fmt.Print(util.ShowCursor)
	}()

	listOfTasks, err := d.GetAllTasks(i)
	if err != nil {
		return false, err
	}

	if len(listOfTasks) == 0 {
		fmt.Print(util.HideCursor)
		fmt.Println(colors.Red + "    There is no open tasks, nothing to Update yet!" + colors.Reset + "\n")
		util.WaitForAnyKey(i, colors.Yellow+"Hit Any button to return to Menu"+colors.Reset)
		fmt.Print(util.ShowCursor)
		return false, nil
	} else {
		fmt.Println(colors.Yellow + colors.Bold + "Which status you want to change?" + colors.Reset)
		fmt.Println(util.NavigationPromptUpdateStatus)
	}

	table := table.Table{Renderer: d.Renderer, TaskList: listOfTasks}

	err = table.DisplayTable(i)
	if err != nil {
		return false, err
	}

	if len(table.SelectedTasks) == 0 {
		return false, nil
	}

	for _, taskToUpdate := range table.SelectedTasks {
		taskToUpdate.IsCompleted = !taskToUpdate.IsCompleted
		taskToUpdate.CompletedAt = time.Now().Unix()

		err := d.db.Update(func(tx *bbolt.Tx) error {
			b := tx.Bucket([]byte(tasksBucket))
			if b == nil {
				return errors.New("Bucket does not exists")
			}

			buf, err := json.Marshal(taskToUpdate)
			if err != nil {
				return err
			}

			err = b.Put(taskToUpdate.ID, buf)

			return err
		})

		if err != nil {
			return false, err
		}
	}

	fmt.Println(colors.Green + "Task(s) updated successfully!")

	return util.AskRepeatOperation(i, "update")
}
