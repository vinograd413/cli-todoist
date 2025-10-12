package storage

import (
	"fmt"
	"time"
)

type Task struct {
	ID        []byte
	CreatedAt int64
	Text      string
	IsDeleted bool
}

func NewTask(text string) *Task {
	return &Task{
		[]byte(fmt.Sprintf("%s%d", text, time.Now().Unix())),
		time.Now().Unix(),
		text,
		false}
}
