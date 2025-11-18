package storage

import (
	"cliTodoist/internal/input"
	"cliTodoist/internal/renderer"
	"cliTodoist/internal/tasks"
	"encoding/json"
	"errors"
	"log"

	"go.etcd.io/bbolt"
)

const dbFile = "dbFile.db"
const tasksBucket = "tasks"

type DB struct {
	db       *bbolt.DB
	Renderer renderer.TableRenderer
}

func NewDB(renderer renderer.TableRenderer) (*DB, error) {
	db, err := bbolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	return &DB{db: db, Renderer: renderer}, err
}

func (d *DB) GetAllTasks(input input.Input) ([]*tasks.Task, error) {
	var listOfTasks []*tasks.Task
	err := d.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(tasksBucket))
		if b == nil {
			return errors.New("Bucket does not exists")
		}

		b.ForEach(func(k, v []byte) error {
			var t tasks.Task
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

	tasks.SortTasksByCreatedAt(listOfTasks)

	return listOfTasks, nil

}
