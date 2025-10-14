package storage

import (
	"fmt"
	"time"
)

type Task struct {
	ID         []byte
	Text       string
	CreatedAt  int64
	ComletedAt int64
	IsDeleted  bool
	IsComleted bool
}

func NewTask(text string) *Task {
	return &Task{
		[]byte(fmt.Sprintf("%s%d", text, time.Now().Unix())),
		text,
		time.Now().Unix(),
		-1,
		false,
		false}
}
