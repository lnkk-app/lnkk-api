package backend

import (
	"context"
	"time"

	"github.com/lnkk-ai/lnkk/internal/types"
	"github.com/lnkk-ai/lnkk/pkg/store"
	"github.com/majordomusio/commons/pkg/util"
	"google.golang.org/appengine/memcache"
)

// GetAuthorization returns the authorization granted for a workspace
func GetAuthorization(ctx context.Context, id string) *types.AuthorizationDS {
	var auth = types.AuthorizationDS{}
	key := "workspace.auth" + id
	_, err := memcache.Gob.Get(ctx, key, &auth)

	if err != nil {
		err = store.Client().Get(ctx, AuthorizationKey(id), &auth)
		if err == nil {
			cache := memcache.Item{}
			cache.Key = key
			cache.Object = auth
			cache.Expiration, _ = time.ParseDuration(DefaultCacheDuration)
			memcache.Gob.Set(ctx, &cache)
		} else {
			return nil
		}
	}

	return &auth
}

// GetAuthToken returns the oauth token of the workspace
func GetAuthToken(ctx context.Context, id string) string {
	auth := GetAuthorization(ctx, id)
	if auth != nil {
		return auth.AccessToken
	}
	return ""
}

// UpdateAuthorization updates the authorization, or creates a new one.
func UpdateAuthorization(ctx context.Context, id, name, token, scope, authorizingUser, installerUser string) error {

	var auth = types.AuthorizationDS{}
	key := AuthorizationKey(id)
	err := store.Client().Get(ctx, key, &auth)

	if err == nil {
		auth.AccessToken = token
		auth.Scope = scope
		auth.Updated = util.Timestamp()
	} else {
		auth = types.AuthorizationDS{
			ID:              id,
			Name:            name,
			AccessToken:     token,
			Scope:           scope,
			AuthorizingUser: authorizingUser,
			InstallerUser:   installerUser,
			Created:         util.Timestamp(),
			Updated:         util.Timestamp(),
		}
	}

	_, err = store.Client().Put(ctx, key, &auth)
	return err
}

// UpdateWorkspace updates the workspace metadata
func UpdateWorkspace(ctx context.Context, id, name string) error {

	var ws = types.WorkspaceDS{}
	key := WorkspaceKey(id)
	err := store.Client().Get(ctx, key, &ws)

	if err == nil {
		ws.NextUpdate = util.Timestamp() + (int64)(ws.UpdateSchedule)
		ws.Updated = util.Timestamp()
	} else {
		ws = types.WorkspaceDS{
			ID:             id,
			Name:           name,
			NextUpdate:     0,
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
		ws.NextUpdate = util.Timestamp() + (int64)(ws.UpdateSchedule)
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
			NextCrawl:       util.Timestamp() + (int64)(util.Random(DefaultCrawlerSchedule)),
			CrawlerSchedule: DefaultCrawlerSchedule,
			Created:         util.Timestamp(),
			Updated:         util.Timestamp(),
		}
	}

	_, err = store.Client().Put(ctx, key, &channel)
	return err
}
