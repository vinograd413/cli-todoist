package storage

import (
	"cliTodoist/colors"
	"cliTodoist/internal/input"
	"cliTodoist/internal/util"

	"encoding/json"
	"fmt"
	"strings"

	"go.etcd.io/bbolt"
)

const minHeaderLen = 3

func (d *DB) AddItem(i input.Input) (string, error) {
	var task *Task
	var input string

	util.ClearScreenPlain()
	fmt.Println(colors.Yellow + colors.Bold + "Add new task" + colors.Reset)
	fmt.Println(colors.Yellow + "First, enter header for your new TODO item" + colors.Reset)
	fmt.Println(colors.Yellow + colors.Bold + "Print 'e' to navigate back to Menu" + colors.Reset)
	for {
		input, err := i.ReadLine(colors.Yellow + "Enter header: " + colors.Reset)
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
			task = NewTask(input)
			break
		}
	}
	err := d.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(tasksBucket))
		if b == nil {
			var err error
			b, err = tx.CreateBucket([]byte(tasksBucket))
			if err != nil {
				return err
			}
		}

		buf, err := json.Marshal(task)
		if err != nil {
			return err
		}

		err = b.Put(task.ID, buf)

		return err
	})

	if err != nil {
		return "", err
	}

	fmt.Println(colors.Green + "Task added successfully!")

	for {
		input, err = i.ReadLine(colors.Yellow + "Do you want to add another task (Y/N): ")
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
