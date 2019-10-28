package backend

import (
	"context"

	"github.com/lnkk-ai/lnkk/internal/types"
	"github.com/lnkk-ai/lnkk/pkg/slack"
	"github.com/majordomusio/commons/pkg/util"
	"github.com/majordomusio/platform/pkg/errorreporting"
	"github.com/majordomusio/platform/pkg/store"
)

// StoreMessage stores a slack message
func StoreMessage(ctx context.Context, id, team string, message *slack.ChannelMessage) error {
	now := util.Timestamp()
	user := message.User
	ts := message.TS
	tsSeconds := slack.Timestamp(ts)
	attachments := false
	reactions := false
	msgID := MessageKeyString(id, ts, user)

	if len(message.Attachements) > 0 {
		attachments = true
	}

	if len(message.Reactions) > 0 {
		reactions = true
	}

	msg := types.MessageDS{
		ChannelID:    id,
		TeamID:       team,
		User:         user,
		Text:         message.Text,
		Timestamp:    slack.TimestampNano(ts),
		Type:         message.Type,
		Subtype:      message.Subtype,
		Attachements: attachments,
		Reactions:    reactions,
		Day:          util.TimestampToWeekday(tsSeconds),
		Hour:         util.TimestampToHour(tsSeconds),
		Created:      now,
		Updated:      now,
	}

	_, err := store.Client().Put(ctx, MessageKey(id, ts, user), &msg)

	// store the attachments if there are any
	if attachments {
		for i := range message.Attachements {
			StoreAttachement(ctx, id, team, msgID, message.Attachements[i].ID, message.Attachements[i].Text, message.Attachements[i].Fallback)
		}
	}

	if reactions {
		for i := range message.Reactions {
			StoreReaction(ctx, id, team, msgID, i, message.Reactions[i].Name, message.Reactions[i].Count, &message.Reactions[i].Users)
		}
	}

	return err
}

// StoreAttachement stores an attachment based on simple params
func StoreAttachement(ctx context.Context, channelID, teamID, msgID string, id int, text, fallback string) error {
	now := util.Timestamp()

	att := types.AttachmentDS{
		MessageID:   msgID,
		ChannelID:   channelID,
		TeamID:      teamID,
		Index:       id,
		Text:        text,
		Alternative: fallback,
		Created:     now,
		Updated:     now,
	}

	_, err := store.Client().Put(ctx, AttachmentKey(msgID, id), &att)
	return err
}

// StoreReaction stores a reaction
func StoreReaction(ctx context.Context, channelID, teamID, msgID string, id int, reaction string, count int, users *[]string) error {
	now := util.Timestamp()

	re := types.ReactionDS{
		MessageID: msgID,
		ChannelID: channelID,
		TeamID:    teamID,
		Reaction:  reaction,
		Count:     count,
		Users:     *users,
		Created:   now,
		Updated:   now,
	}

	_, err := store.Client().Put(ctx, ReactionKey(msgID, id), &re)
	if err != nil {
		// FIXME remove after test
		errorreporting.Report(err)
	}
	return err
}

// MarkChannelCrawled updates the craler data on a channel
func MarkChannelCrawled(ctx context.Context, id, team string, latest int64) {
	var channel = types.ChannelDS{}
	key := ChannelKey(id, team)
	err := store.Client().Get(ctx, key, &channel)

	if err == nil {
		channel.Latest = latest
		channel.Next = util.Timestamp() + (int64)(channel.Schedule)
		channel.Updated = util.Timestamp()
		_, err = store.Client().Put(ctx, key, &channel)
	}
}

// GetChannelLatestCrawled returns the ts of the latest message
func GetChannelLatestCrawled(ctx context.Context, id, team string) int64 {
	var channel = types.ChannelDS{}
	key := ChannelKey(id, team)
	err := store.Client().Get(ctx, key, &channel)

	if err != nil {
		return -1
	}

	return channel.Latest
}
