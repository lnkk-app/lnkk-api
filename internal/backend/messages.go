package backend

import (
	"context"
	"time"

	"github.com/lnkk-ai/lnkk/internal/types"
	"github.com/lnkk-ai/lnkk/pkg/slack"
	"github.com/majordomusio/commons/pkg/util"
	"github.com/majordomusio/platform/pkg/store"

	"google.golang.org/appengine/memcache"
)

// MarkChannelCrawled updates the craler data on a channel
func MarkChannelCrawled(ctx context.Context, id, team string, latest int64) {
	var channel = types.ChannelDS{}
	key := ChannelKey(id, team)
	err := store.Client().Get(ctx, key, &channel)

	if err == nil {
		channel.Latest = latest
		channel.Next = util.Timestamp() + (int64)(channel.CrawlerSchedule)
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

// StoreSlackMessage stores a slack message
func StoreSlackMessage(ctx context.Context, id, team string, message *slack.ChannelMessage) error {
	user := message.User
	ts := message.TS
	attachments := false

	var msg = types.MessageDS{}
	key := MessageKey(id, ts, user)
	err := store.Client().Get(ctx, key, &msg)

	if err == nil {
		msg.Updated = util.Timestamp()
	} else {
		if len(message.Attachements) > 0 {
			attachments = true
		}

		msg = types.MessageDS{
			ChannelID:       id,
			TeamID:          team,
			User:            user,
			TS:              slack.TimestampNano(ts),
			Text:            message.Text,
			HasAttachements: attachments,
			Created:         util.Timestamp(),
			Updated:         util.Timestamp(),
		}
	}

	_, err = store.Client().Put(ctx, key, &msg)

	// store the attachments if there are any
	if attachments {
		msgID := MessageKeyString(id, ts, user)
		for i := range message.Attachements {
			StoreAttachement(ctx, id, team, msgID, message.Attachements[i].ID, message.Attachements[i].Text, message.Attachements[i].Fallback)
		}
	}

	return err
}

// StoreAttachement stores an attachment based on simple params
func StoreAttachement(ctx context.Context, channelID, teamID, msgID string, id int, text, fallback string) error {
	var att = types.AttachmentDS{}
	key := AttachmentKey(msgID, id)
	err := store.Client().Get(ctx, key, &att)

	if err == nil {
		att.Updated = util.Timestamp()
	} else {
		att = types.AttachmentDS{
			MessageID: msgID,
			ChannelID: channelID,
			TeamID:    teamID,
			ID:        id,
			Text:      text,
			Fallback:  fallback,
			Created:   util.Timestamp(),
			Updated:   util.Timestamp(),
		}
	}

	_, err = store.Client().Put(ctx, key, &att)
	return err
}

// GetAttachment gets and caches an attachment
func GetAttachment(ctx context.Context, msgID string, id int) (*types.AttachmentDS, error) {
	var att = types.AttachmentDS{}
	key := AttachmentKeyString(msgID, id)
	_, err := memcache.Gob.Get(ctx, key, &att)

	if err != nil {
		err := store.Client().Get(ctx, AttachmentKey(msgID, id), &att)
		if err == nil {
			cache := memcache.Item{}
			cache.Key = key
			cache.Object = att
			cache.Expiration, _ = time.ParseDuration(DefaultCacheDuration)
			memcache.Gob.Set(ctx, &cache)
		} else {
			return nil, err
		}
	}

	return &att, nil
}
