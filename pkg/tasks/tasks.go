package tasks

import (
	"fmt"

	"golang.org/x/net/context"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2beta3"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2beta3"

	"cloud.google.com/go/datastore"
	"github.com/lnkk-ai/lnkk/pkg/errorreporting"
	"github.com/lnkk-ai/lnkk/pkg/store"
	"github.com/majordomusio/commons/pkg/env"
)

const (
	datastoreTasks string = "TASKS"
)

var ctClient *cloudtasks.Client

func init() {
	ctx := context.Background()
	c, err := cloudtasks.NewClient(ctx)
	if err != nil {
		errorreporting.Report(err)
	}
	ctClient = c
}

// Close does the cleanup
func Close() {
	if ctClient == nil {
		return
	}
	ctClient.Close()
	ctClient = nil
}

// taskDS holds information about a (periodic) task
type taskDS struct {
	Name    string
	Count   int
	LastRun int64
}

// LastRun retrieves the timestamp a task has last run
func LastRun(ctx context.Context, name string) int64 {
	key := datastore.NameKey(datastoreTasks, name, nil)
	var task taskDS

	err := store.Client().Get(ctx, key, &task)
	if err != nil {
		return 0
	}
	return task.LastRun
}

// UpdateLastRun updates the task with a timestamp when it has last run
func UpdateLastRun(ctx context.Context, name string, ts int64) error {
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

// Schedule creates a background task
func Schedule(ctx context.Context, queue, req string) {
	// Build the Task queue path.
	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", env.Getenv("PROJECT_ID", ""), env.Getenv("LOCATION_ID", ""), queue)

	tq := &taskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &taskspb.Task{
			// https://godoc.org/google.golang.org/genproto/googleapis/cloud/tasks/v2beta3#HttpRequest
			PayloadType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					HttpMethod: taskspb.HttpMethod_POST,
					Url:        req,
				},
			},
		},
	}

	_, err := ctClient.CreateTask(ctx, tq)
	if err != nil {
		errorreporting.Report(err)
	}

}
