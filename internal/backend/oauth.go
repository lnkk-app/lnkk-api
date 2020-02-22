package backend

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/appengine/memcache"

	"github.com/lnkk-ai/lnkk/pkg/errors"
	"github.com/lnkk-ai/lnkk/pkg/platform"
	"github.com/majordomusio/commons/pkg/env"
	"github.com/majordomusio/commons/pkg/util"
)

// FIXME change to platform cache

// GetAuthorization returns the authorization granted to an app
func GetAuthorization(ctx context.Context, id string) (*AuthorizationDS, error) {
	var auth = AuthorizationDS{}
	key := "workspace.auth" + id
	_, err := memcache.Gob.Get(ctx, key, &auth)

	if err != nil {
		err = platform.DataStore().Get(ctx, AuthorizationKey(id), &auth)
		if err == nil {
			cache := memcache.Item{}
			cache.Key = key
			cache.Object = auth
			cache.Expiration, _ = time.ParseDuration(DefaultCacheDuration)
			memcache.Gob.Set(ctx, &cache)
		} else {
			return nil, err
		}
	}

	return &auth, nil
}

// GetAuthToken returns the oauth token of the workspace integration
func GetAuthToken(ctx context.Context, id string) (string, error) {
	token := env.Getenv("SLACK_AUTH_TOKEN", "")
	if token != "" {
		return token, nil
	}

	auth, err := GetAuthorization(ctx, id)
	if err != nil {
		return "", err
	}
	if auth == nil {
		return "", errors.New(fmt.Sprintf("No authorization token for workspace '%s'", id))
	}
	return auth.AccessToken, nil
}

// UpdateAuthorization updates the authorization, or creates a new one.
func UpdateAuthorization(ctx context.Context, id, name, token, tokenType, scope, appID, botID string) error {
	now := util.Timestamp()
	var auth = AuthorizationDS{}
	key := AuthorizationKey(id)
	err := platform.DataStore().Get(ctx, key, &auth)

	if err == nil {
		auth.AccessToken = token
		auth.Scope = scope
		auth.Updated = now
	} else {
		auth = AuthorizationDS{
			ID:          id,
			Name:        name,
			AccessToken: token,
			TokenType:   tokenType,
			Scope:       scope,
			AppID:       appID,
			BotUserID:   botID,
			Created:     now,
			Updated:     now,
		}
	}

	_, err = platform.DataStore().Put(ctx, key, &auth)
	return err
}
