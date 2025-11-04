package storage

import (
	"cliTodoist/internal/input"
	"cliTodoist/internal/table"
	"cliTodoist/internal/util"
	"encoding/json"
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"go.etcd.io/bbolt"
)

const dbFile = "dbFile.db"
const tasksBucket = "tasks"

type DB struct {
	db *bbolt.DB
}

func NewDB() (*DB, error) {
	db, err := bbolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	return &DB{db}, err
}

func (d *DB) GetAllTasks(input input.Input) ([]*Task, error) {
	var listOfTasks []*Task
	err := d.db.View(func(tx *bbolt.Tx) error {
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

	if err != nil {
		return nil, err
	}

	SortTasksByCreatedAt(listOfTasks)

	return listOfTasks, nil

}

func PrintTasksAsTable(file *os.File, renderer table.TableRenderer, tasks []*Task) error {
	header := []string{"#", "Text", "Created At", "Completed", "Completed At"}
	var rows [][]string
	for i, t := range tasks {
		var row []string
		row = append(row, strconv.Itoa(i+1))
		row = append(row, t.Text)
		row = append(row, time.Unix(t.CreatedAt, 0).Format(time.RFC1123))
		if !t.IsComleted {
			row = append(row, util.SymbGreenCross)
			row = append(row, "")
		}
		rows = append(rows, row)
	}

	return renderer.RenderTable(file, header, rows)
}
