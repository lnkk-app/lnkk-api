package backend

import (
	"context"

	"github.com/lnkk-ai/lnkk/internal/types"
	"github.com/majordomusio/commons/pkg/util"
	"github.com/majordomusio/platform/pkg/store"
)

// UpdateWorkspace updates the workspace metadata
func UpdateWorkspace(ctx context.Context, id, name string) error {

	var ws = types.WorkspaceDS{}
	key := WorkspaceKey(id)
	err := store.Client().Get(ctx, key, &ws)

	if err == nil {
		ws.Next = util.Timestamp() + (int64)(ws.UpdateSchedule)
		ws.Updated = util.Timestamp()
	} else {
		ws = types.WorkspaceDS{
			ID:             id,
			Name:           name,
			Next:           0,
			UpdateSchedule: DefaultUpdateSchedule,
			Created:        util.Timestamp(),
			Updated:        util.Timestamp(),
		}
	}

	_, err = store.Client().Put(ctx, key, &ws)
	return err
}

// MarkWorkspaceUpdated marks the auth record as updated
func MarkWorkspaceUpdated(ctx context.Context, id string) {
	var ws = types.WorkspaceDS{}
	key := WorkspaceKey(id)

	err := store.Client().Get(ctx, key, &ws)
	if err == nil {
		ws.Next = util.Timestamp() + (int64)(ws.UpdateSchedule)
		ws.Updated = util.Timestamp()
		_, err = store.Client().Put(ctx, key, &ws)
	}
}

// UpdateChannel updates the channel metadata
func UpdateChannel(ctx context.Context, id, team, name, topic, purpose string, archived, private, deleted bool) error {

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
		channel.Updated = util.Timestamp()
	} else {
		channel = types.ChannelDS{
			ID:              id,
			TeamID:          team,
			Name:            name,
			Topic:           topic,
			Purpose:         purpose,
			IsArchived:      archived,
			IsPrivate:       private,
			IsDeleted:       deleted,
			Latest:          0,
			Next:            util.Timestamp() + (int64)(util.Random(DefaultCrawlerSchedule)),
			CrawlerSchedule: DefaultCrawlerSchedule,
			Created:         util.Timestamp(),
			Updated:         util.Timestamp(),
		}
	}

	_, err = store.Client().Put(ctx, key, &channel)
	return err
}
