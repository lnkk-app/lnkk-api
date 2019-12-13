package main

import (
	"fmt"

	"cloud.google.com/go/datastore"
	"github.com/gin-gonic/gin"
	"github.com/majordomusio/commons/pkg/env"
	"github.com/majordomusio/commons/pkg/util"
	"google.golang.org/appengine"

	"github.com/lnkk-ai/lnkk/pkg/api"
	"github.com/majordomusio/platform/pkg/errorreporting"
	"github.com/majordomusio/platform/pkg/store"
	"github.com/majordomusio/platform/pkg/tasks"

	"github.com/lnkk-ai/lnkk/internal/backend"
	"github.com/lnkk-ai/lnkk/internal/types"
)

// updateWorkspaces schedules all workspaces that need updating
func updateWorkspaces(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	now := util.Timestamp()
	var workspaces []types.WorkspaceDS

	q := datastore.NewQuery(backend.DatastoreWorkspaces).Filter("Next <", now)
	_, err := store.Client().GetAll(ctx, q, &workspaces)

	if err == nil {
		for i := range workspaces {

			tasks.Schedule(ctx, backend.BackgroundWorkQueue, env.Getenv("BASE_URL", "")+api.JobsBaseURL+"/users?id="+workspaces[i].ID)
			tasks.Schedule(ctx, backend.BackgroundWorkQueue, env.Getenv("BASE_URL", "")+api.JobsBaseURL+"/channels?id="+workspaces[i].ID)

			backend.MarkWorkspaceUpdated(ctx, workspaces[i].ID)
		}
	} else {
		errorreporting.Report(err)
	}
}

// collectMessages schedules the collection of messages in a given team & channel
func collectMessages(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	now := util.Timestamp()
	var channels []types.ChannelDS

	q := datastore.NewQuery(backend.DatastoreChannels).Filter("Next <", now)
	_, err := store.Client().GetAll(ctx, q, &channels)

	if err == nil {
		for i := range channels {

			id := channels[i].TeamID
			channel := channels[i].ID
			tasks.Schedule(ctx, backend.BackgroundWorkQueue, fmt.Sprintf("%s%s/messages?id=%s&c=%s", env.Getenv("BASE_URL", ""), api.JobsBaseURL, id, channel))
		}
	} else {
		errorreporting.Report(err)
	}
}

// updateStats schedules all workspaces that need updating
func updateStats(c *gin.Context) {
}
