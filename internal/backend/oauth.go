package backend

import (
	"context"
	"fmt"

	"github.com/lnkk-ai/lnkk/pkg/platform"
	"github.com/majordomusio/commons/pkg/env"
	"github.com/majordomusio/commons/pkg/util"
)

// GetAuthorization returns the authorization granted to an app
func GetAuthorization(ctx context.Context, id string) (*AuthorizationDS, error) {
	var auth = AuthorizationDS{}

	// just load it, let caching be handled elsewhere ...
	err := platform.DataStore().Get(ctx, AuthorizationKey(id), &auth)
	if err != nil {
		return nil, err
	}

	return &auth, nil
}

// GetAuthToken returns the oauth token of the workspace integration
func GetAuthToken(ctx context.Context, id string) (string, error) {
	// ENV always overrides anything stored...
	token := env.Getenv("SLACK_AUTH_TOKEN", "")
	if token != "" {
		return token, nil
	}

	// check the in-memory cache
	key := cacheKey(id)
	token, err := platform.Get(ctx, key)
	if token != "" {
		return token, nil
	}

	auth, err := GetAuthorization(ctx, id)
	if err != nil {
		return "", err
	}
	if auth == nil {
		return "", fmt.Errorf("No authorization token for workspace '%s'", id)
	}

	// add the token to the cache
	platform.Set(ctx, key, auth.AccessToken, 1800)

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

	// remove the entry from the cache if it is already there ...
	platform.Invalidate(ctx, cacheKey(id))

	_, err = platform.DataStore().Put(ctx, key, &auth)
	return err
}

func cacheKey(id string) string {
	return "workspace.oauth." + id
}
