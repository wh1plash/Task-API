package fixtures

import (
	"context"
	"log"
	"task/store"
	"task/types"
	"time"
)

func AddTasks(store store.TaskStore, title, description, status string) *types.Task {
	task := types.Task{
		Title:       title,
		Description: description,
		Status:      status,
		CreatedAt:   uint64(time.Now().UnixNano()),
		UpdatedAt:   uint64(time.Now().UnixNano()),
	}

	insertedTask, err := store.InsertTask(context.Background(), &task)
	if err != nil {
		log.Fatal(err)
	}
	return insertedTask
}
