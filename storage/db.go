package storage

import (
	"log"

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
