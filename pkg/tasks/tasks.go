package tasks

import (
	"golang.org/x/net/context"

	"cloud.google.com/go/datastore"
	"github.com/majordomusio/platform/pkg/store"
)

const (
	datastoreTasks string = "TASKS"
)

// taskDS holds information about a (periodic) task
type taskDS struct {
	Name    string
	Count   int
	LastRun int64
}

// Last retrieves the timestamp a task has last run
func Last(ctx context.Context, name string) int64 {
	key := datastore.NameKey(datastoreTasks, name, nil)
	var task taskDS

	err := store.Client().Get(ctx, key, &task)
	if err != nil {
		return 0
	}
	return task.LastRun
}

// Update updates the task with a timestamp when it has last run
func Update(ctx context.Context, name string, ts int64) error {
	key := datastore.NameKey(datastoreTasks, name, nil)
	var task taskDS

	err := store.Client().Get(ctx, key, &task)
	if err != nil {
		task = taskDS{
			name,
			1,
			ts,
		}
	} else {
		task.Count = task.Count + 1
		task.LastRun = ts
	}
	_, err = store.Client().Put(ctx, key, &task)
	return err
}
