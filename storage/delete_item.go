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

func (d *DB) DeleteItem(i util.Input) (string, error) {
	var input string
	var err error
	var empty bool = false
	var numbersToDelete []int

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
		fmt.Println(colors.Red + "    There is no open tasks, nothing to Delete yet!" + colors.Reset + "\n")
		util.WaitForAnyKey(i, colors.Yellow+"Hit Any button to return to Menu"+colors.Reset)
		fmt.Print(util.ShowCursor)
		return "", nil
	} else {
		fmt.Println(colors.Yellow + colors.Bold + "Which task you want to delete?" + colors.Reset)
		fmt.Println(colors.Yellow + "You could select several numbers, just separate them by comma (Like this: 1,2,3)" + colors.Reset)
		fmt.Println(colors.Yellow + colors.Bold + "Print 'e' to navigate back to Menu" + colors.Reset)
		fmt.Println(colors.Yellow + colors.Bold + "Open tasks:" + colors.Reset)
		for n, t := range listOfTasks {
			fmt.Printf("    %d. %s	Created at: %v\n", n, t.Text, time.Unix(t.CreatedAt, 0).Format("2006-01-02 15:04:05"))
		}
	}

	for {
		input, err = i.ReadLine(colors.Yellow + "Enter numbers: " + colors.Reset)
		if err != nil {
			return "", err
		}
		if inputLen := len(input); inputLen == 0 {
			if empty {
				for i := 0; i < 5; i++ {
					fmt.Printf("\033[2K\033[1A")
				}
			} else {
				for i := 0; i < 4; i++ {
					fmt.Printf("\033[2K\033[1A")
				}
			}
			fmt.Printf("\033[2K\r")
			fmt.Println(colors.Red +
				"Your should enter at least one number!" +
				colors.Reset)
			empty = true
		} else if input == "e" {
			return "", nil
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
					for i := 0; i < 5; i++ {
						fmt.Printf("\033[2K\033[1A")
					}
				} else {
					for i := 0; i < 4; i++ {
						fmt.Printf("\033[2K\033[1A")
					}
				}
				fmt.Printf("\033[2K\r")
				fmt.Println(colors.Red + "You should enter a valid number" + colors.Reset)
				empty = true
			} else if len(listOfTasks) < n {
				if empty {
					for i := 0; i < 5; i++ {
						fmt.Printf("\033[2K\033[1A")
					}
				} else {
					for i := 0; i < 4; i++ {
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

	var taskToDelete *Task
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
			return "", err
		}
	}

	fmt.Println(colors.Green + "Task deleted successfully!")

	for {

		input, err = i.ReadLine(colors.Yellow + "Do you want to delete another task (Y/N): ")
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
