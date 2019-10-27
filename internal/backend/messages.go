package backend

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
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
		channel.NextCrawl = util.Timestamp() + (int64)(channel.CrawlerSchedule)
		channel.Updated = util.Timestamp()
		_, err = store.Client().Put(ctx, key, &channel)
	}
}

// GetChannelLatestCrawled returns the ts of the latest message
func GetChannelLatestCrawled(ctx context.Context, id, team string) int64 {
	var channel = types.ChannelDS{}
	key := ChannelKey(id, team)
	err := store.Client().Get(ctx, key, &channel)

	if err == nil {
		return channel.Latest
	} else {
		return -1
	}
}

// StoreSimpleMessage creates a new message
func StoreSimpleMessage(ctx context.Context, id, team, user, ts, message string) error {

	var msg = types.Message{}
	key := MessageKey(id, ts, user)
	err := store.Client().Get(ctx, key, &msg)

	if err == nil {
		msg.Updated = util.Timestamp()
	} else {
		msg = types.Message{
			ChannelID:       id,
			TeamID:          team,
			User:            user,
			TS:              slack.TimestampNano(ts),
			Text:            message,
			HasAttachements: false,
			Created:         util.Timestamp(),
			Updated:         util.Timestamp(),
		}
	}

	_, err = store.Client().Put(ctx, key, &msg)
	return err
}

// StoreSlackMessage stores a slack message
func StoreSlackMessage(ctx context.Context, id, team string, message *slack.ChannelMessage) error {
	user := message.User
	ts := message.TS
	attachments := false

	var msg = types.Message{}
	key := MessageKey(id, ts, user)
	err := store.Client().Get(ctx, key, &msg)

	if err == nil {
		msg.Updated = util.Timestamp()
	} else {
		if len(message.Attachements) > 0 {
			attachments = true
		}

		msg = types.Message{
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
	var att = types.Attachment{}
	key := AttachmentKey(msgID, id)
	err := store.Client().Get(ctx, key, &att)

	if err == nil {
		att.Updated = util.Timestamp()
	} else {
		att = types.Attachment{
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
func GetAttachment(ctx context.Context, msgID string, id int) (*types.Attachment, error) {
	var att = types.Attachment{}
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

// GetMessages retrieves a page of archived messages
func GetMessages(ctx context.Context, channelID, teamID string, page, pageSize int) (*[]types.ArchivedMessage, error) {
	var messages []types.Message

	q := datastore.NewQuery(DatastoreMessages).Filter("ChannelID =", channelID).Order("-TS").Offset((page - 1) * pageSize).Limit(pageSize)
	_, err := store.Client().GetAll(ctx, q, &messages)

	if err != nil {
		return nil, err
	}

	// the new array of messages
	m := make([]types.ArchivedMessage, len(messages))

	// render the messages
	for i := range messages {
		m[i].Text = messages[i].Text
		m[i].User = GetUserName(ctx, messages[i].User, teamID)
		m[i].Timestamp = messages[i].TS
		m[i].Created = util.TimestampToUTC(messages[i].TS / 1000000) // Slack TS is in nanoseconds !
		m[i].HasAttachments = messages[i].HasAttachements

		if messages[i].HasAttachements {
			k := MessageKeyString(channelID, slack.TimestampNanoString(messages[i].TS), messages[i].User)
			att, err := GetAttachment(ctx, k, 1)
			if err == nil {
				m[i].Attachments = make([]types.ArchivedMessageAttachment, 1)
				m[i].Attachments[0] = types.ArchivedMessageAttachment{
					Text:         att.Text,
					FallbackText: att.Fallback,
				}
			}
		}
	}

	// done
	return &m, nil
}
