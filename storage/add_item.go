package storage

import (
	"cliTodoist/internal/colors"
	"cliTodoist/internal/input"
	"cliTodoist/internal/tasks"
	"cliTodoist/internal/util"

	"encoding/json"
	"fmt"

	"go.etcd.io/bbolt"
)

const minHeaderLen = 3

func (d *DB) AddItem(i input.Input) (bool, error) {
	var task *tasks.Task

	util.ClearScreenPlain()
	fmt.Println(colors.Yellow + colors.Bold + "Add new task" + colors.Reset)
	fmt.Println(colors.Gray + "Enter header for your new TO DO item" + colors.Reset)
	fmt.Println(util.PromptBackToMenu)
	for {
		input, err := i.ReadLine(colors.Yellow + "Enter header: " + colors.Reset)
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
			task = tasks.NewTask(input)
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
		return false, err
	}

	fmt.Println(colors.Green + "Task added successfully!")

	return util.AskRepeatOperation(i, "add")
}
