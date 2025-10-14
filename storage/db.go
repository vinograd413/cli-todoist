package storage

import (
	"cliTodoist/internal/input"
	"encoding/json"
	"errors"
	"log"
	"slices"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
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

func (d *DB) ListAllTasks(input input.Input) (*tablewriter.Table, []*Task, error) {
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
		return nil, nil, err
	}

	table := tablewriter.NewWriter(input.File())
	table.Header([]string{"#", "Text", "Created At", "Completed", "Completed At"})

	SortTasksByCreatedAt(listOfTasks)

	for i, t := range listOfTasks {
		var row []string
		row = append(row, strconv.Itoa(i+1))
		row = append(row, t.Text)
		row = append(row, time.Unix(t.CreatedAt, 0).Format(time.RFC1123))
		if !t.IsComleted {
			row = append(row, "")
			row = append(row, "")
		}
		table.Append(row)
	}

	return table, listOfTasks, nil
}

func SortTasksByCreatedAt(tasks []*Task) {
	slices.SortStableFunc(tasks, func(a, b *Task) int {
		if a.CreatedAt < b.CreatedAt {
			return -1
		}
		if a.CreatedAt > b.CreatedAt {
			return 1
		}
		return 0
	})
}
