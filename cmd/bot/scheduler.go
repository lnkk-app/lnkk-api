package main

import (
	"fmt"

	"cloud.google.com/go/datastore"
	"google.golang.org/appengine"

	"github.com/gin-gonic/gin"

	"github.com/majordomusio/commons/pkg/env"
	"github.com/majordomusio/commons/pkg/util"

	"github.com/majordomusio/platform/pkg/errorreporting"
	"github.com/majordomusio/platform/pkg/store"
	"github.com/majordomusio/platform/pkg/tasks"

	"github.com/lnkk-ai/lnkk/pkg/api"

	"github.com/lnkk-ai/lnkk/internal/backend"
	"github.com/lnkk-ai/lnkk/internal/types"
)

// scheduleUpdateWorkspaces schedules all workspaces that need updating
func scheduleUpdateWorkspaces(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	now := util.Timestamp()
	var workspaces []types.WorkspaceDS

	q := datastore.NewQuery(backend.DatastoreWorkspaces).Filter("Next <", now)
	_, err := store.Client().GetAll(ctx, q, &workspaces)

	if err == nil {
		for i := range workspaces {

			tasks.Schedule(ctx, backend.BackgroundWorkQueue, env.Getenv("BASE_URL", "")+api.JobsBaseURL+"/j/users?id="+workspaces[i].ID)
			tasks.Schedule(ctx, backend.BackgroundWorkQueue, env.Getenv("BASE_URL", "")+api.JobsBaseURL+"/j/channels?id="+workspaces[i].ID)

			backend.MarkWorkspaceUpdated(ctx, workspaces[i].ID)
		}
	} else {
		errorreporting.Report(err)
	}
}

// scheduleCollectMessages schedules the collection of messages in a given team & channel
func scheduleCollectMessages(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	now := util.Timestamp()
	var channels []types.ChannelDS

	q := datastore.NewQuery(backend.DatastoreChannels).Filter("Next <", now)
	_, err := store.Client().GetAll(ctx, q, &channels)

	if err == nil {
		for i := range channels {

			id := channels[i].TeamID
			channel := channels[i].ID
			tasks.Schedule(ctx, backend.BackgroundWorkQueue, fmt.Sprintf("%s%s/j/messages?id=%s&c=%s", env.Getenv("BASE_URL", ""), api.JobsBaseURL, id, channel))
		}
	} else {
		errorreporting.Report(err)
	}
}

// scheduleHourlyTasks schedules all workspaces that need updating
func scheduleHourlyTasks(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	var workspaces []types.WorkspaceDS

	q := datastore.NewQuery(backend.DatastoreWorkspaces)
	_, err := store.Client().GetAll(ctx, q, &workspaces)

	if err == nil {
		for i := range workspaces {
			tasks.Schedule(ctx, backend.BackgroundWorkQueue, env.Getenv("BASE_URL", "")+api.JobsBaseURL+"/j/hourly?id="+workspaces[i].ID)
		}
	} else {
		errorreporting.Report(err)
	}
}

// scheduleDailyTasks schedules all workspaces that need updating
func scheduleDailyTasks(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	var workspaces []types.WorkspaceDS

	q := datastore.NewQuery(backend.DatastoreWorkspaces)
	_, err := store.Client().GetAll(ctx, q, &workspaces)

	if err == nil {
		for i := range workspaces {
			tasks.Schedule(ctx, backend.BackgroundWorkQueue, env.Getenv("BASE_URL", "")+api.JobsBaseURL+"/j/daily?id="+workspaces[i].ID)
		}
	} else {
		errorreporting.Report(err)
	}
}
