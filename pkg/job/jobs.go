package job

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/taskqueue"

	"cloud.google.com/go/datastore"
	"github.com/lnkk-ai/lnkk/pkg/errorreporting"
	"github.com/lnkk-ai/lnkk/pkg/store"
)

const (
	datastoreJobs string = "JOBS"
)

// Job holds information about a (periodic) job
type Job struct {
	Name    string
	Count   int
	LastRun int64
}

// LastRun retrieves the timestamp a job has last run
func LastRun(ctx context.Context, name string) int64 {
	key := datastore.NameKey(datastoreJobs, name, nil)
	var job Job

	err := store.Client().Get(ctx, key, &job)
	if err != nil {
		return 0
	}
	return job.LastRun
}

// UpdateLastRun updates the job with a timestamp when it has last run
func UpdateLastRun(ctx context.Context, name string, ts int64) error {
	key := datastore.NameKey(datastoreJobs, name, nil)
	var job Job

	err := store.Client().Get(ctx, key, &job)
	if err != nil {
		job = Job{
			name,
			1,
			ts,
		}
	} else {
		job.Count = job.Count + 1
		job.LastRun = ts
	}
	_, err = store.Client().Put(ctx, key, &job)
	return err
}

// ScheduleJob creates a background job
func ScheduleJob(ctx context.Context, q, req string) {
	t := taskqueue.NewPOSTTask(req, nil)
	_, err := taskqueue.Add(ctx, t, q)
	if err != nil {
		errorreporting.Report(err)
	}
}
