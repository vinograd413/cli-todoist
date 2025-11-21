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
	"strings"

	"go.etcd.io/bbolt"
)

func (d *DB) DeleteItem(i input.Input) (bool, error) {
	var input string
	var err error
	var empty bool = false
	var numbersToDelete []int

	util.ClearScreenPlain()

	listOfTasks, err := d.GetAllTasks(i)
	if err != nil {
		return false, err
	}

	table := table.Table{Renderer: d.Renderer, TaskList: listOfTasks}

	if len(listOfTasks) == 0 {
		fmt.Print(util.HideCursor)
		fmt.Println(colors.Red + "    There is no open tasks, nothing to Delete yet!" + colors.Reset + "\n")
		util.WaitForAnyKey(i, colors.Yellow+"Hit Any button to return to Menu"+colors.Reset)
		fmt.Print(util.ShowCursor)
		return false, nil
	} else {
		fmt.Println(colors.Yellow + colors.Bold + "Which task you want to delete?" + colors.Reset)
		fmt.Println(colors.Gray +
			"You could select several numbers, just separate them by comma " +
			colors.Yellow + "(Like this: 1,2,3)" +
			colors.Reset)
		fmt.Println(colors.Gray + "Print" + colors.Yellow + " 'e' " + colors.Gray + "to navigate back to Menu\n" + colors.Reset)

		fmt.Println(colors.Yellow + colors.Bold + "Open tasks:" + colors.Reset)

		table.PrintTasksAsTable(i.File())

	}

	for {
		input, err = i.ReadLine(colors.Yellow + "Enter numbers: " + colors.Reset)
		if err != nil {
			return false, err
		}
		if inputLen := len(input); inputLen == 0 {
			if empty {
				for i := 0; i < 2; i++ {
					fmt.Printf("\033[2K\033[1A")
				}
			} else {
				fmt.Printf("\033[2K\033[1A")
			}
			fmt.Printf("\033[2K\r")
			fmt.Println(colors.Red +
				"Your should enter at least one number!" +
				colors.Reset)
			empty = true
		} else if input == "e" {
			return false, nil
		} else {
			var inputErr error
			var n int
			for v := range strings.SplitSeq(input, ",") {
				n, inputErr = strconv.Atoi(strings.TrimSpace(v))
				if inputErr == nil && len(listOfTasks) >= n {
					numbersToDelete = append(numbersToDelete, n)
				} else {
					break
				}
			}
			if inputErr != nil {
				if empty {
					for i := 0; i < 2; i++ {
						fmt.Printf("\033[2K\033[1A")
					}
				} else {
					fmt.Printf("\033[2K\033[1A")
				}
				fmt.Printf("\033[2K\r")
				fmt.Println(colors.Red + "You should enter a valid number" + colors.Reset)
				empty = true
			} else if len(listOfTasks) < n {
				if empty {
					for i := 0; i < 2; i++ {
						fmt.Printf("\033[2K\033[1A")
					}
				} else {
					fmt.Printf("\033[2K\033[1A")
				}
				fmt.Printf("\033[2K\r")
				fmt.Println(colors.Red + "You should enter only existing numbers of tasks" + colors.Reset)
				empty = true
			} else {
				break
			}
		}
	}

	var taskToDelete *tasks.Task
	for _, v := range numbersToDelete {
		taskToDelete = listOfTasks[v-1]
		taskToDelete.IsDeleted = true

		err := d.db.Update(func(tx *bbolt.Tx) error {
			b := tx.Bucket([]byte(tasksBucket))
			if b == nil {
				return errors.New("Bucket does not exists")
			}

			buf, err := json.Marshal(taskToDelete)
			if err != nil {
				return err
			}

			err = b.Put(taskToDelete.ID, buf)

			return err
		})

		if err != nil {
			return false, err
		}
	}

	fmt.Println(colors.Green + "Task deleted successfully!")

	return util.AskRepeatOperation(i, "delete")
}

func (d *DB) DeleteItemInteractive(i input.Input) (bool, error) {
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
		fmt.Println(colors.Red + "    There is no open tasks, nothing to Delete yet!" + colors.Reset + "\n")
		util.WaitForAnyKey(i, colors.Yellow+"Hit Any button to return to Menu"+colors.Reset)
		fmt.Print(util.ShowCursor)
		return false, nil
	} else {
		fmt.Println(colors.Yellow + colors.Bold + "Which task you want to delete?" + colors.Reset)
		fmt.Println(util.NavigationPromptDeletion)
	}

	table := table.Table{Renderer: d.Renderer, TaskList: listOfTasks}

	err = table.DisplayTable(i)
	if err != nil {
		return false, err
	}

	if len(table.SelectedTasks) == 0 {
		return false, nil
	}

	for _, taskToDelete := range table.SelectedTasks {
		taskToDelete.IsDeleted = true

		err := d.db.Update(func(tx *bbolt.Tx) error {
			b := tx.Bucket([]byte(tasksBucket))
			if b == nil {
				return errors.New("Bucket does not exists")
			}

			buf, err := json.Marshal(taskToDelete)
			if err != nil {
				return err
			}

			err = b.Put(taskToDelete.ID, buf)

			return err
		})

		if err != nil {
			return false, err
		}
	}

	fmt.Println(colors.Green + "Task(s) deleted successfully!")

	return util.AskRepeatOperation(i, "delete")
}
