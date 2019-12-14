package main

import (
	"fmt"

	"google.golang.org/appengine"

	"github.com/gin-gonic/gin"

	"github.com/majordomusio/commons/pkg/env"
	"github.com/majordomusio/commons/pkg/util"

	"github.com/majordomusio/platform/pkg/errorreporting"
	"github.com/majordomusio/platform/pkg/tasks"

	"github.com/lnkk-ai/lnkk/internal/backend"
	"github.com/lnkk-ai/lnkk/pkg/api"
	"github.com/lnkk-ai/lnkk/pkg/metrics"
	"github.com/lnkk-ai/lnkk/pkg/slack"
)

// updateUsersJob updates the list of users of a workspace
func updateUsersJob(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	id := c.Query("id")
	cursor := c.Query("cursor")

	auth, err := backend.GetAuthToken(ctx, id)
	if err != nil {
		errorreporting.Report(err)
		return
	}

	// update the list of users
	users, err := slack.UsersList(ctx, auth, cursor)
	if err != nil {
		errorreporting.Report(err)
		return
	}

	metrics.Count(ctx, "jobs.slack.users.count", id, len(users.Members))

	for i := range users.Members {
		err = backend.UpdateUser(ctx, users.Members[i].ID, users.Members[i].TeamID, users.Members[i].Name, users.Members[i].RealName, users.Members[i].Profile.FirstName, users.Members[i].Profile.LastName, users.Members[i].Profile.Email, users.Members[i].Deleted, users.Members[i].IsBot)

		if err != nil {
			errorreporting.Report(err)
		}
	}

	nextCursor := users.ResponseMetadata["next_cursor"]
	if nextCursor != "" {
		// there is more data, schedule its retrieval
		tasks.Schedule(ctx, backend.BackgroundWorkQueue, fmt.Sprintf(api.JobsBaseURL+"/users?id=%v&cursor=%v", id, nextCursor))
	}
}

// updateChannelsJob updates the workspace metadata periodically
func updateChannelsJob(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	id := c.Query("id")
	cursor := c.Query("cursor")

	auth, err := backend.GetAuthToken(ctx, id)
	if err != nil {
		errorreporting.Report(err)
		return
	}

	// update the list of channels
	channels, err := slack.ChannelsList(ctx, auth, cursor)
	if err != nil {
		errorreporting.Report(err)
		return
	}

	metrics.Count(ctx, "jobs.slack.channels.count", id, len(channels.Channels))

	for i := range channels.Channels {
		err := backend.UpdateChannel(ctx, channels.Channels[i].ID, id, channels.Channels[i].Name, channels.Channels[i].Topic.Value, channels.Channels[i].Purpose.Value, channels.Channels[i].IsArchived, channels.Channels[i].IsPrivate, false)

		if err != nil {
			errorreporting.Report(err)
		}
	}

	nextCursor := channels.ResponseMetadata["next_cursor"]
	if nextCursor != "" {
		// there is more data, schedule its retrieval
		tasks.Schedule(ctx, backend.BackgroundWorkQueue, fmt.Sprintf(api.JobsBaseURL+"/channels?id=%v&cursor=%v", id, nextCursor))
	}

}

// collectMessagesJob collects all new messages in a team & channel
// /_i/1/jobs/msgs?id=..&c=..&latest=..
func collectMessagesJob(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	id := c.Query("id")
	channel := c.Query("c")
	latest := c.Query("l")
	ts := ""

	auth, err := backend.GetAuthToken(ctx, id)
	if err != nil {
		errorreporting.Report(err)
		return
	}

	// setup the markers
	now := util.Timestamp()
	oldest := backend.GetChannelLatestCrawled(ctx, channel, id)

	if latest == "" {
		ts = fmt.Sprintf("%d", now)
	} else {
		ts = latest
	}

	// collect messages since ts
	msgs, err := slack.ChannelsHistory(ctx, auth, channel, backend.DefaultCrawlerBatchSize, ts)
	if err != nil {
		errorreporting.Report(err)
		return
	}

	n := 0
	for i := range msgs.Messages {

		if slack.Timestamp(msgs.Messages[i].TS) > oldest {
			backend.StoreMessage(ctx, channel, id, &msgs.Messages[i])
			n++
		} else {
			// we have reached the last known message
			backend.MarkChannelCrawled(ctx, channel, id, now)
			metrics.Count(ctx, "jobs.slack.messages.count", channel, n)
			return
		}
	}

	if msgs.HasMore {
		n := len(msgs.Messages)
		tasks.Schedule(ctx, backend.BackgroundWorkQueue, fmt.Sprintf("%s%s/messages?id=%s&c=%s&l=%s", env.Getenv("BASE_URL", ""), api.JobsBaseURL, id, channel, msgs.Messages[n-1].TS))
	} else {
		// we have reached the last known message
		backend.MarkChannelCrawled(ctx, channel, id, now)
	}

	// final auditing
	metrics.Count(ctx, "jobs.slack.messages.count", channel, n)

}
