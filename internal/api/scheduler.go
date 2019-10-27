package api

import (
	"fmt"

	"cloud.google.com/go/datastore"
	"github.com/gin-gonic/gin"
	"github.com/majordomusio/commons/pkg/env"
	"github.com/majordomusio/commons/pkg/util"
	"google.golang.org/appengine"

	"github.com/lnkk-ai/lnkk/pkg/api"
	"github.com/majordomusio/platform/pkg/errorreporting"
	"github.com/majordomusio/platform/pkg/logger"
	"github.com/majordomusio/platform/pkg/store"
	"github.com/majordomusio/platform/pkg/tasks"

	"github.com/lnkk-ai/lnkk/internal/backend"
	"github.com/lnkk-ai/lnkk/internal/types"
)

// UpdateWorkspaces schedules all workspaces that need updating
func UpdateWorkspaces(c *gin.Context) {
	topic := "scheduler.update.workspace"
	ctx := appengine.NewContext(c.Request)

	now := util.Timestamp()
	var workspaces []types.WorkspaceDS

	q := datastore.NewQuery(backend.DatastoreWorkspaces).Filter("NextUpdate <", now)
	_, err := store.Client().GetAll(ctx, q, &workspaces)

	if err == nil {
		for i := range workspaces {

			tasks.Schedule(ctx, backend.BackgroundWorkQueue, env.Getenv("BASE_URL", "")+api.JobsBaseURL+"/users?id="+workspaces[i].ID)
			tasks.Schedule(ctx, backend.BackgroundWorkQueue, env.Getenv("BASE_URL", "")+api.JobsBaseURL+"/channels?id="+workspaces[i].ID)

			backend.MarkWorkspaceUpdated(ctx, workspaces[i].ID)
			logger.Info(topic, "workspace=%s", workspaces[i].ID)
		}
	} else {
		errorreporting.Report(err)
		logger.Critical(topic, err.Error())
	}
}

// CollectMessages schedules the collection of messages in a given team & channel
func CollectMessages(c *gin.Context) {
	topic := "scheduler.update.messages"
	ctx := appengine.NewContext(c.Request)

	now := util.Timestamp()
	var channels []types.ChannelDS

	q := datastore.NewQuery(backend.DatastoreChannels).Filter("NextUpdate <", now)
	_, err := store.Client().GetAll(ctx, q, &channels)

	if err == nil {
		for i := range channels {

			id := channels[i].TeamID
			channel := channels[i].ID
			tasks.Schedule(ctx, backend.BackgroundWorkQueue, fmt.Sprintf("%s%s/messages?id=%s&c=%s", env.Getenv("BASE_URL", ""), api.JobsBaseURL, id, channel))

			logger.Info(topic, "workspace=%s, channel=%s", id, channel)
		}
	} else {
		errorreporting.Report(err)
		logger.Critical(topic, err.Error())
	}
}
