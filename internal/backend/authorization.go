package backend

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"

	"github.com/lnkk-ai/lnkk/internal/types"
	"github.com/lnkk-ai/lnkk/pkg/store"
	"github.com/majordomusio/commons/pkg/util"
	"google.golang.org/appengine/memcache"
)

// GetAuthorization returns the authorization granted for a workspace
func GetAuthorization(ctx context.Context, id string) *types.Authorization {
	var auth = types.Authorization{}
	key := "workspace.auth" + id
	_, err := memcache.Gob.Get(ctx, key, &auth)

	if err != nil {
		err = store.Client().Get(ctx, AuthorizationKey(ctx, id), &auth)
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

	var auth = types.Authorization{}
	key := AuthorizationKey(ctx, id)
	err := store.Client().Get(ctx, key, &auth)

	if err == nil {
		auth.AccessToken = token
		auth.Scope = scope
		auth.Updated = util.Timestamp()
	} else {
		auth = types.Authorization{
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

// AuthorizationKey creates a datastore key for a workspace authorization based on the team_id.
func AuthorizationKey(ctx context.Context, id string) *datastore.Key {
	return datastore.NameKey(DatastoreAuthorizations, id, nil)
}
