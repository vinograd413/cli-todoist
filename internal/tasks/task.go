package tasks

import (
	"fmt"
	"slices"
	"time"
)

type Task struct {
	ID          []byte
	Text        string
	CreatedAt   int64
	CompletedAt int64
	IsDeleted   bool
	IsCompleted bool
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
