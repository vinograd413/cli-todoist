package storage

import (
	"cliTodoist/colors"
	"cliTodoist/internal/util"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go.etcd.io/bbolt"
)

func (d *DB) ListAllItems(i util.Input) {
	util.ClearScreen()

	defer func() {
		fmt.Print(util.ShowCursor)
	}()

	fmt.Println(colors.Yellow + colors.Bold + "List of all tasks:" + colors.Reset)

	d.db.View(func(tx *bbolt.Tx) error {
		var n int = 1
		b := tx.Bucket([]byte(tasksBucket))
		if b == nil {
			return errors.New("Bucket does not exists")
		}

		b.ForEach(func(k, v []byte) error {
			var t Task
			json.Unmarshal(v, &t)

			if !t.IsDeleted {
				fmt.Printf("    %d. %s	Created at: %v\n", n, t.Text, time.Unix(t.CreatedAt, 0).Format("2006-01-02 15:04:05"))
				n++
			}
			return nil
		})

		return nil
	})

	util.WaitForAnyKey(i, colors.Yellow+"Hit Any button to return to Menu"+colors.Reset)
	fmt.Print(util.ShowCursor)

}
