package storage

import (
	"cliTodoist/colors"
	"cliTodoist/internal/util"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.etcd.io/bbolt"
)

func (d *DB) UpdateItem(i util.Input) (string, error) {
	var input string
	var empty bool = false
	var numToDelete int
	var taskToUpdate *Task

	util.ClearScreenPlain()

	var listOfTasks []*Task
	d.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(tasksBucket))
		if b == nil {
			return errors.New("Bucket does not exists")
		}

		b.ForEach(func(k, v []byte) error {
			var t Task
			json.Unmarshal(v, &t)

			if !t.IsDeleted {
				listOfTasks = append(listOfTasks, &t)
			}
			return nil
		})
		return nil
	})

	if len(listOfTasks) == 0 {
		fmt.Print(util.HideCursor)
		fmt.Println(colors.Red + "    There is no open tasks, nothing to Update yet!" + colors.Reset + "\n")
		util.WaitForAnyKey(i, colors.Yellow+"Hit Any button to return to Menu"+colors.Reset)
		fmt.Print(util.ShowCursor)
		return "", nil
	} else {
		fmt.Println(colors.Yellow + colors.Bold + "Which task you want to update?" + colors.Reset)
		fmt.Println(colors.Yellow + colors.Bold + "Print 'e' to navigate back to Menu" + colors.Reset)
		fmt.Println(colors.Yellow + colors.Bold + "Open tasks:" + colors.Reset)
		for n, t := range listOfTasks {
			fmt.Printf("    %d. %s	Created at: %v\n", n+1, t.Text, time.Unix(t.CreatedAt, 0).Format("2006-01-02 15:04:05"))
		}
	}

	for {

		input, err := i.ReadLine(colors.Yellow + "Enter number: " + colors.Reset)
		if err != nil {
			return "", err
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
			return "", nil
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

		input, err := i.ReadLine(colors.Yellow + "Enter new header: " + colors.Reset)
		if err != nil {
			return "", err
		}
		if input == "e" {
			return "", nil
		} else if inputLen := len(input); inputLen <= minHeaderLen {
			fmt.Println(colors.Red +
				fmt.Sprintf("Your header for task should be longer %d chars!", minHeaderLen) +
				colors.Reset)
		} else {
			break
		}
	}

	taskToUpdate.Text = input

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
		return "", err
	}

	fmt.Println(colors.Green + "Task updated successfully!")

	for {

		input, err = i.ReadLine(colors.Yellow + "Do you want to update another task (Y/N): ")
		if err != nil {
			return "", err
		}
		if inputLen := len(input); inputLen != 0 {
			ls := strings.ToLower(input)
			answer := ls[:1]
			if answer == "y" {
				return "y", nil
			} else {
				return "n", nil
			}
		}
	}
}
