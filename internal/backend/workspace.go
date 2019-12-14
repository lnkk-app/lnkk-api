package backend

import (
	"context"

	"github.com/majordomusio/commons/pkg/util"

	"github.com/majordomusio/platform/pkg/store"

	"github.com/lnkk-ai/lnkk/internal/types"
)

// UpdateWorkspace updates the workspace metadata
func UpdateWorkspace(ctx context.Context, id, name string) error {
	now := util.Timestamp()
	var ws = types.WorkspaceDS{}
	key := WorkspaceKey(id)
	err := store.Client().Get(ctx, key, &ws)

	if err == nil {
		ws.Next = now + (int64)(ws.Schedule)
		ws.Updated = now
	} else {
		ws = types.WorkspaceDS{
			ID:       id,
			Name:     name,
			Next:     0,
			Schedule: DefaultUpdateSchedule,
			Created:  now,
			Updated:  now,
		}
	}

	_, err = store.Client().Put(ctx, key, &ws)
	return err
}

// MarkWorkspaceUpdated marks the auth record as updated
func MarkWorkspaceUpdated(ctx context.Context, id string) {
	now := util.Timestamp()
	var ws = types.WorkspaceDS{}
	key := WorkspaceKey(id)

	err := store.Client().Get(ctx, key, &ws)
	if err == nil {
		ws.Next = now + (int64)(ws.Schedule)
		ws.Updated = now
		_, err = store.Client().Put(ctx, key, &ws)
	}
}

// UpdateChannel updates the channel metadata
func UpdateChannel(ctx context.Context, id, team, name, topic, purpose string, archived, private, deleted bool) error {
	now := util.Timestamp()
	var channel = types.ChannelDS{}
	key := ChannelKey(id, team)
	err := store.Client().Get(ctx, key, &channel)

	if err == nil {
		channel.Name = name
		channel.Topic = topic
		channel.Purpose = purpose
		channel.IsArchived = archived
		channel.IsPrivate = private
		channel.IsDeleted = deleted
		channel.Updated = now
	} else {
		channel = types.ChannelDS{
			ID:         id,
			TeamID:     team,
			Name:       name,
			Topic:      topic,
			Purpose:    purpose,
			IsArchived: archived,
			IsPrivate:  private,
			IsDeleted:  deleted,
			Latest:     0,
			Next:       now + (int64)(util.Random(DefaultCrawlerSchedule)),
			Schedule:   DefaultCrawlerSchedule,
			Created:    now,
			Updated:    now,
		}
	}

	_, err = store.Client().Put(ctx, key, &channel)
	return err
}
